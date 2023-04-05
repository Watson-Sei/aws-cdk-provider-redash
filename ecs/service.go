package ecs

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/jsii-runtime-go"
)

func (e *ECS) MakeService() {
	// awsecs.NewFargateService(e.scope, jsii.String("redash-createdb-service"), &awsecs.FargateServiceProps{
	// 	Cluster:        e.cluster,
	// 	TaskDefinition: e.tasks["createDB"],
	// 	DesiredCount:   jsii.Number(1),
	// })

	sg := awsec2.NewSecurityGroup(e.scope, jsii.String("redash-server-sg"), &awsec2.SecurityGroupProps{
		SecurityGroupName: jsii.String("redash-server-sg"),
		Vpc:               e.vpc,
		AllowAllOutbound:  jsii.Bool(true),
	})
	sg.AddIngressRule(
		awsec2.Peer_AnyIpv4(),
		awsec2.Port_AllTcp(),
		jsii.String("Allow all inbound TCP"),
		jsii.Bool(false),
	)

	awsecs.NewFargateService(e.scope, jsii.String("redash-worker-service"), &awsecs.FargateServiceProps{
		Cluster:        e.cluster,
		TaskDefinition: e.tasks["worker"],
		DesiredCount:   jsii.Number(1),
	})

	awsecs.NewFargateService(e.scope, jsii.String("redash-scheduler-service"), &awsecs.FargateServiceProps{
		Cluster:        e.cluster,
		TaskDefinition: e.tasks["scheduler"],
		DesiredCount:   jsii.Number(1),
	})

	awsecs.NewFargateService(e.scope, jsii.String("redash-server-service"), &awsecs.FargateServiceProps{
		Cluster:              e.cluster,
		DesiredCount:         jsii.Number(1),
		ServiceName:          new(string),
		TaskDefinition:       e.tasks["server"],
		AssignPublicIp:       jsii.Bool(true),
		PlatformVersion:      awsecs.FargatePlatformVersion_VERSION1_4,
		SecurityGroups:       &[]awsec2.ISecurityGroup{sg},
		VpcSubnets:           &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PUBLIC},
		MaxHealthyPercent:    jsii.Number(200),
		MinHealthyPercent:    jsii.Number(100),
		EnableExecuteCommand: jsii.Bool(true),
	})
}
