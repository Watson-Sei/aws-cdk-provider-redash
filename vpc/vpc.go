package vpc

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type VPC struct {
	scope constructs.Construct
	vpc   awsec2.Vpc
}

func NewVPC(scope constructs.Construct) *VPC {
	return &VPC{
		scope: scope,
	}
}

func (v *VPC) Make() {
	v.vpc = awsec2.NewVpc(v.scope, jsii.String("redash-vpc"), &awsec2.VpcProps{
		Cidr:   jsii.String("10.0.0.0/21"),
		MaxAzs: jsii.Number(2),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				Name:       jsii.String("redash-subnet-public"),
				SubnetType: awsec2.SubnetType_PUBLIC,
			},
			{
				Name:       jsii.String("redash-subnet-private"),
				SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
			},
		},
	})

	v.vpc.AddInterfaceEndpoint(jsii.String("EndpointEcr"), &awsec2.InterfaceVpcEndpointOptions{
		Service: awsec2.InterfaceVpcEndpointAwsService_ECR(),
	})
	v.vpc.AddInterfaceEndpoint(jsii.String("EndpointEcrDocker"), &awsec2.InterfaceVpcEndpointOptions{
		Service: awsec2.InterfaceVpcEndpointAwsService_ECR_DOCKER(),
	})
	v.vpc.AddInterfaceEndpoint(jsii.String("EndpointEcs"), &awsec2.InterfaceVpcEndpointOptions{
		Service: awsec2.InterfaceVpcEndpointAwsService_ECS(),
	})
	v.vpc.AddInterfaceEndpoint(jsii.String("EndpointRds"), &awsec2.InterfaceVpcEndpointOptions{
		Service: awsec2.InterfaceVpcEndpointAwsService_RDS(),
	})
	v.vpc.AddInterfaceEndpoint(jsii.String("EndpointElastiCache"), &awsec2.InterfaceVpcEndpointOptions{
		Service: awsec2.InterfaceVpcEndpointAwsService_ELASTICACHE(),
	})
}

func (v *VPC) Get() awsec2.Vpc {
	return v.vpc
}
