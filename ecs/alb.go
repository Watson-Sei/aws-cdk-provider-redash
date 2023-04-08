package ecs

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	"github.com/aws/jsii-runtime-go"
)

func (e *ECS) MakeAlb() {
	securityGroup := awsec2.NewSecurityGroup(e.scope, jsii.String("AlbSecurityGroup"), &awsec2.SecurityGroupProps{
		Vpc:               e.vpc,
		SecurityGroupName: jsii.String("ALB Security Group"),
		Description:       jsii.String("Security group for Redash ALB"),
	})

	securityGroup.AddIngressRule(awsec2.Peer_Ipv4(jsii.String("10.0.0.0/8")), awsec2.Port_Tcp(jsii.Number(80)), jsii.String("Allow Redash devices HTTP access"), jsii.Bool(false))

	e.alb = awselasticloadbalancingv2.NewApplicationLoadBalancer(e.scope, jsii.String("redash-alb"), &awselasticloadbalancingv2.ApplicationLoadBalancerProps{
		LoadBalancerName: jsii.String("redash-alb"),
		InternetFacing:   jsii.Bool(true),
		SecurityGroup:    securityGroup,
		Vpc:              e.vpc,
		IpAddressType:    awselasticloadbalancingv2.IpAddressType_IPV4,
	})
}
