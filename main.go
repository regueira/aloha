package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3notifications"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AlohaEmailProcessorStackProps struct {
	awscdk.StackProps
}

type AlohaEmailProcessorStack struct {
	awscdk.Stack
}

func NewAlohaEmailProcessorStack(scope constructs.Construct, id string, props *AlohaEmailProcessorStackProps) AlohaEmailProcessorStack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Create S3 bucket for raw SES emails
	rawEmailsBucket := awss3.NewBucket(stack, jsii.String("RawEmailsBucket"), &awss3.BucketProps{
		Versioned:     jsii.Bool(true),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		LifecycleRules: &[]*awss3.LifecycleRule{
			{
				Id:      jsii.String("DeleteOldEmails"),
				Enabled: jsii.Bool(true),
				Expiration: awscdk.Duration_Days(jsii.Number(90)),
			},
		},
	})

	// Create S3 bucket for processed emails
	processedEmailsBucket := awss3.NewBucket(stack, jsii.String("ProcessedEmailsBucket"), &awss3.BucketProps{
		Versioned:     jsii.Bool(true),
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		LifecycleRules: &[]*awss3.LifecycleRule{
			{
				Id:      jsii.String("DeleteOldProcessedEmails"),
				Enabled: jsii.Bool(true),
				Expiration: awscdk.Duration_Days(jsii.Number(365)),
			},
		},
	})

	// Create IAM role for Lambda function
	lambdaRole := awsiam.NewRole(stack, jsii.String("EmailProcessorLambdaRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), nil),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AWSLambdaBasicExecutionRole")),
		},
		InlinePolicies: &map[string]awsiam.PolicyDocument{
			"S3Access": awsiam.NewPolicyDocument(&awsiam.PolicyDocumentProps{
				Statements: &[]awsiam.PolicyStatement{
					awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
						Effect: awsiam.Effect_ALLOW,
						Actions: &[]*string{
							jsii.String("s3:GetObject"),
							jsii.String("s3:PutObject"),
							jsii.String("s3:DeleteObject"),
						},
						Resources: &[]*string{
							rawEmailsBucket.ArnForObjects(jsii.String("*")),
							processedEmailsBucket.ArnForObjects(jsii.String("*")),
						},
					}),
					awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
						Effect: awsiam.Effect_ALLOW,
						Actions: &[]*string{
							jsii.String("s3:ListBucket"),
						},
						Resources: &[]*string{
							rawEmailsBucket.BucketArn(),
							processedEmailsBucket.BucketArn(),
						},
					}),
				},
			}),
		},
	})

	// Create Lambda function for email processing
	emailProcessorLambda := awslambda.NewFunction(stack, jsii.String("EmailProcessorFunction"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Handler:      jsii.String("bootstrap"),
		Code:         awslambda.Code_FromAsset(jsii.String("lambda/email-processor"), nil),
		Role:         lambdaRole,
		FunctionName: jsii.String("aloha-email-processor"),
		Description:  jsii.String("Processes raw emails from SES and moves them to processed bucket"),
		Timeout:      awscdk.Duration_Minutes(jsii.Number(5)),
		Environment: &map[string]*string{
			"RAW_BUCKET":       rawEmailsBucket.BucketName(),
			"PROCESSED_BUCKET": processedEmailsBucket.BucketName(),
		},
	})

	// Grant S3 bucket permissions to Lambda
	rawEmailsBucket.GrantRead(emailProcessorLambda, jsii.String("*"))
	processedEmailsBucket.GrantReadWrite(emailProcessorLambda, jsii.String("*"))

	// Add S3 event trigger for Lambda
	rawEmailsBucket.AddEventNotification(
		awss3.EventType_OBJECT_CREATED,
		awss3notifications.NewLambdaDestination(emailProcessorLambda),
		&awss3.NotificationKeyFilter{
			Prefix: jsii.String("emails/"),
			Suffix: jsii.String(".eml"),
		},
	)

	// Note: SES configuration needs to be done manually through AWS Console or CLI
	// as the SES receipt rules API in CDK requires additional setup

	// Output the bucket names
	awscdk.NewCfnOutput(stack, jsii.String("RawEmailsBucketName"), &awscdk.CfnOutputProps{
		Value:       rawEmailsBucket.BucketName(),
		Description: jsii.String("Name of the S3 bucket for raw emails"),
	})

	awscdk.NewCfnOutput(stack, jsii.String("ProcessedEmailsBucketName"), &awscdk.CfnOutputProps{
		Value:       processedEmailsBucket.BucketName(),
		Description: jsii.String("Name of the S3 bucket for processed emails"),
	})

	awscdk.NewCfnOutput(stack, jsii.String("EmailProcessorLambdaName"), &awscdk.CfnOutputProps{
		Value:       emailProcessorLambda.FunctionName(),
		Description: jsii.String("Name of the email processor Lambda function"),
	})

	return AlohaEmailProcessorStack{
		Stack: stack,
	}
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewAlohaEmailProcessorStack(app, "AlohaEmailProcessorStack", &AlohaEmailProcessorStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	account := os.Getenv("CDK_DEFAULT_ACCOUNT")
	region := os.Getenv("CDK_DEFAULT_REGION")

	if len(account) == 0 {
		account = "123456789012" // Default placeholder account
	}
	if len(region) == 0 {
		region = "us-east-1" // Default region
	}

	return &awscdk.Environment{
		Account: jsii.String(account),
		Region:  jsii.String(region),
	}
}