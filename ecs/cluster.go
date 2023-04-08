package ecs

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/jsii-runtime-go"
)

func (e *ECS) MakeCluster() {
	e.cluster = awsecs.NewCluster(e.scope, jsii.String("redash-cluster"), &awsecs.ClusterProps{
		Vpc:                            e.vpc,
		ClusterName:                    jsii.String("redash-cluster"),
		EnableFargateCapacityProviders: jsii.Bool(true),
	})
}
