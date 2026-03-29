package resources

import (
	"github.com/pulumi/pulumi-awsx/sdk/v3/go/awsx/ecr"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ECROutput struct {
	Repository *ecr.Repository
	Image      *ecr.Image
	ImageURI   pulumi.StringOutput
}

func CreateECR(ctx *pulumi.Context) (*ECROutput, error) {
	repository, err := ecr.NewRepository(ctx, "ecom-api-repo", &ecr.RepositoryArgs{
		ForceDelete: pulumi.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	image, err := ecr.NewImage(ctx, "ecom-api-image", &ecr.ImageArgs{
		RepositoryUrl: repository.Url,
		Context:       pulumi.String(".."),
		Dockerfile:    pulumi.String("../Dockerfile"),
		Platform:      pulumi.String("linux/amd64"),
	})
	if err != nil {
		return nil, err
	}

	return &ECROutput{
		Repository: repository,
		Image:      image,
		ImageURI:   image.ImageUri,
	}, nil
}

func GetECRRepositoryUrl(ecrOutput *ECROutput) pulumi.StringOutput {
	return ecrOutput.Repository.Url
}
