package db

import (
	"github.com/bww/go-acl/v1"

	"github.com/bww/go-dbx/v1"
	"github.com/bww/go-dbx/v1/entity"
	"github.com/bww/go-dbx/v1/persist"
	"github.com/bww/go-dbx/v1/persist/ident"
	"github.com/bww/go-dbx/v1/persist/registry"
	"github.com/bww/go-util/uuid"
)

type authorizationStore struct {
	persist.Persister
}

func NewAuthorizationStore(cxt dbx.Context) authorizationStore {
	return authorizationStore{persist.New(cxt, entity.DefaultFieldMapper(), registry.DefaultRegistry(), ident.UUID)}
}

func (s authorizationStore) With(cxt dbx.Context) authorizationStore {
	return authorizationStore{s.Persister.With(cxt)}
}

func (s authorizationStore) StoreAuthorization(v *acl.Authorization) error {
	return s.Store(authTable, v, nil)
}

func (s authorizationStore) CycleAuthorizationCredentials(id uuid.UUID, oldKey, oldSecret, newKey, newSecret string) error {
	_, err := s.Exec(`UPDATE `+authTable+` SET key = $1, secret = $2 WHERE id = $3 AND key = $4 AND secret = $5`, newKey, newSecret, id, oldKey, oldSecret)
	return err
}

func (s authorizationStore) CountAuthorizations() (int, error) {
	return s.Count(`SELECT COUNT(*) FROM ` + authTable)
}

func (s authorizationStore) FetchAuthorization(id uuid.UUID) (*acl.Authorization, error) {
	v := &acl.Authorization{}
	err := s.Fetch(authTable, v, id)
	if err != nil {
		return nil, err
	}
	return v, err
}

func (s authorizationStore) FetchAuthorizationForKeyAndSecret(key, secret string) (*acl.Authorization, error) {
	v := &acl.Authorization{}
	err := s.Select(v, `SELECT {*} FROM `+authTable+` WHERE key = $1 AND secret = $2`, key, secret)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (s authorizationStore) DeleteAuthorization(v *acl.Authorization) error {
	return s.Delete(authTable, v)
}
