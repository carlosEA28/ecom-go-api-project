package resources

import (
	awsecr "github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ecr"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ECROutput struct {
	Repository *awsecr.Repository
	ImageURI   pulumi.StringOutput
}

func CreateECR(ctx *pulumi.Context) (*ECROutput, error) {
	repository, err := awsecr.NewRepository(ctx, "ecom-api-repo", &awsecr.RepositoryArgs{
		ForceDelete: pulumi.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	// Use the repository URL as the image URI (will be tagged with 'latest' or git sha in CI/CD)
	imageUri := repository.RepositoryUrl.ApplyT(func(url interface{}) string {
		return url.(string) + ":latest"
	}).(pulumi.StringOutput)

	return &ECROutput{
		Repository: repository,
		ImageURI:   imageUri,
	}, nil
}

func GetECRRepositoryUrl(ecrOutput *ECROutput) pulumi.StringOutput {
	return ecrOutput.Repository.RepositoryUrl
}
