package acl

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test policies
func TestPolicy(t *testing.T) {
	var r *http.Request
	var p Policy

	r = request(t, "GET", "/companies/ABC123/employees", nil)
	p = AllowResource(Read, "/")
	checkPolicy(t, r, p, Pass, nil)

	r = request(t, "GET", "/companies/ABC123/employees", nil)
	p = AllowResource(Every, "/")
	checkPolicy(t, r, p, Pass, nil)

	r = request(t, "GET", "/companies/ABC123/employees", nil)
	p = AllowResource(Read, "/*")
	checkPolicy(t, r, p, Allow, nil)

	r = request(t, "GET", "/companies/ABC123/employees", nil)
	p = AllowResource(Read, "/*/")
	checkPolicy(t, r, p, Allow, nil)

	r = request(t, "GET", "/companies/ABC123/employees", nil)
	p = AllowResource(Read, "*")
	checkPolicy(t, r, p, Allow, nil)

	r = request(t, "POST", "/companies/ABC123/employees", nil)
	p = AllowResource(Read, "*")
	checkPolicy(t, r, p, Pass, nil)

	r = request(t, "POST", "/companies/ABC123/employees", nil)
	p = AllowResource(Write, "*")
	checkPolicy(t, r, p, Allow, nil)

	r = request(t, "POST", "/companies/ABC123/employees", nil)
	p = AllowResource(Write, "/companies/*/") // empty final component is ignored
	checkPolicy(t, r, p, Allow, nil)

	r = request(t, "POST", "/companies/ABC123/employees", nil)
	p = AllowResource(Write, "/companies/*")
	checkPolicy(t, r, p, Allow, nil)

	r = request(t, "POST", "/companies/ABC123/employees", nil)
	p = AllowResource(Write, "/companies/*/employees")
	checkPolicy(t, r, p, Allow, nil)

	r = request(t, "POST", "/companies/ABC123/employees", nil)
	p = AllowResource(Write, "/companies/*/locations")
	checkPolicy(t, r, p, Pass, nil)

	r = request(t, "POST", "/companies/ABC123/employees", nil)
	p = AllowResource(Write, "/companies/*/employees/too/specific")
	checkPolicy(t, r, p, Pass, nil)

	r = request(t, "POST", "/companies/ABC123/employees", nil)
	p = AllowResource(Write, "/companies/ABC123/employees")
	checkPolicy(t, r, p, Allow, nil)

	r = request(t, "GET", "/companies/ABC123/employees", nil)
	p = ResourcePolicy{policy{}, ActionSet{Read, Write}, PathSet{Path("/companies/*/employees")}, Allow}
	checkPolicy(t, r, p, Allow, nil)

	r = request(t, "POST", "/companies/ABC123/employees", nil)
	p = ResourcePolicy{policy{}, ActionSet{Read, Write}, PathSet{Path("/companies/*/employees")}, Allow}
	checkPolicy(t, r, p, Allow, nil)

	r = request(t, "DELETE", "/companies/ABC123/employees", nil)
	p = ResourcePolicy{policy{}, ActionSet{Read, Write}, PathSet{Path("/companies/*/employees")}, Allow}
	checkPolicy(t, r, p, Pass, nil)

	r = request(t, "DELETE", "/companies/ABC123/employees", nil)
	p = ResourcePolicy{policy{}, ActionSet{Read, Write, Delete}, PathSet{Path("/companies/*/employees")}, Allow}
	checkPolicy(t, r, p, Allow, nil)

	r = request(t, "UNSUPPORTED", "/companies/ABC123/employees", nil)
	p = ResourcePolicy{policy{}, ActionSet{Read, Write}, PathSet{Path("/companies/*/employees")}, Allow}
	checkPolicy(t, r, p, Pass, ErrMethodNotSupported)

	p = ResourcePolicy{policy{}, ActionSet{Every}, PathSet{Path("/companies/*/employees")}, Allow}
	r = request(t, "GET", "/companies/ABC123/employees", nil)
	checkPolicy(t, r, p, Allow, nil)
	r = request(t, "PUT", "/companies/ABC123/employees", nil)
	checkPolicy(t, r, p, Allow, nil)
	r = request(t, "POST", "/companies/ABC123/employees", nil)
	checkPolicy(t, r, p, Allow, nil)
	r = request(t, "DELETE", "/companies/ABC123/employees", nil)
	checkPolicy(t, r, p, Allow, nil)

	// --

	r = request(t, "GET", "/companies/ABC123/employees", nil)
	p = DenyResource(Read, "/")
	checkPolicy(t, r, p, Pass, nil)

	r = request(t, "GET", "/companies/ABC123/employees", nil)
	p = DenyResource(Read, "/*")
	checkPolicy(t, r, p, Deny, nil)

	r = request(t, "GET", "/companies/ABC123/employees", nil)
	p = DenyResource(Read, "/*/")
	checkPolicy(t, r, p, Deny, nil)

	r = request(t, "GET", "/companies/ABC123/employees", nil)
	p = DenyResource(Read, "/companies/*/employees")
	checkPolicy(t, r, p, Deny, nil)

	r = request(t, "GET", "/companies/ABC123/employees", nil)
	p = DenyResource(Read, "/companies/*/locations")
	checkPolicy(t, r, p, Pass, nil)

	r = request(t, "GET", "/companies/ABC123/employees", nil)
	p = DenyResource(Read, "/companies/*/employees/too/specific")
	checkPolicy(t, r, p, Pass, nil)

	r = request(t, "UNSUPPORTED", "/companies/ABC123/employees", nil)
	p = ResourcePolicy{policy{}, ActionSet{Read, Write}, PathSet{Path("/companies/*/employees")}, Deny}
	checkPolicy(t, r, p, Pass, ErrMethodNotSupported)

	p = ResourcePolicy{policy{}, ActionSet{Every}, PathSet{Path("/companies/*/employees")}, Deny}
	r = request(t, "GET", "/companies/ABC123/employees", nil)
	checkPolicy(t, r, p, Deny, nil)
	r = request(t, "PUT", "/companies/ABC123/employees", nil)
	checkPolicy(t, r, p, Deny, nil)
	r = request(t, "POST", "/companies/ABC123/employees", nil)
	checkPolicy(t, r, p, Deny, nil)
	r = request(t, "DELETE", "/companies/ABC123/employees", nil)
	checkPolicy(t, r, p, Deny, nil)

	// --

	r = request(t, "GET", "/companies/ABC123/employees", nil)
	p = ResourcePolicy{policy{}, ActionSet{Read, Write}, PathSet{"/locations", "/companies/*/employees"}, Allow}
	checkPolicy(t, r, p, Allow, nil)

	r = request(t, "GET", "/locations", nil)
	p = ResourcePolicy{policy{}, ActionSet{Read, Write}, PathSet{"/locations", "/companies/*/employees"}, Allow}
	checkPolicy(t, r, p, Allow, nil)

	r = request(t, "PUT", "/companies/ABC123/employees", nil)
	p = ResourcePolicy{policy{}, ActionSet{Read, Write}, PathSet{"/locations", "/companies/*/employees"}, Allow}
	checkPolicy(t, r, p, Allow, nil)

	r = request(t, "PUT", "/locations", nil)
	p = ResourcePolicy{policy{}, ActionSet{Read, Write}, PathSet{"/locations", "/companies/*/employees"}, Allow}
	checkPolicy(t, r, p, Allow, nil)

	r = request(t, "GET", "/elsewhere", nil)
	p = ResourcePolicy{policy{}, ActionSet{Read, Write}, PathSet{"/locations", "/companies/*/employees"}, Allow}
	checkPolicy(t, r, p, Pass, nil)

}

func request(t *testing.T, m, p string, b []byte) *http.Request {
	req, err := http.NewRequest(m, p, bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	return req
}

func checkPolicy(t *testing.T, r *http.Request, p Policy, e Effect, f error) {
	z, err := p.Eval(r)
	fmt.Printf("--> %v %v (%v) -> (%v, %v)\n", r.Method, r.URL.Path, p, z, err)
	if f != nil {
		assert.Equal(t, f, err)
	} else if assert.Nil(t, err, fmt.Sprintf("%v", err)) {
		assert.Equal(t, e, z, fmt.Sprintf("%v %v (%v)\n", r.Method, r.URL.Path, p))
	}
}
