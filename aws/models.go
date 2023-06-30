package aws

import (
	"github.com/raito-io/cli/base/tag"
)

type GroupEntity struct {
	ARN        string
	ExternalId string
	Name       string
	Members    []string
}

type UserEntity struct {
	ARN        string
	ExternalId string
	Name       string
	Email      string //not natively used in AWS
	Tags       []*tag.Tag
}
