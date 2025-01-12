package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3notifications"
	"github.com/lpernett/godotenv"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type SummarizerAiStackProps struct {
	awscdk.StackProps
}

func NewSummarizerAiStack(scope constructs.Construct, id string, props *SummarizerAiStackProps) awscdk.Stack {
	godotenv.Load()
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	fileBucket := awss3.NewBucket(stack, jsii.String("my-bucket"), &awss3.BucketProps{
		BucketName:        jsii.String(os.Getenv("BUCKET_NAME")),
		Versioned:         jsii.Bool(true),
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
	})

	apiGateway := awsapigateway.NewRestApi(stack, jsii.String("apiGateway"), &awsapigateway.RestApiProps{
		CloudWatchRole: jsii.Bool(true),
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowHeaders: jsii.Strings("Content-Type"),
			AllowMethods: jsii.Strings("GET", "POST", "PUT", "DELETE", "OPTIONS"),
			AllowOrigins: jsii.Strings("*"),
		},
		DeployOptions: &awsapigateway.StageOptions{
			LoggingLevel: awsapigateway.MethodLoggingLevel_INFO,
		},
	})
	fileUploadLambda := awslambda.NewFunction(stack, jsii.String("fileUploadLambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("lambda/function.zip"), nil),
		Handler: jsii.String("main"),
		Environment: &map[string]*string{
			"BUCKET_NAME": jsii.String(os.Getenv("BUCKET_NAME")),
		},
	})
	summarizerLambda := awslambda.NewFunction(stack, jsii.String("summarizer"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PYTHON_3_12(),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("python-service"), nil),
		Handler: jsii.String("summarize.summarizer_handler"),
		Environment: &map[string]*string{
			"GEMINI_API_KEY": jsii.String(os.Getenv("GEMINI_API_KEY")),
			"BUCKET_NAME":    jsii.String(os.Getenv("BUCKET_NAME")),
		},
	})
	fileBucket.GrantReadWrite(fileUploadLambda, nil)
	fileBucket.GrantReadWrite(summarizerLambda, nil)
	apiV1 := apiGateway.Root().AddResource(jsii.String("api"), nil).AddResource(jsii.String("v1"), nil)
	fileUploadIntegration := awsapigateway.NewLambdaIntegration(fileUploadLambda, nil)
	summarizerLambdaIntegration := awsapigateway.NewLambdaIntegration(summarizerLambda, nil)
	summarizer := apiV1.AddResource(jsii.String("get-summary"), nil)
	summarizer.AddMethod(jsii.String("GET"), summarizerLambdaIntegration, nil)

	fileBucket.AddEventNotification(awss3.EventType_OBJECT_CREATED, awss3notifications.NewLambdaDestination(summarizerLambda))
	fileUpload := apiV1.AddResource(jsii.String("upload"), nil)
	fileUpload.AddMethod(jsii.String("POST"), fileUploadIntegration, nil)
	return stack
}


func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewSummarizerAiStack(app, "SummarizerAIStack", &SummarizerAiStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return nil
}
