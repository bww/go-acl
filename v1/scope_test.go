package acl

import (
	"encoding/json"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseScopes(t *testing.T) {
	tests := []struct {
		Input  string
		Expect Scope
		Error  error
	}{
		{
			"a", Scope{nil, "a"}, nil,
		},
		{
			":a", Scope{nil, "a"}, nil,
		},
		{
			"*:a", Scope{Actions{Every}, "a"}, nil,
		},
		{
			"read:a", Scope{Actions{Read}, "a"}, nil,
		},
		{
			"read,write:a", Scope{Actions{Read, Write}, "a"}, nil,
		},
		{
			"read,write,delete:a", Scope{Actions{Read, Write, Delete}, "a"}, nil,
		},
		{
			"read,write,list,delete:a", Scope{Actions{Read, Write, List, Delete}, "a"}, nil,
		},
		{
			"read,write,delete,foobar:a", Scope{}, errInvalidAction,
		},
		{
			"read,write,*:a", Scope{Actions{Every}, "a"}, nil,
		},
		{
			"read,*,write,delete:a", Scope{Actions{Every}, "a"}, nil,
		},
		{
			"read:a/b**C_d@3FG ANYTHING ELSE // whatever you want~~~~", Scope{Actions{Read}, "a/b**C_d@3FG ANYTHING ELSE // whatever you want~~~~"}, nil,
		},
		{
			"read,", Scope{}, errEmptyResource,
		},
		{
			",", Scope{}, errEmptyResource,
		},
		{
			",:", Scope{}, errEmptyResource,
		},
		{
			",:foo", Scope{nil, "foo"}, nil,
		},
		{
			",,,:::foo", Scope{nil, "::foo"}, nil,
		},
	}

	for _, e := range tests {
		s, err := ParseScope(e.Input)
		if e.Error != nil {
			fmt.Println("***", err)
			assert.Equal(t, e.Error, err)
		} else if assert.Nil(t, err, fmt.Sprint(err)) {
			fmt.Println("-->", e.Input, "/", s)
			assert.Equal(t, e.Expect, s)
		}
	}
}

func TestParseScopeSets(t *testing.T) {
	tests := []struct {
		Scopes  Scopes
		Require Scopes
		Expect  bool
	}{
		{
			Scopes{{Actions{Read}, "a"}},
			Scopes{{Actions{Read}, "a"}},
			true,
		},
		{
			Scopes{{Actions{Read}, "a"}},
			Scopes{{Actions{Read}, "a"}, {Actions{Read}, "b"}},
			false,
		},
		{
			Scopes{{Actions{Read}, "a"}},
			Scopes{{Actions{Read}, "a"}, {Actions{Read}, "a"}},
			true,
		},
		{
			Scopes{{Actions{Read}, "a"}},
			Scopes{{Actions{Read}, "b"}, {Actions{Every}, "a"}},
			false,
		},
		{
			Scopes{{Actions{Read, Write}, "a"}},
			Scopes{{Actions{Read}, "a"}, {Actions{Write}, "a"}},
			true,
		},
		{
			Scopes{{Actions{Read}, "a"}, {Actions{Write}, "a"}},
			Scopes{{Actions{Read, Write}, "a"}},
			false,
		},
		{
			Scopes{{Actions{Every}, "a"}},
			Scopes{{Actions{Read, Write, Delete}, "a"}},
			true,
		},
		{
			Scopes{{Actions{Every, Read}, "a"}},
			Scopes{{Actions{Read, Write, Delete}, "a"}},
			true,
		},
		{
			Scopes{{Actions{Read, Write}, "a"}},
			Scopes{{Actions{Every}, "a"}},
			false,
		},
		{
			Scopes{{Actions{Every}, "a"}},
			Scopes{{Actions{Every}, "a"}},
			true,
		},
		{
			Scopes{{Actions{Read}, "a"}, {Actions{Read}, "b"}},
			Scopes{{Actions{Read, Write}, "a"}},
			false,
		},
		{
			Scopes{{Actions{Read}, "a"}, {Actions{Read, Write}, "b"}},
			Scopes{{Actions{Read, Write}, "a"}},
			false,
		},
		{
			Scopes{{Actions{Read, Write}, "a"}},
			Scopes{{Actions{Read, Write}, "a"}},
			true,
		},
		{
			Scopes{{Actions{Every}, "a"}},
			Scopes{{Actions{Read, Write}, "a"}},
			true,
		},
		{
			Scopes{{Actions{Every}, "a"}},
			Scopes{{Actions{}, "a"}},
			false,
		},
		{
			Scopes{{Actions{}, "a"}},
			Scopes{{Actions{Read}, "a"}},
			false,
		},
		{
			Scopes{{Actions{Read}, "b"}},
			Scopes{{Actions{Read}, "a"}},
			false,
		},
		{
			Scopes{{Actions{Read}, "b"}},
			Scopes{{Actions{Read}, "a"}},
			false,
		},
		{
			Scopes{{Actions{Read}, "b"}},
			Scopes{{Actions{Read}, ""}},
			false,
		},
		{
			Scopes{{Actions{Read}, ""}},
			Scopes{{Actions{Read}, "a"}},
			false,
		},
		{
			Scopes{{Actions{Read}, ""}},
			Scopes{{Actions{Read}, ""}},
			false,
		},
		{
			Scopes{{Actions{}, ""}},
			Scopes{{Actions{}, ""}},
			false,
		},
	}

	for _, e := range tests {
		v := e.Scopes.Satisfies(e.Require...)
		fmt.Println("-->", e.Scopes, "/", e.Require)
		assert.Equal(t, e.Expect, v)
	}
}

func TestMergedScopes(t *testing.T) {
	tests := []struct {
		Scopes Scopes
		Expect Scopes
	}{
		{
			Scopes{{Actions{Read}, "a"}, {Actions{Write}, "a"}},
			Scopes{{Actions{Read, Write}, "a"}},
		},
		{
			Scopes{{Actions{Read}, "a"}, {Actions{Write}, "b"}},
			Scopes{{Actions{Read}, "a"}, {Actions{Write}, "b"}},
		},
		{
			Scopes{{Actions{Read}, "a"}, {Actions{Write}, "a"}, {Actions{Every}, "a"}},
			Scopes{{Actions{Every}, "a"}},
		},
		{
			Scopes{{Actions{Read}, "a"}, {Actions{Write}, "a"}, {Actions{Delete}, "a"}},
			Scopes{{Actions{Delete, Read, Write}, "a"}},
		},
		{
			Scopes{{Actions{Read}, "a"}, {Actions{Write}, "a"}, {Actions{Delete}, "a"}, {Actions{Read, Write}, "b"}, {Actions{Delete}, "b"}},
			Scopes{{Actions{Delete, Read, Write}, "a"}, {Actions{Delete, Read, Write}, "b"}},
		},
	}

	for _, e := range tests {
		m := e.Scopes.Merged()
		for _, e := range m {
			sort.Sort(e.Actions)
		}
		sort.Sort(m)
		fmt.Println("-->", e.Scopes, "=", m)
		assert.Equal(t, e.Expect, m)
	}
}

func TestMergeScopes(t *testing.T) {
	tests := []struct {
		A, B   Scopes
		Expect Scopes
	}{
		{
			Scopes{{Actions{Read}, "a"}},
			Scopes{{Actions{Write}, "a"}},
			Scopes{{Actions{Read, Write}, "a"}},
		},
		{
			Scopes{{Actions{Read}, "a"}, {Actions{Write}, "a"}},
			Scopes{{Actions{Write}, "a"}},
			Scopes{{Actions{Read, Write}, "a"}},
		},
		{
			Scopes{{Actions{Read}, "a"}, {Actions{Write}, "b"}},
			Scopes{{Actions{Write}, "a"}},
			Scopes{{Actions{Read, Write}, "a"}, {Actions{Write}, "b"}},
		},
	}

	for _, e := range tests {
		var m Scopes

		// using Merge()
		m = e.A.Merge(e.B)
		for _, e := range m {
			sort.Sort(e.Actions)
		}
		sort.Sort(m)
		fmt.Println("--> (Merge())", e.A, "+", e.B, "=", m)
		assert.Equal(t, e.Expect, m)

		// using Union()
		m = Union(e.A, e.B)
		for _, e := range m {
			sort.Sort(e.Actions)
		}
		sort.Sort(m)
		fmt.Println("--> (Union())", e.A, "+", e.B, "=", m)
		assert.Equal(t, e.Expect, m)
	}
}

func TestMarshalScope(t *testing.T) {
	tests := []struct {
		Scope  Scope
		Expect string
		Error  error
	}{
		{
			Scope{Actions{Read}, "a"}, `"read:a"`, nil,
		},
		{
			Scope{Actions{Every}, "a"}, `"*:a"`, nil,
		},
		{
			Scope{Actions{Read, Write, Delete}, "a"}, `"read,write,delete:a"`, nil,
		},
	}

	for _, e := range tests {
		fmt.Println("-->", e.Scope)
		d, err := json.Marshal(e.Scope)
		if e.Error != nil {
			if assert.NotNil(t, err, fmt.Sprintf("Expected an error; got: %v", string(d))) {
				assert.Equal(t, e.Error, err.(*json.MarshalerError).Err)
			}
			continue
		}
		if assert.Nil(t, err, fmt.Sprint(err)) {
			assert.Equal(t, e.Expect, string(d))
		}
		var s Scope
		err = json.Unmarshal(d, &s)
		if assert.Nil(t, err, fmt.Sprint(err)) {
			assert.Equal(t, e.Scope, s)
		}
	}
}
