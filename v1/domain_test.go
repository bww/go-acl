package acl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalDomain(t *testing.T) {
	tests := []struct {
		Domain Domain
		Expect string
		Error  error
	}{
		{
			nil, "", nil,
		},
		{
			Domain{{Type: "a", Name: "hello"}}, "a:hello", nil,
		},
		{
			Domain{{Type: "a", Name: "hello"}, {Type: "b", Name: "another"}}, "a:hello/b:another", nil,
		},
		{
			Domain{{Type: "a", Name: "///"}, {Type: "b", Name: "???"}}, "a:%2F%2F%2F/b:%3F%3F%3F", nil,
		},
		{
			Domain{{Type: "wk", Name: "00000000000000000000"}, {Type: "pj", Name: "11111111111111111111"}}, "wk:00000000000000000000/pj:11111111111111111111", nil,
		},
	}
	for _, e := range tests {
		d, err := e.Domain.MarshalText()
		if e.Error != nil {
			fmt.Printf("%s -> %v\n", e.Domain, err)
			assert.ErrorIs(t, err, e.Error)
		} else if assert.NoError(t, err) {
			fmt.Printf("%s -> %v\n", e.Domain, string(d))
			assert.Equal(t, e.Expect, string(d))
		}
	}
}

func TestUnmarshalDomain(t *testing.T) {
	tests := []struct {
		Input  string
		Expect Domain
		Error  error
	}{
		{
			"", nil, nil,
		},
		{
			"a:hello", Domain{{Type: "a", Name: "hello"}}, nil,
		},
		{
			"a:hello/b:another", Domain{{Type: "a", Name: "hello"}, {Type: "b", Name: "another"}}, nil,
		},
		{
			"a:%2F%2F%2F/b:%3F%3F%3F", Domain{{Type: "a", Name: "///"}, {Type: "b", Name: "???"}}, nil,
		},
		{
			"wk:00000000000000000000/pj:11111111111111111111", Domain{{Type: "wk", Name: "00000000000000000000"}, {Type: "pj", Name: "11111111111111111111"}}, nil,
		},
		{
			"no/component/delimiter", nil, errInvalidDomain,
		},
		{
			"invalid:%%%encoding", nil, errInvalidDomain,
		},
	}
	for _, e := range tests {
		var d Domain
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
