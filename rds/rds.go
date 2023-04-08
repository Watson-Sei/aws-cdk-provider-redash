package rds

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssecretsmanager"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type RDS struct {
	scope constructs.Construct
	vpc   awsec2.Vpc
}

func NewRDS(scope constructs.Construct, vpc awsec2.Vpc) *RDS {

	return &RDS{
		scope: scope,
		vpc:   vpc,
	}
}

func (r *RDS) Make() {
	subnetSelection := awsec2.SubnetSelection{
		SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
	}

	dbSubnetGroup := awsrds.NewSubnetGroup(r.scope, jsii.String("dbSubnetGroup"), &awsrds.SubnetGroupProps{
		SubnetGroupName: jsii.String("redash"),
		VpcSubnets:      &subnetSelection,
		Vpc:             r.vpc,
		Description:     jsii.String("for redash"),
	})

	securityGroup := awsec2.NewSecurityGroup(r.scope, jsii.String("redashSecurityGroup"), &awsec2.SecurityGroupProps{
		Vpc:               r.vpc,
		SecurityGroupName: jsii.String("rds-postgres-redash"),
		Description:       jsii.String("for redash"),
		AllowAllOutbound:  jsii.Bool(false),
	})

	securityGroup.AddIngressRule(
		awsec2.Peer_AnyIpv4(),
		awsec2.Port_Tcp(jsii.Number(5432)),
		jsii.String("Allow Redash devices HTTP access"),
		jsii.Bool(false),
	)
	securityGroup.AddEgressRule(
		awsec2.Peer_AnyIpv4(),
		awsec2.Port_AllTcp(),
		jsii.String("Allow Redash devices HTTP access"),
		jsii.Bool(false),
	)

	rdsSecret := awssecretsmanager.NewSecret(r.scope, jsii.String("redash-db-secret"), &awssecretsmanager.SecretProps{
		SecretName: jsii.String("RDSSecret"),
		GenerateSecretString: &awssecretsmanager.SecretStringGenerator{
			SecretStringTemplate: jsii.String(`{"username": "postgres"}`),
			GenerateStringKey:    jsii.String("password"),
			ExcludePunctuation:   jsii.Bool(true),
		},
	})

	awsrds.NewDatabaseInstance(r.scope, jsii.String("redash-db"), &awsrds.DatabaseInstanceProps{
		Engine: awsrds.DatabaseInstanceEngine_Postgres(&awsrds.PostgresInstanceEngineProps{
			Version: awsrds.PostgresEngineVersion_VER_14_1(),
		}),
		InstanceType:            awsec2.InstanceType_Of(awsec2.InstanceClass_M5, awsec2.InstanceSize_XLARGE),
		Credentials:             awsrds.Credentials_FromSecret(rdsSecret, nil),
		DatabaseName:            jsii.String("redash_postgresql"),
		Vpc:                     r.vpc,
		SubnetGroup:             dbSubnetGroup,
		PubliclyAccessible:      jsii.Bool(false),
		SecurityGroups:          &[]awsec2.ISecurityGroup{securityGroup},
		InstanceIdentifier:      jsii.String("redash-db"),
		StorageType:             awsrds.StorageType_IO1,
		AllocatedStorage:        jsii.Number(400),
		Iops:                    jsii.Number(3000),
		MultiAz:                 jsii.Bool(false),
		AutoMinorVersionUpgrade: jsii.Bool(true),
		DeleteAutomatedBackups:  jsii.Bool(false),
	})
}
