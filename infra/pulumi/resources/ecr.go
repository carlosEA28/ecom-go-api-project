package resources

import (
	"fmt"

	awsecr "github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ecr"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

type ECROutput struct {
	Repository *awsecr.Repository
	ImageURI   pulumi.StringOutput
}

func CreateECR(ctx *pulumi.Context) (*ECROutput, error) {
	cfg := config.New(ctx, "")
	environment := cfg.Get("environment")

	// ForceDelete should be true only for dev (false for staging/prod to prevent data loss)
	forceDelete := environment == "dev"

	repository, err := awsecr.NewRepository(ctx, "ecom-api-repo", &awsecr.RepositoryArgs{
		// Only force delete in dev environment to prevent accidental data loss in prod
		ForceDelete: pulumi.Bool(forceDelete),

		// Enable image scanning on push for security vulnerability detection
		ImageScanningConfiguration: &awsecr.RepositoryImageScanningConfigurationArgs{
			ScanOnPush: pulumi.Bool(true),
		},

		// Enable image tag mutability in production/staging to prevent accidental overwrites
		ImageTagMutability: pulumi.String(func() string {
			if environment == "prod" || environment == "staging" {
				return "IMMUTABLE"
			}
			return "MUTABLE"
		}()),

		Tags: pulumi.StringMap{
			"Name":        pulumi.String("ecom-api-repo"),
			"Environment": pulumi.String(environment),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ECR repository: %w", err)
	}

	// TODO: Add lifecycle policy to clean up old images using AWS CLI or manual configuration
	// Lifecycle policy example:
	// {
	//   "rules": [
	//     {
	//       "rulePriority": 1,
	//       "description": "Keep last 10 tagged images",
	//       "selection": {
	//         "tagStatus": "tagged",
	//         "tagPrefixList": ["v", "sha-"],
	//         "countType": "imageCountMoreThan",
	//         "countNumber": 10
	//       },
	//       "action": { "type": "expire" }
	//     }
	//   ]
	// }

	// Use the repository URL as the image URI (will be tagged with 'latest' or git sha in CI/CD)
	imageUri := repository.RepositoryUrl.ApplyT(func(url interface{}) string {
		return url.(string) + ":latest"
	}).(pulumi.StringOutput)

	ctx.Log.Info(fmt.Sprintf("ECR repository created with ForceDelete=%v, ImageScanning=enabled for %s environment", forceDelete, environment), nil)

	return &ECROutput{
		Repository: repository,
		ImageURI:   imageUri,
	}, nil
}

func GetECRRepositoryUrl(ecrOutput *ECROutput) pulumi.StringOutput {
	return ecrOutput.Repository.RepositoryUrl
}
