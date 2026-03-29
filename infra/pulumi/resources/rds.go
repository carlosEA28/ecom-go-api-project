package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type RDSOutput struct {
	Database     *rds.Instance
	Endpoint     pulumi.StringOutput
	Port         pulumi.IntOutput
	DatabaseName pulumi.StringOutput
	Username     pulumi.StringOutput
}

func CreateRDS(ctx *pulumi.Context, rdsSubnetID pulumi.Output, rdsSecurityGroupID pulumi.StringOutput) (*RDSOutput, error) {
	cfg := config.New(ctx, "")
	environment := cfg.Get("environment")
	if environment == "" {
		environment = "dev"
	}

	// Get configuration from Pulumi config - passwords must come from secrets
	dbPassword := cfg.RequireSecret("dbPassword")

	dbName := cfg.Get("dbName")
	if dbName == "" {
		dbName = "ecom"
	}

	dbUsername := cfg.Get("dbUsername")
	if dbUsername == "" {
		dbUsername = "postgres"
	}

	dbInstanceClass := cfg.Get("dbInstanceClass")
	if dbInstanceClass == "" {
		if environment == "prod" {
			dbInstanceClass = "db.r5.large"
		} else {
			dbInstanceClass = "db.t3.micro"
		}
	}

	dbAllocatedStorage := 20
	if environment == "prod" {
		dbAllocatedStorage = 100
	}
	if storage := cfg.GetInt("dbAllocatedStorage"); storage > 0 {
		dbAllocatedStorage = storage
	}

	// Create DB Subnet Group for RDS
	dbSubnetGroup, err := rds.NewSubnetGroup(ctx, "rds-subnet-group", &rds.SubnetGroupArgs{
		SubnetIds: pulumi.StringArray{
			rdsSubnetID.ApplyT(func(id interface{}) string {
				return id.(string)
			}).(pulumi.StringOutput),
		},
		Tags: pulumi.StringMap{
			"Name":        pulumi.String("ecom-db-subnet-group"),
			"Environment": pulumi.String(environment),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create RDS subnet group: %w", err)
	}

	// Determine skip final snapshot based on environment
	skipFinalSnapshot := environment == "dev"

	// Create RDS PostgreSQL instance
	database, err := rds.NewInstance(ctx, "ecom-database", &rds.InstanceArgs{
		AllocatedStorage:  pulumi.Int(dbAllocatedStorage),
		Engine:            pulumi.String("postgres"),
		EngineVersion:     pulumi.String("15"),
		InstanceClass:     pulumi.String(dbInstanceClass),
		DbName:            pulumi.String(dbName),
		Username:          pulumi.String(dbUsername),
		Password:          dbPassword,
		SkipFinalSnapshot: pulumi.Bool(skipFinalSnapshot),
		DbSubnetGroupName: dbSubnetGroup.Name,
		VpcSecurityGroupIds: pulumi.StringArray{
			rdsSecurityGroupID,
		},
		PubliclyAccessible:               pulumi.Bool(false),
		BackupRetentionPeriod:            pulumi.Int(30),                    // Keep 30 days of backups
		MultiAz:                          pulumi.Bool(environment != "dev"), // HA for non-dev
		StorageEncrypted:                 pulumi.Bool(true),                 // Enable encryption
		IamDatabaseAuthenticationEnabled: pulumi.Bool(true),                 // Use IAM auth
		DeletionProtection:               pulumi.Bool(environment == "prod"),
		AutoMinorVersionUpgrade:          pulumi.Bool(true),
		Tags: pulumi.StringMap{
			"Name":        pulumi.String("ecom-database"),
			"Environment": pulumi.String(environment),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create RDS instance: %w", err)
	}

	return &RDSOutput{
		Database:     database,
		Endpoint:     database.Endpoint,
		Port:         database.Port,
		DatabaseName: database.DbName,
		Username:     database.Username,
	}, nil
}
