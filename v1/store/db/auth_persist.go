package db

import (
	"reflect"

	"github.com/bww/go-acl/v1"

	"github.com/bww/go-dbx/v1/persist"
	"github.com/bww/go-dbx/v1/persist/registry"
	"github.com/bww/go-util/uuid"
)

func init() {
	registry.Set(reflect.ValueOf((*acl.Authorization)(nil)).Type(), &authorizationPersister{})
}

type authorizationPersister struct{}

func (p authorizationPersister) StoreRelated(pst persist.Persister, val interface{}) error {
	z := val.(*acl.Authorization)
	ps := NewPolicyStore(pst)

	for i, e := range z.Policies {
		p, err := ps.StorePolicy(e)
		if err != nil {
			return err
		}
		z.Policies[i] = p
	}

	return nil
}

func (p authorizationPersister) StoreReferences(pst persist.Persister, val interface{}) error {
	z := val.(*acl.Authorization)

	err := p.DeleteReferences(pst, val)
	if err != nil {
		return err
	}

	for _, e := range z.Policies {
		_, err := pst.Exec(`INSERT INTO acl_authorization_r_policy (auth_id, policy_id) VALUES ($1, $2)`, z.Id, e.Id())
		if err != nil {
			return err
		}
	}

	return nil
}

func (p authorizationPersister) FetchRelated(pst persist.Persister, val interface{}) error {
	z := val.(*acl.Authorization)

	rows, err := pst.Query(`SELECT policy_id FROM acl_authorization_r_policy WHERE auth_id = $1`, z.Id)
	if err != nil {
		return err
	}

	policies := make([]acl.Policy, 0)
	ps := NewPolicyStore(pst)

	for rows.Next() {
		var id uuid.UUID
		err := rows.Scan(&id)
		if err != nil {
			return err
		}
		p, err := ps.FetchPolicy(id)
		if err != nil {
			return err
		}
		policies = append(policies, p)
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	z.Policies = policies
	return nil
}

func (e authorizationPersister) DeleteRelated(pst persist.Persister, val interface{}) error {
	z := val.(*acl.Authorization)

	ps := NewPolicyStore(pst)
	for _, e := range z.Policies {
		err := ps.DeletePolicy(e)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e authorizationPersister) DeleteReferences(pst persist.Persister, val interface{}) error {
	z := val.(*acl.Authorization)

	_, err := pst.Exec(`DELETE FROM acl_authorization_r_policy WHERE auth_id = $1`, z.Id)
	if err != nil {
		return err
	}

	return nil
}
