package ecs

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/jsii-runtime-go"
)

func (e *ECS) MakeTask() {
	executionRole := awsiam.NewRole(e.scope, jsii.String("ExecRole"), &awsiam.RoleProps{
		RoleName:  jsii.String("esc-exec-role"),
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("ecs-tasks.amazonaws.com"), nil),
		ManagedPolicies: &[]awsiam.IManagedPolicy{
			awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AmazonECSTaskExecutionRolePolicy")),
		},
	})

	statements := []awsiam.PolicyStatement{
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Effect:    awsiam.Effect_ALLOW,
			Actions:   jsii.Strings("s3:*"),
			Resources: jsii.Strings("*"),
		}),
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Effect: awsiam.Effect_ALLOW,
			Actions: jsii.Strings(
				"ssmmessages:CreateControlChannel",
				"ssmmessages:CreateDataChannel",
				"ssmmessages:OpenControlChannel",
				"ssmmessages:OpenDataChannel",
				"logs:CreateLogStream",
				"logs:PutLogEvents",
			),
			Resources: jsii.Strings("*"),
		}),
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Effect:    awsiam.Effect_ALLOW,
			Actions:   jsii.Strings("secretsmanager:GetSecretValue"),
			Resources: jsii.Strings("*"),
		}),
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Effect:    awsiam.Effect_ALLOW,
			Actions:   jsii.Strings("logs:CreateExportTask"),
			Resources: jsii.Strings("*"),
		}),
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Effect:    awsiam.Effect_ALLOW,
			Actions:   jsii.Strings("ecs:DescribeTaskDefinition"),
			Resources: jsii.Strings("*"),
		}),
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Effect: awsiam.Effect_ALLOW,
			Actions: jsii.Strings(
				"events:PutRule",
				"events:PutTargets",
			),
			Resources: jsii.Strings("*"),
		}),
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Effect:    awsiam.Effect_ALLOW,
			Actions:   jsii.Strings("iam:PassRole"),
			Resources: jsii.Strings("*"),
		}),
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Effect:    awsiam.Effect_ALLOW,
			Actions:   jsii.Strings("cloudfront:*"),
			Resources: jsii.Strings("*"),
		}),
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Effect:    awsiam.Effect_ALLOW,
			Actions:   jsii.Strings("ecr:GetAuthorizationToken"),
			Resources: jsii.Strings("*"),
		}),
	}

	policy := awsiam.NewManagedPolicy(e.scope, jsii.String("TaskPolicy"), &awsiam.ManagedPolicyProps{
		ManagedPolicyName: jsii.String("ecs-task-policy"),
		Statements:        &statements,
	})
	taskRole := awsiam.NewRole(e.scope, jsii.String("ecs-task-role"), &awsiam.RoleProps{
		RoleName:        jsii.String("ecs-task-role"),
		AssumedBy:       awsiam.NewServicePrincipal(jsii.String("ecs-tasks.amazonaws.com"), nil),
		ManagedPolicies: &[]awsiam.IManagedPolicy{policy},
	})

	repository := awsecr.Repository_FromRepositoryName(e.scope, jsii.String("redash-repository"), jsii.String("ecr-redash"))
	// rdsSecret := awssecretsmanager.Secret_FromSecretNameV2(e.scope, jsii.String("redash-secret"), jsii.String("RDSSecret"))

	createDbDefinition := awsecs.NewFargateTaskDefinition(e.scope, jsii.String("CreateDBTaskDefinition"), &awsecs.FargateTaskDefinitionProps{
		ExecutionRole:  executionRole,
		TaskRole:       taskRole,
		MemoryLimitMiB: jsii.Number(2048),
		Cpu:            jsii.Number(1024),
	})
	createDbDefinition.AddContainer(jsii.String("CreateDBContainer"), &awsecs.ContainerDefinitionOptions{
		Image:          awsecs.ContainerImage_FromEcrRepository(repository, jsii.String("latest")),
		MemoryLimitMiB: jsii.Number(2048),
		Cpu:            jsii.Number(1024),
		Command: &[]*string{
			jsii.String("create_db"),
		},
		Environment: &map[string]*string{
			"PYTHONUNBUFFERED":     jsii.String(os.Getenv("PYTHONUNBUFFERED")),
			"REDASH_LOG_LEVEL":     jsii.String(os.Getenv("REDASH_LOG_LEVEL")),
			"REDASH_REDIS_URL":     jsii.String(os.Getenv("REDASH_REDIS_URL")),
			"REDASH_DATABASE_URL":  jsii.String(os.Getenv("REDASH_DATABASE_URL")),
			"REDASH_COOKIE_SECRET": jsii.String(os.Getenv("REDASH_COOKIE_SECRET")),
			"REDASH_SECRET_KEY":    jsii.String(os.Getenv("REDASH_SECRET_KEY")),
		},
		Logging: awsecs.LogDrivers_AwsLogs(&awsecs.AwsLogDriverProps{
			StreamPrefix: jsii.String("redash"),
		}),
	})
	e.tasks["createDB"] = createDbDefinition

	workerTaskDefinition := awsecs.NewFargateTaskDefinition(e.scope, jsii.String("WorkerTaskDefinition"), &awsecs.FargateTaskDefinitionProps{
		ExecutionRole:  executionRole,
		TaskRole:       taskRole,
		MemoryLimitMiB: jsii.Number(2048),
		Cpu:            jsii.Number(1024),
	})
	workerTaskDefinition.AddContainer(jsii.String("WorkerContainer"), &awsecs.ContainerDefinitionOptions{
		Image:          awsecs.ContainerImage_FromEcrRepository(repository, jsii.String("latest")),
		MemoryLimitMiB: jsii.Number(2048),
		Cpu:            jsii.Number(1024),
		Command: &[]*string{
			jsii.String("worker"),
		},
		Environment: &map[string]*string{
			"PYTHONUNBUFFERED":              jsii.String(os.Getenv("PYTHONUNBUFFERED")),
			"REDASH_LOG_LEVEL":              jsii.String(os.Getenv("REDASH_LOG_LEVEL")),
			"REDASH_REDIS_URL":              jsii.String(os.Getenv("REDASH_REDIS_URL")),
			"REDASH_DATABASE_URL":           jsii.String(os.Getenv("REDASH_DATABASE_URL")),
			"REDASH_COOKIE_SECRET":          jsii.String(os.Getenv("REDASH_COOKIE_SECRET")),
			"REDASH_SECRET_KEY":             jsii.String(os.Getenv("REDASH_SECRET_KEY")),
			"WORKERS_COUNT":                 jsii.String("4"),
			"QUEUES":                        jsii.String("queries,scheduled_queries,celery"),
			"REDASH_PASSWORD_LOGIN_ENABLED": jsii.String("true"),
			"REDASH_LDAP_LOGIN_ENABLED":     jsii.String("false"),
		},
		Logging: awsecs.LogDrivers_AwsLogs(&awsecs.AwsLogDriverProps{
			StreamPrefix: jsii.String("redash"),
		}),
		Essential: jsii.Bool(true),
	})
	e.tasks["worker"] = workerTaskDefinition

	schedulerTaskDefinition := awsecs.NewFargateTaskDefinition(e.scope, jsii.String("SchedulerTaskDefinition"), &awsecs.FargateTaskDefinitionProps{
		ExecutionRole:  executionRole,
		TaskRole:       taskRole,
		MemoryLimitMiB: jsii.Number(2048),
		Cpu:            jsii.Number(1024),
	})
	schedulerTaskDefinition.AddContainer(jsii.String("SchedulerContainer"), &awsecs.ContainerDefinitionOptions{
		Image:          awsecs.ContainerImage_FromEcrRepository(repository, jsii.String("latest")),
		MemoryLimitMiB: jsii.Number(2048),
		Cpu:            jsii.Number(1024),
		Command: &[]*string{
			jsii.String("scheduler"),
		},
		Environment: &map[string]*string{
			"PYTHONUNBUFFERED":     jsii.String(os.Getenv("PYTHONUNBUFFERED")),
			"REDASH_LOG_LEVEL":     jsii.String(os.Getenv("REDASH_LOG_LEVEL")),
			"REDASH_REDIS_URL":     jsii.String(os.Getenv("REDASH_REDIS_URL")),
			"REDASH_DATABASE_URL":  jsii.String(os.Getenv("REDASH_DATABASE_URL")),
			"REDASH_COOKIE_SECRET": jsii.String(os.Getenv("REDASH_COOKIE_SECRET")),
			"REDASH_SECRET_KEY":    jsii.String(os.Getenv("REDASH_SECRET_KEY")),
		},
		Logging: awsecs.LogDrivers_AwsLogs(&awsecs.AwsLogDriverProps{
			StreamPrefix: jsii.String("redash"),
		}),
		Essential: jsii.Bool(true),
	})
	e.tasks["scheduler"] = schedulerTaskDefinition

	serverTaskDefinition := awsecs.NewFargateTaskDefinition(e.scope, jsii.String("ServerTaskDefinition"), &awsecs.FargateTaskDefinitionProps{
		ExecutionRole:  executionRole,
		TaskRole:       taskRole,
		MemoryLimitMiB: jsii.Number(2048),
		Cpu:            jsii.Number(1024),
	})
	serverTaskDefinition.AddContainer(jsii.String("ServerContainer"), &awsecs.ContainerDefinitionOptions{
		Image:          awsecs.ContainerImage_FromEcrRepository(repository, jsii.String("redash-nginx")),
		MemoryLimitMiB: jsii.Number(2048),
		Cpu:            jsii.Number(1024),
		Command: &[]*string{
			jsii.String("server"),
		},
		Environment: &map[string]*string{
			"PYTHONUNBUFFERED":     jsii.String(os.Getenv("PYTHONUNBUFFERED")),
			"REDASH_LOG_LEVEL":     jsii.String(os.Getenv("REDASH_LOG_LEVEL")),
			"REDASH_REDIS_URL":     jsii.String(os.Getenv("REDASH_REDIS_URL")),
			"REDASH_DATABASE_URL":  jsii.String(os.Getenv("REDASH_DATABASE_URL")),
			"REDASH_COOKIE_SECRET": jsii.String(os.Getenv("REDASH_COOKIE_SECRET")),
			"REDASH_SECRET_KEY":    jsii.String(os.Getenv("REDASH_SECRET_KEY")),
			"REDASH_WEB_WORKER":    jsii.String("4"),
		},
		Logging: awsecs.LogDrivers_AwsLogs(&awsecs.AwsLogDriverProps{
			StreamPrefix: jsii.String("redash"),
		}),
		Essential: jsii.Bool(true),
		PortMappings: &[]*awsecs.PortMapping{
			{
				HostPort:      jsii.Number(5000),
				ContainerPort: jsii.Number(5000),
			},
		},
	})
	e.tasks["server"] = serverTaskDefinition
}
