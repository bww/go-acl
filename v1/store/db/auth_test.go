package db

import (
	"fmt"
	"testing"

	"github.com/bww/go-acl/v1"

	"github.com/bww/go-dbx/v1"
	"github.com/bww/go-dbx/v1/test"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizationCRUD(t *testing.T) {
	pa := NewAuthorizationStore(test.DB())
	now := dbx.Now()
	var err error

	a1 := &acl.Authorization{
		Key:    "every",
		Secret: "ABC123",
		Policies: []acl.Policy{
			acl.AllowResource(acl.Read, "/companies/*"),
			acl.AllowResource(acl.Write, "/partners/*"),
		},
		Created: now,
	}

	err = pa.StoreAuthorization(a1)
	assert.Nil(t, err, fmt.Sprint(err))

	a2, err := pa.FetchAuthorization(a1.Id)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, a1, a2)
	}

	a2.Policies = []acl.Policy{
		acl.AllowResource(acl.Read, "/companies/THIS_ONE"),
		acl.AllowResource(acl.Write, "/partners/THAT_ONE"),
	}
	err = pa.StoreAuthorization(a2)
	assert.Nil(t, err, fmt.Sprint(err))

	a3, err := pa.FetchAuthorization(a2.Id)
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, a2, a3)
	}

	err = pa.DeleteAuthorization(a3)
	assert.Nil(t, err, fmt.Sprint(err))

	_, err = pa.FetchAuthorization(a3.Id)
	if assert.NotNil(t, err, "Expected a not-found error") {
		assert.Equal(t, dbx.ErrNotFound, err, "Expected a not-found error")
	}

}
