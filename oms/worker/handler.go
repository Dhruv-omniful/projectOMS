package worker

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	// "io"
	"path"
	"strconv"
	"time"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	// "github.com/omniful/go_commons/config"
	commoncsv "github.com/omniful/go_commons/csv"
	"github.com/omniful/go_commons/log"
	// gooms3 "github.com/omniful/go_commons/s3"
	"github.com/omniful/go_commons/sqs"
	"github.com/dhruv/oms/client"
	"github.com/dhruv/oms/model"
)

type queueHandler struct{}

func (h *queueHandler) Process(ctx context.Context, msgs *[]sqs.Message) error {
	logger := log.DefaultLogger()

	s3Client, err := client.NewS3Client(ctx)
if err != nil {
	logger.Errorf("❌ Failed to init S3 client: %v", err)
	return err
}

	if err != nil {
		logger.Errorf("❌ Failed to init S3 client: %v", err)
		return err
	}

	for _, msg := range *msgs {
		var evt struct {
			Bucket string `json:"Bucket"`
			Key    string `json:"Key"`
		}
		if err := json.Unmarshal(msg.Value, &evt); err != nil {
			logger.Errorf("❌ Invalid SQS JSON: %v", err)
			return err
		}
		if evt.Bucket == "" || evt.Key == "" {
		logger.Errorf("❌ Missing Bucket or Key in SQS message: %+v", evt)
		continue
	}


		logger.Infof("📥 Processing file from S3: %s/%s", evt.Bucket, evt.Key)

		out, err := s3Client.Client.GetObject(ctx, &awss3.GetObjectInput{
			Bucket: awsV2.String(evt.Bucket),
			Key:    awsV2.String(evt.Key),
		})
		if err != nil {
			logger.Errorf("❌ Failed to get S3 object: %v", err)
			return err
		}
		defer out.Body.Close()

		r := csv.NewReader(out.Body)
		r.Comma = commoncsv.CsvDelimiter
		r.LazyQuotes = true

		header, err := r.Read()
		if err != nil {
			logger.Errorf("❌ Failed to read header: %v", err)
			return err
		}
		rows, err := r.ReadAll()
		if err != nil {
			logger.Errorf("❌ Failed to read CSV rows: %v", err)
			return err
		}

		idx := make(map[string]int)
		for i, col := range header {
			idx[col] = i
		}

		var invalid [][]string

		for _, row := range rows {
			qty, err := strconv.Atoi(row[idx["quantity"]])
			if err != nil || qty <= 0 {
				invalid = append(invalid, row)
				continue
			}

			order := &model.Order{
				TenantID: row[idx["tenant_id"]],
				SellerID: row[idx["seller_id"]],
				HubID:    row[idx["hub_id"]],
				SKUID:    row[idx["sku_id"]],
				Quantity: int64(qty),
			}

			if err := client.SaveOrder(ctx, order); err != nil {
				invalid = append(invalid, row)
				continue
			}

			logger.Infof("✅ Order processed: %+v", order)
		}

		if len(invalid) > 0 {
			buf := &bytes.Buffer{}
			w := csv.NewWriter(buf)
			w.Write(header)
			w.WriteAll(invalid)
			w.Flush()

			errKey := fmt.Sprintf("errors/%s-%d.csv", path.Base(evt.Key), time.Now().Unix())

			if _, err := s3Client.Client.PutObject(ctx, &awss3.PutObjectInput{
				Bucket: awsV2.String(evt.Bucket),
				Key:    awsV2.String(errKey),
				Body:   bytes.NewReader(buf.Bytes()),
			}); err != nil {
				logger.Errorf("❌ Failed to upload invalid CSV: %v", err)
			} else {
				logger.Infof("⚠️ Invalid rows saved to: s3://%s/%s", evt.Bucket, errKey)
			}
		}
	}

	return nil
}
