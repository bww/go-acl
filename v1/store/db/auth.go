package acl

import (
	"github.com/bww/go-acl/v1"

	"github.com/bww/go-dbx/v1"
	"github.com/bww/go-dbx/v1/entity"
	"github.com/bww/go-dbx/v1/persist"
	"github.com/bww/go-dbx/v1/persist/ident"
	"github.com/bww/go-util/uuid"
)

const authTable = "acl_authorization"

type authorizationPersister struct {
	persist.Persister
}

func NewAuthorizationPersister(cxt dbx.Context) authorizationPersister {
	return authorizationPersister{persist.New(cxt, entity.DefaultFieldMapper(), ident.UUID)}
}

func (p authorizationPersister) With(cxt dbx.Context) authorizationPersister {
	return authorizationPersister{p.Persister.With(cxt)}
}

func (p authorizationPersister) StoreAuthorization(v *acl.Authorization) error {
	return p.Store(authTable, v, nil)
}

func (p authorizationPersister) CycleAuthorizationCredentials(id uuid.UUID, oldKey, oldSecret, newKey, newSecret string) error {
	_, err := p.Persister.Exec(`UPDATE acl_authorization SET key = $1, secret = $2 WHERE id = $3 AND key = $4 AND secret = $5`, newKey, newSecret, id, oldKey, oldSecret)
	return err
}

func (p authorizationPersister) CountAuthorizations() (int, error) {
	return p.Count(`SELECT COUNT(*) FROM acl_authorization`)
}

func (p authorizationPersister) FetchAuthorization(id uuid.UUID) (*acl.Authorization, error) {
	v := &acl.Authorization{}
	err := p.Fetch(authTable, v, id)
	if err != nil {
		return nil, err
	}
	return v, err
}

func (p authorizationPersister) FetchAuthorizationForKeyAndSecret(key, secret string) (*acl.Authorization, error) {
	v := &acl.Authorization{}
	err := p.Select(v, `SELECT {*} FROM acl_authorization WHERE key = $1 AND secret = $2`, key, secret)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (p authorizationPersister) DeleteAuthorization(v *acl.Authorization) error {
	return p.Delete(authTable, v)
}

// func (p authorizationPersister) StoreRelated(v interface{}, opts persist.StoreOptions, cxt dbx.Context) error {
// 	z := v.(*Authorization)
// 	pp := NewPolicyPersister(e.DefaultContext())
// 	for i, e := range z.Policies {
// 		p, err := pp.StorePolicy(e, opts, cxt)
// 		if err != nil {
// 			return err
// 		}
// 		z.Policies[i] = p
// 	}
// 	return nil
// }

// func (e authorizationPersister) StoreReferences(v interface{}, opts persist.StoreOptions, cxt dbx.Context) error {
// 	z := v.(*Authorization)

// 	err := e.DeleteReferences(v, opts, cxt)
// 	if err != nil {
// 		return err
// 	}

// 	cxt = e.Context(cxt)

// 	for _, e := range z.Policies {
// 		_, err := cxt.Exec(`INSERT INTO acl_authorization_r_policy (auth_id, policy_id) VALUES ($1, $2)`, z.Id, e.Id())
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (e authorizationPersister) FetchRelated(v interface{}, opts persist.FetchOptions, cxt dbx.Context) error {
// 	cxt = e.DefaultContext()
// 	z := v.(*Authorization)

// 	rows, err := cxt.Query(`SELECT policy_id FROM acl_authorization_r_policy WHERE auth_id = $1`, z.Id)
// 	if err != nil {
// 		return err
// 	}

// 	policies := make([]Policy, 0)
// 	pp := NewPolicyPersister(cxt)

// 	for rows.Next() {
// 		var id uuid.UUID
// 		err := rows.Scan(&id)
// 		if err != nil {
// 			return err
// 		}
// 		p, err := pp.FetchPolicy(id, opts, cxt)
// 		if err != nil {
// 			return err
// 		}
// 		policies = append(policies, p)
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		return err
// 	}

// 	z.Policies = policies
// 	return nil
// }

// func (e authorizationPersister) DeleteRelated(v interface{}, opts persist.StoreOptions, cxt dbx.Context) error {
// 	cxt = e.Context(cxt)
// 	z := v.(*Authorization)

// 	pp := NewPolicyPersister(e.DefaultContext())
// 	for _, e := range z.Policies {
// 		err := pp.DeletePolicy(e, opts, cxt)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (e authorizationPersister) DeleteReferences(v interface{}, opts persist.StoreOptions, cxt dbx.Context) error {
// 	cxt = e.Context(cxt)
// 	z := v.(*Authorization)
// 	_, err := cxt.Exec(`DELETE FROM acl_authorization_r_policy WHERE auth_id = $1`, z.Id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
