package worker

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"path"
	"strconv"
	"time"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	commoncsv "github.com/omniful/go_commons/csv"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/sqs"

	"github.com/dhruv/oms/client"
	"github.com/dhruv/oms/model"
)

type queueHandler struct {
	IMS *client.IMSClient
}

func NewQueueHandler(ims *client.IMSClient) *queueHandler {
	return &queueHandler{
		IMS: ims,
	}
}

func (h *queueHandler) Process(ctx context.Context, msgs *[]sqs.Message) (err error) {
	logger := log.DefaultLogger()
	
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("🔥 Recovered from panic in Process: %v", r)
		}
	}()

	s3Client, err := client.NewS3Client(ctx)
	if err != nil {
		logger.Errorf("❌ Failed to init S3 client: %v", err)
		return err
	}
	logger.Infof("✅ S3 client initialized")

	for _, msg := range *msgs {
		logger.Infof("📨 Processing SQS message: %s", string(msg.Value))

		var evt struct {
			Bucket string `json:"Bucket"`
			Key    string `json:"Key"`
		}
		if err := json.Unmarshal(msg.Value, &evt); err != nil {
			logger.Errorf("❌ Invalid SQS JSON: %v", err)
			continue
		}
		if evt.Bucket == "" || evt.Key == "" {
			logger.Errorf("❌ Missing Bucket or Key in SQS message: %+v", evt)
			continue
		}

		logger.Infof("📥 Fetching file from S3: %s/%s", evt.Bucket, evt.Key)

		out, err := s3Client.Client.GetObject(ctx, &awss3.GetObjectInput{
			Bucket: awsV2.String(evt.Bucket),
			Key:    awsV2.String(evt.Key),
		})
		if err != nil {
			logger.Errorf("❌ Failed to get S3 object: %v", err)
			continue
		}
		defer out.Body.Close()

		r := csv.NewReader(out.Body)
		r.Comma = commoncsv.CsvDelimiter
		r.LazyQuotes = true

		header, err := r.Read()
		if err != nil {
			logger.Errorf("❌ Failed to read header: %v", err)
			continue
		}
		logger.Infof("✅ CSV header read: %v", header)

		rows, err := r.ReadAll()
		if err != nil {
			logger.Errorf("❌ Failed to read CSV rows: %v", err)
			continue
		}
		logger.Infof("✅ CSV rows count: %d", len(rows))

		idx := make(map[string]int)
		for i, col := range header {
			idx[col] = i
		}
		logger.Infof("✅ CSV column index map: %+v", idx)

		var invalid [][]string

		for rowNum, row := range rows {
			fmt.Printf("➡️ Row %d: %v\n", rowNum+1, row) 
			logger.Infof("➡️ Row %d only number log", rowNum+1)
			logger.Debugf("➡️ Row content dump: %#v", row)

			if len(row) < len(header) {
				logger.Warnf("⚠️ Row %d has insufficient columns: %v", rowNum+1, row)
				invalid = append(invalid, row)
				continue
			}

			getVal := func(col string) string {
				i, ok := idx[col]
				if !ok || i >= len(row) {
					logger.Warnf("⚠️ Missing or invalid index for column '%s' at row %d", col, rowNum+1)
					return ""
				}
				return row[i]
			}

			qtyStr := getVal("quantity")
			qty, err := strconv.Atoi(qtyStr)
			if err != nil || qty <= 0 {
				logger.Warnf("⚠️ Invalid quantity at row %d: %s", rowNum+1, qtyStr)
				invalid = append(invalid, row)
				continue
			}

			skuID := getVal("sku_id")
			hubID := getVal("hub_id")

			if h.IMS == nil {
				logger.Errorf("❌ IMS client is nil at row %d", rowNum+1)
				invalid = append(invalid, row)
				continue
			}

			logger.Debugf("🔍 Calling CheckSKU on %s", skuID)
			isValidSKU := h.IMS.CheckSKU(ctx, skuID)
			logger.Debugf("✅ CheckSKU(%s) -> %v", skuID, isValidSKU)

			logger.Debugf("🔍 Calling CheckHub on %s", hubID)
			isValidHub := h.IMS.CheckHub(ctx, hubID)
			logger.Debugf("✅ CheckHub(%s) -> %v", hubID, isValidHub)

			if !isValidSKU || !isValidHub {
				logger.Warnf("⚠️ Invalid SKU or Hub at row %d: SKU=%s Hub=%s", rowNum+1, skuID, hubID)
				invalid = append(invalid, row)
				continue
			}

			order := &model.Order{
				TenantID: getVal("tenant_id"),
				SellerID: getVal("seller_id"),
				HubID:    hubID,
				SKUID:    skuID,
				Quantity: int64(qty),
			}
			logger.Debugf("💾 Attempting to save order: %+v", order)
			if err := client.SaveOrder(ctx, order); err != nil {
				logger.Errorf("❌ Failed to save order at row %d: %v", rowNum+1, err)
				invalid = append(invalid, row)
				continue
			}
			logger.Infof("✅ Order processed at row %d: %+v", rowNum+1, order)
			client.PublishOrderCreated(ctx, order)

			client.NotifyWebhooks(ctx, order.TenantID, "order.created", order)
		}

		if len(invalid) > 0 {
			logger.Warnf("⚠️ Found %d invalid rows, uploading to S3", len(invalid))
			buf := &bytes.Buffer{}
			w := csv.NewWriter(buf)
			w.Write(header)
			w.WriteAll(invalid)
			w.Flush()

			errKey := fmt.Sprintf("errors/%s-%d.csv", path.Base(evt.Key), time.Now().Unix())

			_, err := s3Client.Client.PutObject(ctx, &awss3.PutObjectInput{
				Bucket: awsV2.String(evt.Bucket),
				Key:    awsV2.String(errKey),
				Body:   bytes.NewReader(buf.Bytes()),
			})
			if err != nil {
				logger.Errorf("❌ Failed to upload invalid CSV: %v", err)
			} else {
				logger.Infof("⚠️ Invalid rows saved to: s3://%s/%s", evt.Bucket, errKey)
			}
		}
	}

	return nil
}
