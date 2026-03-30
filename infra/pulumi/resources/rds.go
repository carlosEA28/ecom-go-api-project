package resources

import (
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type RDSOutput struct {
	Database     *rds.Instance
	Endpoint     pulumi.StringOutput
	Port         pulumi.IntOutput
	DatabaseName pulumi.StringOutput
	Username     pulumi.StringOutput
}

func CreateRDS(ctx *pulumi.Context, privateSubnetIDs pulumi.StringArrayOutput, rdsSecurityGroupID pulumi.StringOutput) (*RDSOutput, error) {
	// Criar DB Subnet Group para o RDS
	dbSubnetGroup, err := rds.NewSubnetGroup(ctx, "rds-subnet-group", &rds.SubnetGroupArgs{
		SubnetIds: privateSubnetIDs,
	})
	if err != nil {
		return nil, err
	}

	// Criar instância RDS PostgreSQL
	database, err := rds.NewInstance(ctx, "ecom-database", &rds.InstanceArgs{
		AllocatedStorage:  pulumi.Int(20),
		Engine:            pulumi.String("postgres"),
		EngineVersion:     pulumi.String("15"),
		InstanceClass:     pulumi.String("db.t3.micro"),
		DbName:            pulumi.String("ecom"),
		Username:          pulumi.String("postgres"),
		Password:          pulumi.String("postgres123!"),
		SkipFinalSnapshot: pulumi.Bool(true),
		DbSubnetGroupName: dbSubnetGroup.Name,
		VpcSecurityGroupIds: pulumi.StringArray{
			rdsSecurityGroupID,
		},
		PubliclyAccessible: pulumi.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	return &RDSOutput{
		Database:     database,
		Endpoint:     database.Endpoint,
		Port:         database.Port,
		DatabaseName: database.DbName,
		Username:     database.Username,
	}, nil
}
