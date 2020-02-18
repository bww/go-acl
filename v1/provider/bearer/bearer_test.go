package bearer

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/bww/go-acl/v1"
	"github.com/bww/go-acl/v1/provider"

	"github.com/bww/go-util/v1/crypto"
	"github.com/stretchr/testify/assert"
)

func TestBearerProvider(t *testing.T) {
	b := BearerProvider{crypto.GenerateKey("This is the secret key", "Some salt", crypto.SHA1)}
	var r *http.Request

	grants := []acl.Policy{
		acl.AllowResource(acl.Read, "/companies/*"),
		acl.DenyResource(acl.Every, "/partners/*/secret"),
		acl.AllowResource(acl.Write, "/partners/*"),
	}

	r = request(t, "GET", "/companies/THIS_COMPANY", nil)
	check(t, methodBearer, b, r, grants, nil)

	r = request(t, "GET", "/companies/THIS_COMPANY/employees", nil)
	check(t, methodBearer, b, r, grants, nil)

	r = request(t, "GET", "/companies/THIS_COMPANY", nil)
	check(t, methodBearer, b, r, nil, provider.ErrUnauthorized)

	r = request(t, "PUT", "/companies/THIS_COMPANY", nil)
	check(t, methodBearer, b, r, grants, provider.ErrForbidden)

	r = request(t, "GET", "/elsewhere", nil)
	check(t, methodBearer, b, r, grants, provider.ErrForbidden)

	r = request(t, "GET", "/partners/THIS_PARTNER/secret", nil)
	check(t, methodBearer, b, r, grants, provider.ErrForbidden)

}

func request(t *testing.T, m, p string, b []byte) *http.Request {
	r, err := http.NewRequest(m, p, bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	return r
}

func check(t *testing.T, m string, b BearerProvider, r *http.Request, g []acl.Policy, f error) {
	fmt.Printf("--> %s %s (%v)\n", r.Method, r.URL.Path, g)
	var err error

	if g != nil {
		enc, sig, err := crypto.SignMessage(b.key, crypto.SHA256, g)
		assert.Nil(t, err, fmt.Sprintf("%v", err))
		tok := sig + sep + enc
		r.Header.Set("Authorization", fmt.Sprintf("%s %s", m, tok))
	}

	err = b.Validate(r)
	if f != nil {
		assert.Equal(t, f, err, fmt.Sprintf("(%s %s)", r.Method, r.URL.Path))
	} else {
		assert.Nil(t, err, fmt.Sprintf("%v (%s %s)", err, r.Method, r.URL.Path))
	}
}
