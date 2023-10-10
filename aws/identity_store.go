package aws

import (
	"context"

	"github.com/raito-io/cli/base/util/config"
	"github.com/raito-io/cli/base/wrappers"

	is "github.com/raito-io/cli/base/identity_store"
)

//go:generate go run github.com/vektra/mockery/v2 --name=identityStoreRepository --with-expecter --inpackage
type identityStoreRepository interface {
	GetIdentityStores(ctx context.Context, configMap *config.ConfigMap) ([]string, error)
	GetUsers(ctx context.Context, identityStores []string, configMap *config.ConfigMap) ([]is.User, error)
	GetGroups(ctx context.Context, identityStores []string, configMap *config.ConfigMap) ([]is.Group, map[string][]string, error)
}

type IdentityStoreSyncer struct {
	repoProvider func(configMap *config.ConfigMap) identityStoreRepository
}

func NewIdentityStoreSyncer() *IdentityStoreSyncer {
	return &IdentityStoreSyncer{repoProvider: newRepoProvider}
}

func (s *IdentityStoreSyncer) GetIdentityStoreMetaData(_ context.Context, _ *config.ConfigMap) (*is.MetaData, error) {
	logger.Debug("Returning meta data for AWS identity store")

	return &is.MetaData{
		Type:        "aws-organization",
		CanBeLinked: true,
		CanBeMaster: true,
	}, nil
}

func newRepoProvider(configMap *config.ConfigMap) identityStoreRepository {
	return &AwsIdentityStoreRepository{}
}

func (s *IdentityStoreSyncer) SyncIdentityStore(ctx context.Context, identityHandler wrappers.IdentityStoreIdentityHandler, configMap *config.ConfigMap) error {
	identityStores, err := s.repoProvider(configMap).GetIdentityStores(ctx, configMap)
	if err != nil {
		return err
	}

	groups, parentMap, err := s.repoProvider(configMap).GetGroups(ctx, identityStores, configMap)
	if err != nil {
		return err
	}

	for _, g := range groups {
		err = identityHandler.AddGroups(&is.Group{
			ExternalId:  g.ExternalId,
			Name:        g.Name,
			DisplayName: g.Name,
			// ParentGroupExternalIds: nil, // AWS IAM doesn't support group hierarchies
		})
		if err != nil {
			return err
		}
	}

	// get users
	users, err := s.repoProvider(configMap).GetUsers(ctx, identityStores, configMap)
	if err != nil {
		return err
	}

	for _, u := range users {
		// TODO: figure out what the best mapping is here. No email for AWS users.
		err = identityHandler.AddUsers(&is.User{
			ExternalId:       u.ExternalId,
			UserName:         u.Name,
			Email:            u.Email,
			Name:             u.Name,
			Tags:             u.Tags,
			GroupExternalIds: parentMap[u.ExternalId],
		})
		if err != nil {
			return err
		}
	}

	return nil
}
