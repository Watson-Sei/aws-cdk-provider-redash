package elasticache

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awselasticache"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ElastiCache struct {
	scope constructs.Construct
	vpc   awsec2.Vpc

	cluster awselasticache.CfnReplicationGroup
}

func NewElastiCache(scope constructs.Construct, vpc awsec2.Vpc) *ElastiCache {
	return &ElastiCache{
		scope: scope,
		vpc:   vpc,
	}
}

func (e *ElastiCache) Make() {
	// TODO
	sg := awsec2.NewSecurityGroup(e.scope, jsii.String("redash-redis-securityGroup"), &awsec2.SecurityGroupProps{
		Vpc: e.vpc,
	})

	sg.AddIngressRule(
		awsec2.Peer_Ipv4(jsii.String("0.0.0.0/0")),
		awsec2.Port_Tcp(jsii.Number(6379)),
		jsii.String("Allow Redis from anywhere"),
		jsii.Bool(false),
	)

	privateSubnetsIds := e.vpc.PrivateSubnets()
	var subnetIds []*string
	for _, subnet := range *privateSubnetsIds {
		subnetIds = append(subnetIds, subnet.SubnetId())
	}

	subnetGroup := awselasticache.NewCfnSubnetGroup(e.scope, jsii.String("redash-redis-subnetGroup"), &awselasticache.CfnSubnetGroupProps{
		Description:          jsii.String("Redash Redis subnet group"),
		SubnetIds:            &subnetIds,
		CacheSubnetGroupName: jsii.String("redash-redis-subnetGroup"),
		Tags:                 &[]*awscdk.CfnTag{},
	})

	e.cluster = awselasticache.NewCfnReplicationGroup(e.scope, jsii.String("redash-redis"), &awselasticache.CfnReplicationGroupProps{
		ReplicationGroupId:          jsii.String("redash-redis"),
		ReplicationGroupDescription: jsii.String("Redash Redis cluster"),
		CacheNodeType:               jsii.String("cache.r6g.large"),
		Engine:                      jsii.String("redis"),
		EngineVersion:               jsii.String("7.0"),
		Port:                        jsii.Number(6379),
		NumCacheClusters:            jsii.Number(2),
		SecurityGroupIds:            &[]*string{sg.SecurityGroupId()},
		AtRestEncryptionEnabled:     jsii.Bool(true),
		TransitEncryptionEnabled:    jsii.Bool(true),
		AutomaticFailoverEnabled:    jsii.Bool(true),
		CacheParameterGroupName:     jsii.String("default.redis7"),
		CacheSubnetGroupName:        subnetGroup.CacheSubnetGroupName(),
	})

	awscdk.NewCfnOutput(e.scope, jsii.String("redash-redis-endpoint"), &awscdk.CfnOutputProps{
		Value:      e.cluster.AttrPrimaryEndPointAddress(),
		ExportName: jsii.String("redash-redis-endpoint"),
	})
}

func (e *ElastiCache) GetClusterAddress() *string {
	return e.cluster.AttrPrimaryEndPointAddress()
}
