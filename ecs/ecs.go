package ecs

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	"github.com/aws/constructs-go/constructs/v10"
)

type ECS struct {
	scope constructs.Construct
	vpc   awsec2.Vpc

	cluster             awsecs.Cluster
	tasks               map[string]awsecs.TaskDefinition
	redisClusterAddress *string
	alb                 awselasticloadbalancingv2.ApplicationLoadBalancer
}

func NewECS(scope constructs.Construct, vpc awsec2.Vpc, redisClusterAddress *string) *ECS {
	return &ECS{
		scope:               scope,
		vpc:                 vpc,
		tasks:               make(map[string]awsecs.TaskDefinition),
		redisClusterAddress: redisClusterAddress,
	}
}

func (e *ECS) Make() {
	e.MakeCluster()
	e.MakeTask()
	e.MakeAlb()
	e.MakeService()
}
