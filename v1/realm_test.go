package acl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalDomain(t *testing.T) {
	tests := []struct {
		Realm  Realm
		Expect string
		Error  error
	}{
		{
			nil, "", nil,
		},
		{
			Realm{{Type: "a", Name: "hello"}}, "a:hello", nil,
		},
		{
			Realm{{Type: "a", Name: "hello"}, {Type: "b", Name: "another"}}, "a:hello/b:another", nil,
		},
		{
			Realm{{Type: "a", Name: "///"}, {Type: "b", Name: "???"}}, "a:%2F%2F%2F/b:%3F%3F%3F", nil,
		},
		{
			Realm{{Type: "wk", Name: "00000000000000000000"}, {Type: "pj", Name: "11111111111111111111"}}, "wk:00000000000000000000/pj:11111111111111111111", nil,
		},
	}
	for _, e := range tests {
		d, err := e.Realm.MarshalText()
		if e.Error != nil {
			fmt.Printf("%s -> %v\n", e.Realm, err)
			assert.ErrorIs(t, err, e.Error)
		} else if assert.NoError(t, err) {
			fmt.Printf("%s -> %v\n", e.Realm, string(d))
			assert.Equal(t, e.Expect, string(d))
		}
	}
}

func TestUnmarshalDomain(t *testing.T) {
	tests := []struct {
		Input  string
		Expect Realm
		Error  error
	}{
		{
			"", nil, nil,
		},
		{
			"a:hello", Realm{{Type: "a", Name: "hello"}}, nil,
		},
		{
			"a:hello/b:another", Realm{{Type: "a", Name: "hello"}, {Type: "b", Name: "another"}}, nil,
		},
		{
			"a:%2F%2F%2F/b:%3F%3F%3F", Realm{{Type: "a", Name: "///"}, {Type: "b", Name: "???"}}, nil,
		},
		{
			"wk:00000000000000000000/pj:11111111111111111111", Realm{{Type: "wk", Name: "00000000000000000000"}, {Type: "pj", Name: "11111111111111111111"}}, nil,
		},
		{
			"no/component/delimiter", nil, errInvalidRealm,
		},
		{
			"invalid:%%%encoding", nil, errInvalidRealm,
		},
	}
	for _, e := range tests {
		var d Realm
		err := d.UnmarshalText([]byte(e.Input))
		if e.Error != nil {
			fmt.Printf("%s -> %v\n", e.Input, err)
			assert.ErrorIs(t, err, e.Error)
		} else if assert.NoError(t, err) {
			fmt.Printf("%s -> %v\n", e.Input, d)
			assert.Equal(t, e.Expect, d)
		}
	}
}
