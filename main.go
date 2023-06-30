package main

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/raito-io/cli-plugin-aws-organization/aws"
	"github.com/raito-io/cli/base"
	"github.com/raito-io/cli/base/info"
	"github.com/raito-io/cli/base/util/plugin"
	"github.com/raito-io/cli/base/wrappers"
)

var version = "0.0.0"

var logger hclog.Logger

func main() {
	logger = base.Logger()
	logger.SetLevel(hclog.Debug)

	err := base.RegisterPlugins(
		wrappers.IdentityStoreSync(aws.NewIdentityStoreSyncer()), &info.InfoImpl{
			Info: &plugin.PluginInfo{
				Name:    "AWS Organization",
				Version: plugin.ParseVersion(version),
				Parameters: []*plugin.ParameterInfo{
					{Name: aws.AwsProfile, Description: "The AWS SDK profile to use for connecting to the AWS account where the organization is managed. When not specified, the default profile is used (or what is defined in the AWS_PROFILE environment variable).", Mandatory: false},
					{Name: aws.AwsRegion, Description: "The AWS region to use. If not set, the default region from the AWS SDK is used.", Mandatory: false},
				},
			},
		})

	if err != nil {
		logger.Error(fmt.Sprintf("error while registering plugins: %s", err.Error()))
	}
}
