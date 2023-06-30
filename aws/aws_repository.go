package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/raito-io/cli/base/util/config"
)

func GetAWSConfig(ctx context.Context, configMap *config.ConfigMap) (aws.Config, error) {
	loadOptions := make([]func(*awsconfig.LoadOptions) error, 0)

	profile := configMap.GetStringWithDefault(AwsProfile, "")
	if profile != "" {
		loadOptions = append(loadOptions, awsconfig.WithSharedConfigProfile(profile))
	}

	awsRegion := configMap.GetStringWithDefault(AwsRegion, "")
	if awsRegion != "" {
		loadOptions = append(loadOptions, awsconfig.WithRegion(awsRegion))
	}

	return awsconfig.LoadDefaultConfig(ctx, loadOptions...)
}
