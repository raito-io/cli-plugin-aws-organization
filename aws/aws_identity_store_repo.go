package aws

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/identitystore/types"
	sso "github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	is "github.com/raito-io/cli/base/identity_store"
	"github.com/raito-io/cli/base/util/config"
)

type AwsIdentityStoreRepository struct {
}

func (repo *AwsIdentityStoreRepository) getIdentityStoreClient(ctx context.Context, configMap *config.ConfigMap) (*identitystore.Client, error) {
	cfg, err := GetAWSConfig(ctx, configMap)

	if err != nil {
		return nil, err
	}

	client := identitystore.NewFromConfig(cfg)

	return client, nil
}

func (repo *AwsIdentityStoreRepository) getSSOClient(ctx context.Context, configMap *config.ConfigMap) (*sso.Client, error) {
	cfg, err := GetAWSConfig(ctx, configMap)

	if err != nil {
		return nil, err
	}

	client := sso.NewFromConfig(cfg)

	return client, nil
}

func (repo *AwsIdentityStoreRepository) GetIdentityStores(ctx context.Context, configMap *config.ConfigMap) ([]string, error) {
	ssoClient, err := repo.getSSOClient(ctx, configMap)
	if err != nil {
		return nil, err
	}

	moreObjectsAvailable := true
	var nextToken *string
	identityStores := make([]string, 0)

	for moreObjectsAvailable {
		response, err := ssoClient.ListInstances(ctx, &sso.ListInstancesInput{
			NextToken: nextToken,
		})

		if err != nil {
			return nil, fmt.Errorf("error while listing identity stores: %s", err.Error())
		}

		moreObjectsAvailable = response.NextToken != nil
		nextToken = response.NextToken

		for _, instance := range response.Instances {
			if instance.IdentityStoreId != nil {
				identityStores = append(identityStores, *instance.IdentityStoreId)
			}
		}
	}

	logger.Info(fmt.Sprintf("Found identity stores: %+v", identityStores))

	return identityStores, nil
}

func (repo *AwsIdentityStoreRepository) GetUsers(ctx context.Context, identityStores []string, configMap *config.ConfigMap) ([]is.User, error) {
	client, err := repo.getIdentityStoreClient(ctx, configMap)
	if err != nil {
		return nil, err
	}

	var result []is.User

	for _, identityStoreId := range identityStores {
		isID := identityStoreId
		moreObjectsAvailable := true
		var nextToken *string

		for moreObjectsAvailable {
			input := identitystore.ListUsersInput{
				NextToken:       nextToken,
				IdentityStoreId: &isID,
			}

			response, err := client.ListUsers(ctx, &input)
			if err != nil {
				return nil, fmt.Errorf("error while listing users: %s", err.Error())
			}

			moreObjectsAvailable = response.NextToken != nil
			nextToken = response.NextToken

			for i := range response.Users {
				user := response.Users[i]

				emailAddress := *user.UserName

				for _, email := range user.Emails {
					if email.Primary && email.Value != nil {
						emailAddress = *email.Value
					}
				}

				result = append(result, is.User{
					ExternalId: *user.UserId,
					UserName:   *user.UserName,
					Name:       *user.UserName,
					Email:      emailAddress,
				})
			}
		}
	}

	logger.Info(fmt.Sprintf("%d users have been found", len(result)))

	return result, nil
}

func (repo *AwsIdentityStoreRepository) GetGroups(ctx context.Context, identityStores []string, configMap *config.ConfigMap) ([]is.Group, map[string][]string, error) {
	client, err := repo.getIdentityStoreClient(ctx, configMap)
	if err != nil {
		return nil, nil, err
	}

	var result []is.Group
	parentMap := make(map[string][]string)

	for _, identityStoreId := range identityStores {
		isID := identityStoreId
		moreObjectsAvailable := true
		var nextToken *string

		for moreObjectsAvailable {
			input := identitystore.ListGroupsInput{
				NextToken:       nextToken,
				IdentityStoreId: &isID,
			}

			response, err := client.ListGroups(ctx, &input)
			if err != nil {
				return nil, nil, fmt.Errorf("error while listing groups: %s", err.Error())
			}

			moreObjectsAvailable = response.NextToken != nil
			nextToken = response.NextToken

			for _, group := range response.Groups {
				if err = repo.addGroupToParentMap(ctx, client, group.GroupId, group.IdentityStoreId, parentMap); err != nil {
					return nil, nil, fmt.Errorf("error while fetching group members: %s", err.Error())
				}

				description := ""
				if group.Description != nil {
					description = *group.Description
				}

				result = append(result, is.Group{
					ExternalId:  *group.GroupId,
					Name:        *group.DisplayName,
					DisplayName: *group.DisplayName,
					Description: description,
				})
			}
		}
	}

	logger.Info(fmt.Sprintf("%d groups have been found", len(result)))

	return result, parentMap, nil
}

func (repo *AwsIdentityStoreRepository) addGroupToParentMap(ctx context.Context, client *identitystore.Client, groupId, identityStoreId *string, parentMap map[string][]string) error {
	moreObjectsAvailable := true
	var nextToken *string

	for moreObjectsAvailable {
		input := identitystore.ListGroupMembershipsInput{
			NextToken:       nextToken,
			GroupId:         groupId,
			IdentityStoreId: identityStoreId,
		}

		response, err := client.ListGroupMemberships(ctx, &input)
		if err != nil {
			return err
		}

		moreObjectsAvailable = response.NextToken != nil
		nextToken = response.NextToken

		for _, memberShip := range response.GroupMemberships {
			if um, f := memberShip.MemberId.(*types.MemberIdMemberUserId); f {
				member := um.Value

				parents := parentMap[member]
				parentMap[member] = append(parents, *groupId)
			}
		}
	}

	return nil
}
