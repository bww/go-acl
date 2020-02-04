package db

import (
	"reflect"
	"time"

	"github.com/bww/go-acl/v1"

	"github.com/bww/go-dbx/v1"
	"github.com/bww/go-dbx/v1/entity"
	"github.com/bww/go-dbx/v1/persist"
	"github.com/bww/go-dbx/v1/persist/ident"
	"github.com/bww/go-dbx/v1/persist/registry"
	"github.com/bww/go-util/uuid"
)

type policyStore struct {
	persist.Persister
}

func NewPolicyStore(cxt dbx.Context) policyStore {
	return policyStore{persist.New(cxt, entity.DefaultFieldMapper(), registry.DefaultRegistry(), ident.UUID)}
}

func (s policyStore) StorePolicy(v acl.Policy) (acl.Policy, error) {
	id, t, d, err := acl.MarshalPolicy(v)
	if err != nil {
		return nil, err
	}
	if !ident.IsZero(id) {
		return nil, dbx.ErrImmutable
	}

	c := &acl.PersistentPolicy{
		id, t, d, time.Now(),
	}

	err = s.Store(policyTable, c, nil)
	if err != nil {
		return nil, err
	}

	return v.WithId(c.Id), nil
}

func (s policyStore) CountPolicies() (int, error) {
	return s.Count(`SELECT COUNT(*) FROM ` + policyTable)
}

func (s policyStore) FetchPolicy(id uuid.UUID) (acl.Policy, error) {
	c := &acl.PersistentPolicy{}
	err := s.Fetch(policyTable, c, id)
	if err != nil {
		return nil, err
	}
	return acl.UnmarshalPolicy(c.Id, c.Type, c.Data)
}

func (s policyStore) DeletePolicy(v acl.Policy) error {
	return s.DeleteWithID(policyTable, reflect.ValueOf((*acl.PersistentPolicy)(nil)).Type(), v.Id())
}
