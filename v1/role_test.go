package acl

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalRole(t *testing.T) {
	tests := []struct {
		Role   Role
		Expect string
		Error  error
	}{
		{
			None, `"none"`, nil,
		},
		{
			Owner, `"owner"`, nil,
		},
		{
			Role("never heard of it"), `"never heard of it"`, nil,
		},
	}

	for _, e := range tests {
		d, err := json.Marshal(e.Role)
		if e.Error != nil {
			if assert.NotNil(t, err, fmt.Sprintf("Expected an error; got: %v", string(d))) {
				assert.Equal(t, e.Error, err.(*json.MarshalerError).Err)
			}
		} else {
			if assert.Nil(t, err, fmt.Sprint(err)) {
				assert.Equal(t, e.Expect, string(d))
			}
		}
	}
}

func TestUnmarshalRole(t *testing.T) {
	tests := []struct {
		Input  string
		Expect Role
		Error  error
	}{
		{
			`"none"`, None, nil,
		},
		{
			`"owner"`, Owner, nil,
		},
		{
			`"never heard of it"`, "", errInvalidRole,
		},
	}

	for _, e := range tests {
		var c Role
		err := json.Unmarshal([]byte(e.Input), &c)
		if e.Error != nil {
			_ = assert.NotNil(t, err, e.Input) && assert.Equal(t, e.Error, err)
		} else if assert.Nil(t, err, fmt.Sprint(err)) {
			assert.Equal(t, e.Expect, c)
		}
	}
}

func TestRoleCreation(t *testing.T) {
	tests := []struct {
		Roles  Roles
		Grant  Role
		Expect bool
	}{
		{
			nil, None, false,
		},
		{
			Roles{}, None, false,
		},
		{
			Roles{None}, None, false,
		},
		{
			Roles{Member}, Member, false,
		},
		{
			Roles{Admin}, Member, true,
		},
		{
			Roles{Admin}, Admin, true,
		},
		{
			Roles{Admin}, Owner, false,
		},
		{
			Roles{Member, Admin}, Member, true,
		},
		{
			Roles{Owner}, Owner, true,
		},
	}

	for _, e := range tests {
		v := e.Roles.CanGrant(e.Grant)
		fmt.Printf("--> %v (%v): %v\n", e.Roles, e.Grant, v)
		assert.Equal(t, e.Expect, v)
	}
}
