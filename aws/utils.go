package aws

import (
	"github.com/raito-io/cli/base/tag"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/hashicorp/go-hclog"
	"github.com/raito-io/cli/base"
)

var logger hclog.Logger

func init() {
	logger = base.Logger()
}

func getTags(input []types.Tag) []*tag.Tag {
	tags := make([]*tag.Tag, 0)

	for _, t := range input {
		if t.Key != nil && t.Value != nil {
			tags = append(tags, &tag.Tag{
				Key:    *t.Key,
				Value:  *t.Value,
				Source: "AWS",
			})
		}
	}

	return tags
}

func getEmailAddressFromTags(tags []*tag.Tag, defaultValue string) string {
	for _, t := range tags {
		if strings.Contains(t.Key, "email") {
			return t.Value
		}
	}

	return defaultValue
}

func findTagValue(tags []*tag.Tag, key string) *string {
	for _, t := range tags {
		if t.Key == key {
			return &t.Value
		}
	}

	return nil
}
