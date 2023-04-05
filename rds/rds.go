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
	sg := awsec2.NewSecurityGroup(r.scope, jsii.String("redash-db-securityGroup"), &awsec2.SecurityGroupProps{
		Vpc: r.vpc,
	})

	sg.AddIngressRule(
		awsec2.Peer_Ipv4(jsii.String("0.0.0.0/0")),
		awsec2.Port_Tcp(jsii.Number(5432)),
		jsii.String("Allow PostgreSQL from anywhere"),
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
		InstanceType: awsec2.InstanceType_Of(awsec2.InstanceClass_M5, awsec2.InstanceSize_XLARGE),
		Credentials:  awsrds.Credentials_FromSecret(rdsSecret, nil),
		DatabaseName: jsii.String("redash_postgresql"),
		Vpc:          r.vpc,
		SubnetGroup: awsrds.NewSubnetGroup(r.scope, jsii.String("redash-db-subnetGroup"), &awsrds.SubnetGroupProps{
			Vpc:             r.vpc,
			Description:     jsii.String("db subnet group"),
			SubnetGroupName: jsii.String("redash-db-subnetGroup"),
		}),
		PubliclyAccessible:      jsii.Bool(false),
		SecurityGroups:          &[]awsec2.ISecurityGroup{sg},
		InstanceIdentifier:      jsii.String("redash-db"),
		StorageType:             awsrds.StorageType_IO1,
		AllocatedStorage:        jsii.Number(400),
		Iops:                    jsii.Number(3000),
		MultiAz:                 jsii.Bool(false),
		AutoMinorVersionUpgrade: jsii.Bool(true),
		DeleteAutomatedBackups:  jsii.Bool(false),
	})
}
