package apikey

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/bww/go-acl/v1"
	"github.com/bww/go-acl/v1/provider"

	"github.com/bww/go-dbx/v1"
	"github.com/stretchr/testify/assert"
)

type credentials struct {
	key, secret string
}

type dataSource struct {
	auth *acl.Authorization
}

func (s dataSource) AuthorizationForKeyAndSecret(key, secret string) (*acl.Authorization, error) {
	if s.auth.Key == key && s.auth.Secret == secret {
		return s.auth, nil
	} else {
		return nil, dbx.ErrNotFound
	}
}

func TestAPIKeyProvider(t *testing.T) {
	now := time.Now()
	a1 := &acl.Authorization{
		Key:    "test",
		Secret: "ABC123",
		Policies: []acl.Policy{
			acl.AllowResource(acl.Read, "/companies/*"),
			acl.DenyResource(acl.Every, "/partners/*/secret"),
			acl.AllowResource(acl.Write, "/partners/*"),
		},
		Active:  true,
		Created: now,
	}

	a := APIKeyProvider{dataSource{a1}}
	var r *http.Request

	r = request(t, "GET", "/companies/THIS_COMPANY", grant("test", "ABC123"), nil)
	check(t, a, r, a1, nil)

	r = request(t, "GET", "/companies/THIS_COMPANY/employees", grant("test", "ABC123"), nil)
	check(t, a, r, a1, nil)

	r = request(t, "GET", "/companies/THIS_COMPANY", grant("test", "XYZ987"), nil)
	check(t, a, r, a1, provider.ErrUnauthorized)

	r = request(t, "GET", "/companies/THIS_COMPANY", nil, nil)
	check(t, a, r, a1, provider.ErrUnauthorized)

	r = request(t, "PUT", "/companies/THIS_COMPANY", grant("test", "ABC123"), nil)
	check(t, a, r, a1, provider.ErrForbidden)

	r = request(t, "GET", "/elsewhere", grant("test", "ABC123"), nil)
	check(t, a, r, a1, provider.ErrForbidden)

	r = request(t, "GET", "/partners/THIS_PARTNER/secret", grant("test", "ABC123"), nil)
	check(t, a, r, a1, provider.ErrForbidden)

}

func grant(k, s string) *credentials {
	return &credentials{k, s}
}

func request(t *testing.T, m, p string, c *credentials, b []byte) *http.Request {
	r, err := http.NewRequest(m, p, bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	if c != nil {
		r.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.key, c.secret)))))
	}
	return r
}

func check(t *testing.T, a APIKeyProvider, r *http.Request, g *acl.Authorization, f error) {
	fmt.Printf("--> %s %s (%v)\n", r.Method, r.URL.Path, g.Policies)
	err := a.Validate(r)
	if f != nil {
		assert.Equal(t, f, err, fmt.Sprintf("(%s %s)", r.Method, r.URL.Path))
	} else {
		assert.Nil(t, err, fmt.Sprintf("%v (%s %s)", err, r.Method, r.URL.Path))
	}
}
