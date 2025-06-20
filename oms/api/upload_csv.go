package api

import (
	// "fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/log"
	gooms3 "github.com/omniful/go_commons/s3"
)

func (h *Handlers) UploadLocalCSVs(c *gin.Context) {
	// 1) Local CSV directory path
	localFolder := "C:/Users/dhruv/Desktop/omni_project/omni_project/oms/csv"

	// 2) Get list of *.csv files
	files, err := filepath.Glob(filepath.Join(localFolder, "*.csv"))
	if err != nil {
		log.Errorf("❌ Failed to list local CSVs: %v", err)
		c.JSON(500, gin.H{"error": "Failed to read local folder"})
		return
	}

	if len(files) == 0 {
		c.JSON(200, gin.H{"message": "No CSV files found"})
		return
	}

	// 3) Init GoCommons S3 client
	s3Client, err := gooms3.NewDefaultAWSS3Client()
	if err != nil {
		log.Errorf("❌ S3 client init failed: %v", err)
		c.JSON(500, gin.H{"error": "S3 client init failed"})
		return
	}

	bucket := h.OrderService.S3.Bucket
	var uploaded []string
	var failed []string

	// 4) Upload and validate each file
	for _, filePath := range files {
		fileName := filepath.Base(filePath)
		key := "uploads/" + fileName

		f, err := os.Open(filePath)
		if err != nil {
			log.Warnf("⚠️ Could not open %s: %v", fileName, err)
			failed = append(failed, fileName)
			continue
		}
		defer f.Close()

		// Upload to S3
		_, err = s3Client.PutObject(c.Request.Context(), &s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   f,
		})
		if err != nil {
			log.Warnf("❌ Upload failed for %s: %v", fileName, err)
			failed = append(failed, fileName)
			continue
		}
		log.Infof("✅ Uploaded to S3: %s", key)

		// Validate using OrderService
		if err := h.OrderService.ProcessCSV(c.Request.Context(), key); err != nil {
			log.Warnf("❌ Validation failed for %s: %v", fileName, err)
			failed = append(failed, fileName)
			continue
		}

		uploaded = append(uploaded, fileName)
	}

	// 5) Return result
	c.JSON(202, gin.H{
		"uploaded": uploaded,
		"failed":   failed,
	})
}
