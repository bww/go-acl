package db

import (
	"fmt"
	"testing"

	"github.com/bww/go-acl/v1"

	"github.com/bww/go-dbx/v1"
	"github.com/bww/go-dbx/v1/test"
	"github.com/stretchr/testify/assert"
)

func TestPolicyCRUD(t *testing.T) {
	ps := NewPolicyStore(test.DB())

	pX := acl.AllowResource(acl.Read, "/companies/*")
	p1, err := ps.StorePolicy(pX)
	assert.Nil(t, err, fmt.Sprint(err))

	p2, err := ps.FetchPolicy(p1.Id())
	if assert.Nil(t, err, fmt.Sprint(err)) {
		assert.Equal(t, p1, p2)
	}

	err = ps.DeletePolicy(p2)
	assert.Nil(t, err, fmt.Sprint(err))

	_, err = ps.FetchPolicy(p2.Id())
	if assert.NotNil(t, err, "Expected a not-found error") {
		assert.Equal(t, dbx.ErrNotFound, err, "Expected a not-found error")
	}
}
