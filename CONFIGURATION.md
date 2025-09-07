# Aloha Email Processor Configuration Examples

## Environment Variables

Set these environment variables before deployment:

```bash
# AWS Account and Region (for CDK deployment)
export CDK_DEFAULT_ACCOUNT=123456789012
export CDK_DEFAULT_REGION=us-east-1

# AWS Credentials (if not using AWS CLI profiles)
export AWS_ACCESS_KEY_ID=your-access-key
export AWS_SECRET_ACCESS_KEY=your-secret-key
```

## CDK Context Configuration

You can customize the CDK deployment by modifying `cdk.json`:

```json
{
  "app": "go run main.go",
  "context": {
    "email-processor:domain": "yourdomain.com",
    "email-processor:environment": "production",
    "email-processor:retention-days": 90
  }
}
```

## Lambda Environment Variables

The Lambda function automatically receives these environment variables from CDK:

- `RAW_BUCKET`: Name of the S3 bucket for raw emails
- `PROCESSED_BUCKET`: Name of the S3 bucket for processed emails
- `AWS_REGION`: AWS region where the function is deployed

## SES Receipt Rule Configuration

Example AWS CLI command to create SES receipt rules:

```bash
# Replace YOUR_DOMAIN and YOUR_BUCKET_NAME with actual values
DOMAIN="yourdomain.com"
RAW_BUCKET="aloha-email-processor-rawemailsbucket-xxxxx"

# Create receipt rule set
aws ses create-receipt-rule-set --rule-set-name aloha-email-rules

# Create receipt rule for main inbox
aws ses create-receipt-rule \
  --rule-set-name aloha-email-rules \
  --rule '{
    "Name": "main-inbox",
    "Recipients": ["inbox@'${DOMAIN}'", "hello@'${DOMAIN}'"],
    "Enabled": true,
    "Actions": [
      {
        "S3Action": {
          "BucketName": "'${RAW_BUCKET}'",
          "ObjectKeyPrefix": "emails/"
        }
      }
    ]
  }'

# Create receipt rule for support emails
aws ses create-receipt-rule \
  --rule-set-name aloha-email-rules \
  --rule '{
    "Name": "support-inbox",
    "Recipients": ["support@'${DOMAIN}'"],
    "Enabled": true,
    "Actions": [
      {
        "S3Action": {
          "BucketName": "'${RAW_BUCKET}'",
          "ObjectKeyPrefix": "emails/support/"
        }
      }
    ]
  }'

# Set as active rule set
aws ses set-active-receipt-rule-set --rule-set-name aloha-email-rules
```

## Domain DNS Configuration

Add these DNS records to your domain:

### MX Record
```
Type: MX
Name: @  (or your domain)
Value: inbound-smtp.us-east-1.amazonaws.com
Priority: 10
TTL: 3600
```

### TXT Records for Domain Verification
```
Type: TXT
Name: _amazonses.yourdomain.com
Value: [provided by AWS SES console]
TTL: 1800
```

### SPF Record (Optional)
```
Type: TXT
Name: @
Value: "v=spf1 include:amazonses.com ~all"
TTL: 3600
```

### DKIM Records (Optional)
```
Type: CNAME
Name: [dkim-key-1]._domainkey
Value: [dkim-value-1].dkim.amazonses.com
TTL: 1800

Type: CNAME  
Name: [dkim-key-2]._domainkey
Value: [dkim-value-2].dkim.amazonses.com
TTL: 1800

Type: CNAME
Name: [dkim-key-3]._domainkey  
Value: [dkim-value-3].dkim.amazonses.com
TTL: 1800
```

## S3 Bucket Organization

The system organizes emails in the following structure:

### Raw Emails Bucket
```
raw-emails-bucket/
├── emails/
│   ├── email-001.eml
│   ├── email-002.eml
│   └── support/
│       ├── support-001.eml
│       └── support-002.eml
```

### Processed Emails Bucket
```
processed-emails-bucket/
├── processed/
│   ├── 2024/01/15/14/
│   │   ├── email-001.eml
│   │   ├── email-001.json
│   │   ├── email-002.eml
│   │   └── email-002.json
│   └── 2024/01/15/15/
│       ├── support-001.eml
│       ├── support-001.json
│       ├── support-002.eml
│       └── support-002.json
```

## Monitoring and Alerting

Example CloudWatch alarm configuration:

```bash
# Lambda function error rate alarm
aws cloudwatch put-metric-alarm \
  --alarm-name "EmailProcessor-ErrorRate" \
  --alarm-description "Email processor Lambda error rate" \
  --metric-name Errors \
  --namespace AWS/Lambda \
  --statistic Sum \
  --period 300 \
  --threshold 1 \
  --comparison-operator GreaterThanOrEqualToThreshold \
  --evaluation-periods 2 \
  --dimensions Name=FunctionName,Value=aloha-email-processor

# S3 bucket size growth alarm  
aws cloudwatch put-metric-alarm \
  --alarm-name "RawEmailsBucket-SizeGrowth" \
  --alarm-description "Raw emails bucket size growth" \
  --metric-name BucketSizeBytes \
  --namespace AWS/S3 \
  --statistic Average \
  --period 86400 \
  --threshold 10737418240 \
  --comparison-operator GreaterThanThreshold \
  --evaluation-periods 1 \
  --dimensions Name=BucketName,Value=your-raw-emails-bucket Name=StorageType,Value=StandardStorage
```

## Local Development

For local testing and development:

```bash
# Set up local AWS credentials
aws configure

# Run tests
cd lambda/email-processor
go test -v

# Build locally
./build-lambda.sh

# Lint code
go vet ./...
go fmt ./...
```

## Production Considerations

1. **Security**: Enable S3 encryption and VPC endpoints
2. **Monitoring**: Set up comprehensive CloudWatch dashboards
3. **Backup**: Configure cross-region replication for S3 buckets
4. **Scaling**: Monitor Lambda concurrency limits
5. **Cost**: Set up S3 lifecycle policies for cost optimization