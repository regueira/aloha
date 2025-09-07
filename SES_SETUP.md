# SES Setup Guide

This guide explains how to configure Amazon Simple Email Service (SES) to work with the Aloha Email Processor.

## Prerequisites

1. Deploy the CDK stack first: `cdk deploy`
2. Note the S3 bucket name from the stack output (`RawEmailsBucketName`)
3. Have a verified domain in SES

## Step 1: Verify Your Domain

1. Go to the AWS SES Console
2. Navigate to "Verified identities"
3. Click "Create identity"
4. Choose "Domain" and enter your domain name
5. Follow the DNS verification steps

## Step 2: Create Receipt Rule Set

1. In the SES Console, go to "Email receiving"
2. Click "Create rule set"
3. Name it "aloha-email-rules"
4. Make it the active rule set

## Step 3: Create Receipt Rule

1. Click "Create rule" in your rule set
2. Configure the rule:
   - **Recipients**: Add the email addresses you want to process (e.g., `inbox@yourdomain.com`)
   - **Actions**: Add an S3 action
     - **S3 bucket**: Select the raw emails bucket from your CDK deployment
     - **Object key prefix**: `emails/`
     - **Object key suffix**: `.eml`

## Step 4: Configure MX Records

Add MX records to your domain's DNS:
- **Record Type**: MX
- **Name**: Your domain (e.g., `yourdomain.com`)
- **Value**: The appropriate SES inbound mail endpoint for your region:
  - `us-east-1`: `inbound-smtp.us-east-1.amazonaws.com` (Priority: 10)
  - `us-west-2`: `inbound-smtp.us-west-2.amazonaws.com` (Priority: 10)
  - `eu-west-1`: `inbound-smtp.eu-west-1.amazonaws.com` (Priority: 10)

## Step 5: Test Email Processing

1. Send a test email to your configured address
2. Check the raw emails S3 bucket for the incoming email
3. Check the processed emails S3 bucket for the processed email and metadata
4. Monitor CloudWatch Logs for the Lambda function

## CLI Commands

You can also set up SES using the AWS CLI:

```bash
# Create rule set
aws ses create-receipt-rule-set --rule-set-name aloha-email-rules

# Create receipt rule
aws ses create-receipt-rule \
  --rule-set-name aloha-email-rules \
  --rule '{
    "Name": "save-to-s3",
    "Recipients": ["inbox@yourdomain.com"],
    "Enabled": true,
    "Actions": [
      {
        "S3Action": {
          "BucketName": "YOUR_RAW_EMAILS_BUCKET_NAME",
          "ObjectKeyPrefix": "emails/"
        }
      }
    ]
  }'

# Set as active rule set
aws ses set-active-receipt-rule-set --rule-set-name aloha-email-rules
```

## Troubleshooting

### Email Not Received
- Verify your domain is confirmed in SES
- Check MX records are properly configured
- Ensure the receipt rule set is active
- Check the rule recipients match the email address used

### Lambda Not Triggering
- Verify S3 event notifications are configured on the raw emails bucket
- Check the object key prefix/suffix filters
- Review Lambda function logs in CloudWatch

### Permissions Issues
- Ensure the Lambda execution role has proper S3 permissions
- Verify SES can write to the S3 bucket

## Security Considerations

- Consider enabling S3 bucket encryption
- Review IAM policies for least privilege access
- Enable CloudTrail for audit logging
- Consider implementing email content scanning for security

## Monitoring

Set up CloudWatch alarms for:
- Lambda function errors
- S3 bucket size growth
- Failed email processing metrics