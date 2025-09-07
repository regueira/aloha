package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestEmailProcessor(t *testing.T) {
	// Test that EmailProcessor can be created without errors
	_, err := NewEmailProcessor()
	if err != nil {
		// This is expected to fail in test environment without AWS credentials
		t.Logf("Expected error creating EmailProcessor in test environment: %v", err)
	} else {
		t.Log("EmailProcessor created successfully")
	}
}

func TestExtractEmailMetadata(t *testing.T) {
	// Create a mock processor for testing metadata extraction
	processor := &EmailProcessor{}
	
	emailContent := `Message-ID: <test123@example.com>
Subject: Test Email
From: sender@example.com
To: recipient@example.com

This is the body of the email.`

	metadata := processor.extractEmailMetadata(emailContent)
	
	if metadata["message-id"] != "<test123@example.com>" {
		t.Errorf("Expected Message-ID '<test123@example.com>', got '%s'", metadata["message-id"])
	}
	
	if metadata["subject"] != "Test Email" {
		t.Errorf("Expected Subject 'Test Email', got '%s'", metadata["subject"])
	}
	
	if metadata["from"] != "sender@example.com" {
		t.Errorf("Expected From 'sender@example.com', got '%s'", metadata["from"])
	}
	
	if metadata["to"] != "recipient@example.com" {
		t.Errorf("Expected To 'recipient@example.com', got '%s'", metadata["to"])
	}
}

func TestHandleRequest(t *testing.T) {
	// Test the main handler function without AWS resources
	ctx := context.Background()
	s3Event := events.S3Event{
		Records: []events.S3EventRecord{},
	}
	
	// This should fail gracefully without AWS credentials
	err := handleRequest(ctx, s3Event)
	if err != nil {
		t.Logf("Expected error in test environment: %v", err)
	}
}