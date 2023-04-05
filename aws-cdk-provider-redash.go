package main

import (
	"aws-cdk-provider-redash/ecs"
	"aws-cdk-provider-redash/elasticache"
	"aws-cdk-provider-redash/rds"
	"aws-cdk-provider-redash/vpc"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/joho/godotenv"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AwsCdkProviderRedashStackProps struct {
	awscdk.StackProps
}

func NewAwsCdkProviderRedashStack(scope constructs.Construct, id string, props *AwsCdkProviderRedashStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	v := vpc.NewVPC(stack)
	v.Make()

	r := rds.NewRDS(stack, v.Get())
	r.Make()

	er := elasticache.NewElastiCache(stack, v.Get())
	er.Make()

	ecs := ecs.NewECS(stack, v.Get(), er.GetClusterAddress())
	ecs.Make()

	return stack
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewAwsCdkProviderRedashStack(app, "AwsCdkProviderRedashStack", &AwsCdkProviderRedashStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	// return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
