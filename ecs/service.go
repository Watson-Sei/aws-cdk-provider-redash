package ecs

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	"github.com/aws/jsii-runtime-go"
)

func (e *ECS) MakeService() {
	sg := awsec2.NewSecurityGroup(e.scope, jsii.String("RedashServerServiceSecurityGroup"), &awsec2.SecurityGroupProps{
		Vpc: e.vpc,
	})
	sg.AddIngressRule(
		awsec2.Peer_AnyIpv4(),
		awsec2.Port_Tcp(jsii.Number(80)),
		jsii.String("Allow HTTP access"),
		jsii.Bool(false),
	)

	awsecs.NewFargateService(e.scope, jsii.String("RedashWorkerService"), &awsecs.FargateServiceProps{
		Cluster:              e.cluster,
		DesiredCount:         jsii.Number(1),
		TaskDefinition:       e.tasks["worker"],
		AssignPublicIp:       jsii.Bool(true),
		EnableExecuteCommand: jsii.Bool(true),
	})

	serverService := awsecs.NewFargateService(e.scope, jsii.String("RedashServerService"), &awsecs.FargateServiceProps{
		Cluster:                e.cluster,
		DesiredCount:           jsii.Number(1),
		TaskDefinition:         e.tasks["server"],
		AssignPublicIp:         jsii.Bool(true),
		SecurityGroups:         &[]awsec2.ISecurityGroup{sg},
		EnableExecuteCommand:   jsii.Bool(true),
		MaxHealthyPercent:      jsii.Number(100),
		MinHealthyPercent:      jsii.Number(50),
		HealthCheckGracePeriod: awscdk.Duration_Minutes(jsii.Number(3)),
	})

	targetGroup := awselasticloadbalancingv2.NewApplicationTargetGroup(e.scope, jsii.String("RedashServerTargetGroup"), &awselasticloadbalancingv2.ApplicationTargetGroupProps{
		Vpc:             e.vpc,
		TargetGroupName: jsii.String("redash-alb-target-group"),
		Port:            jsii.Number(5000),
		Targets:         &[]awselasticloadbalancingv2.IApplicationLoadBalancerTarget{serverService},
		Protocol:        awselasticloadbalancingv2.ApplicationProtocol_HTTP,
		HealthCheck: &awselasticloadbalancingv2.HealthCheck{
			Path:                  jsii.String("/ping"),
			Protocol:              awselasticloadbalancingv2.Protocol_HTTP,
			HealthyHttpCodes:      jsii.String("200"),
			HealthyThresholdCount: jsii.Number(2),
			Interval:              awscdk.Duration_Seconds(jsii.Number(30)),
			Timeout:               awscdk.Duration_Seconds(jsii.Number(15)),
		},
		DeregistrationDelay: awscdk.Duration_Seconds(jsii.Number(10)),
	})

	awselasticloadbalancingv2.NewApplicationListener(e.scope, jsii.String("RedashServerListener"), &awselasticloadbalancingv2.ApplicationListenerProps{
		LoadBalancer:        e.alb,
		Port:                jsii.Number(80),
		Protocol:            awselasticloadbalancingv2.ApplicationProtocol_HTTP,
		DefaultTargetGroups: &[]awselasticloadbalancingv2.IApplicationTargetGroup{targetGroup},
		Open:                jsii.Bool(false),
	})
}
