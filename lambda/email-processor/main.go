package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type EmailProcessor struct {
	s3Client        *s3.S3
	rawBucket       string
	processedBucket string
}

type ProcessedEmail struct {
	OriginalKey   string    `json:"original_key"`
	ProcessedAt   time.Time `json:"processed_at"`
	Size          int64     `json:"size"`
	ContentType   string    `json:"content_type"`
	MessageID     string    `json:"message_id,omitempty"`
	Subject       string    `json:"subject,omitempty"`
	From          string    `json:"from,omitempty"`
	To            string    `json:"to,omitempty"`
}

func NewEmailProcessor() (*EmailProcessor, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	return &EmailProcessor{
		s3Client:        s3.New(sess),
		rawBucket:       os.Getenv("RAW_BUCKET"),
		processedBucket: os.Getenv("PROCESSED_BUCKET"),
	}, nil
}

func (ep *EmailProcessor) handleS3Event(ctx context.Context, s3Event events.S3Event) error {
	for _, record := range s3Event.Records {
		if err := ep.processEmailObject(ctx, record.S3.Bucket.Name, record.S3.Object.Key); err != nil {
			log.Printf("Error processing object %s/%s: %v", record.S3.Bucket.Name, record.S3.Object.Key, err)
			return err
		}
	}
	return nil
}

func (ep *EmailProcessor) processEmailObject(ctx context.Context, bucketName, objectKey string) error {
	log.Printf("Processing email object: %s/%s", bucketName, objectKey)

	// Get the email object from S3
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	result, err := ep.s3Client.GetObjectWithContext(ctx, getObjectInput)
	if err != nil {
		return fmt.Errorf("failed to get object: %v", err)
	}
	defer result.Body.Close()

	// Read email content for basic parsing
	emailContent := make([]byte, *result.ContentLength)
	_, err = result.Body.Read(emailContent)
	if err != nil {
		return fmt.Errorf("failed to read email content: %v", err)
	}

	// Extract basic email metadata
	emailStr := string(emailContent)
	metadata := ep.extractEmailMetadata(emailStr)

	// Create processed email record
	processedEmail := ProcessedEmail{
		OriginalKey: objectKey,
		ProcessedAt: time.Now(),
		Size:        *result.ContentLength,
		ContentType: aws.StringValue(result.ContentType),
		MessageID:   metadata["message-id"],
		Subject:     metadata["subject"],
		From:        metadata["from"],
		To:          metadata["to"],
	}

	// Generate processed object key
	processedKey := ep.generateProcessedKey(objectKey)

	// Save processed email metadata
	metadataJSON, err := json.Marshal(processedEmail)
	if err != nil {
		return fmt.Errorf("failed to marshal processed email metadata: %v", err)
	}

	// Copy original email to processed bucket
	copySource := fmt.Sprintf("%s/%s", bucketName, objectKey)
	_, err = ep.s3Client.CopyObjectWithContext(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(ep.processedBucket),
		Key:        aws.String(processedKey),
		CopySource: aws.String(copySource),
		Metadata: map[string]*string{
			"processed-at":   aws.String(time.Now().Format(time.RFC3339)),
			"original-key":   aws.String(objectKey),
			"message-id":     aws.String(metadata["message-id"]),
			"subject":        aws.String(metadata["subject"]),
		},
		MetadataDirective: aws.String(s3.MetadataDirectiveReplace),
	})
	if err != nil {
		return fmt.Errorf("failed to copy email to processed bucket: %v", err)
	}

	// Save metadata as JSON file
	metadataKey := strings.Replace(processedKey, ".eml", ".json", 1)
	_, err = ep.s3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(ep.processedBucket),
		Key:         aws.String(metadataKey),
		Body:        strings.NewReader(string(metadataJSON)),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return fmt.Errorf("failed to save metadata: %v", err)
	}

	log.Printf("Successfully processed email: %s -> %s", objectKey, processedKey)
	return nil
}

func (ep *EmailProcessor) extractEmailMetadata(emailContent string) map[string]string {
	metadata := make(map[string]string)
	lines := strings.Split(emailContent, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			break // End of headers
		}

		if strings.HasPrefix(strings.ToLower(line), "message-id:") {
			metadata["message-id"] = strings.TrimSpace(line[11:])
		} else if strings.HasPrefix(strings.ToLower(line), "subject:") {
			metadata["subject"] = strings.TrimSpace(line[8:])
		} else if strings.HasPrefix(strings.ToLower(line), "from:") {
			metadata["from"] = strings.TrimSpace(line[5:])
		} else if strings.HasPrefix(strings.ToLower(line), "to:") {
			metadata["to"] = strings.TrimSpace(line[3:])
		}
	}

	return metadata
}

func (ep *EmailProcessor) generateProcessedKey(originalKey string) string {
	// Generate processed key with timestamp
	timestamp := time.Now().Format("2006/01/02/15")
	filename := filepath.Base(originalKey)
	return fmt.Sprintf("processed/%s/%s", timestamp, filename)
}

func handleRequest(ctx context.Context, s3Event events.S3Event) error {
	processor, err := NewEmailProcessor()
	if err != nil {
		return fmt.Errorf("failed to create email processor: %v", err)
	}

	return processor.handleS3Event(ctx, s3Event)
}

func main() {
	lambda.Start(handleRequest)
}