# Aloha Email Processor

A AWS CDK project built with Go 1.24 that sets up an automated email processing system using SES, S3, and Lambda.

## Architecture

This project creates the following AWS infrastructure:

```
SES Email Reception
        ↓
S3 Raw Emails Bucket (RAW)
        ↓ (S3 Event Trigger)
Lambda Email Processor
        ↓
S3 Processed Emails Bucket
```

### Components

1. **S3 Raw Emails Bucket**: Stores incoming emails from SES in their raw format
2. **S3 Processed Emails Bucket**: Stores processed emails with extracted metadata
3. **Lambda Email Processor**: Go-based function that processes raw emails and extracts metadata
4. **SES Configuration**: Email receiving rules to save emails to S3
5. **IAM Roles**: Proper permissions for Lambda to access S3 and SES

## Prerequisites

- Go 1.24+
- AWS CDK v2
- AWS CLI configured with appropriate permissions
- Node.js (for CDK CLI)

## Installation

1. Clone this repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Install AWS CDK CLI:
   ```bash
   npm install -g aws-cdk
   ```

4. Bootstrap CDK (if not already done):
   ```bash
   cdk bootstrap
   ```

## Build and Deploy

1. Build the Lambda function:
   ```bash
   ./build-lambda.sh
   ```

2. Deploy the infrastructure:
   ```bash
   cdk deploy
   ```

3. Verify the deployment:
   ```bash
   cdk diff
   ```

## Configuration

### Environment Variables

The Lambda function uses the following environment variables (automatically set by CDK):

- `RAW_BUCKET`: Name of the S3 bucket for raw emails
- `PROCESSED_BUCKET`: Name of the S3 bucket for processed emails
- `AWS_REGION`: AWS region where the function is deployed

### SES Domain Setup

To receive emails, you need to:

1. Verify your domain in SES
2. Update the SES receipt rule in `main.go` to match your domain:
   ```go
   Recipients: &[]*string{
       jsii.String("inbox@yourdomain.com"), // Replace with your domain
   },
   ```
3. Set up MX records for your domain to point to AWS SES

## Email Processing Flow

1. **Email Reception**: SES receives an email and stores it in the raw emails S3 bucket
2. **Trigger**: S3 event triggers the Lambda function when a new email is stored
3. **Processing**: Lambda function:
   - Reads the raw email from S3
   - Extracts metadata (Message-ID, Subject, From, To)
   - Copies the email to the processed bucket with timestamp-based organization
   - Saves metadata as a JSON file
4. **Storage**: Processed emails are organized by date and time in the processed bucket

### Processed Email Structure

```
processed-emails-bucket/
├── processed/2024/01/15/14/
│   ├── email-12345.eml          # Original email
│   └── email-12345.json         # Extracted metadata
```

### Metadata Format

```json
{
  "original_key": "emails/email-12345.eml",
  "processed_at": "2024-01-15T14:30:00Z",
  "size": 1024,
  "content_type": "text/plain",
  "message_id": "<12345@example.com>",
  "subject": "Test Email",
  "from": "sender@example.com",
  "to": "recipient@example.com"
}
```

## Monitoring and Logs

- Lambda function logs are available in CloudWatch Logs
- S3 access logs can be enabled for audit trails
- CloudWatch metrics are automatically created for Lambda executions

## Cost Optimization

- Raw emails are automatically deleted after 90 days
- Processed emails are retained for 365 days
- S3 versioning is enabled for data protection

## Security

- Lambda function uses least-privilege IAM roles
- S3 buckets are private by default
- Email processing uses secure AWS SDK operations

## Development

### Project Structure

```
.
├── main.go                    # CDK application entry point
├── cdk.json                   # CDK configuration
├── build-lambda.sh            # Lambda build script
├── lambda/
│   └── email-processor/
│       ├── main.go           # Lambda function code
│       ├── go.mod            # Lambda dependencies
│       └── bootstrap         # Compiled Lambda binary
├── go.mod                    # CDK dependencies
└── README.md                 # This file
```

### Testing Locally

You can test the Lambda function locally using AWS SAM or by creating test S3 events.

### Adding New Features

1. Update the Lambda function in `lambda/email-processor/main.go`
2. Rebuild with `./build-lambda.sh`
3. Update CDK infrastructure if needed in `main.go`
4. Deploy with `cdk deploy`

## Cleanup

To remove all resources:

```bash
cdk destroy
```

## Troubleshooting

### Common Issues

1. **Lambda function fails to process emails**: Check CloudWatch logs for detailed error messages
2. **SES not saving emails to S3**: Verify SES receipt rules and S3 bucket permissions
3. **CDK deploy fails**: Ensure AWS CLI is configured and you have necessary permissions

### Debugging

Enable detailed logging in the Lambda function by setting the log level:

```go
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## Support

For issues and questions, please open a GitHub issue or contact the maintainers.