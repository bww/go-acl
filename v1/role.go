package acl

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

var errInvalidRole = fmt.Errorf("Invalid role")

type Role string

const (
	None   = Role("none")
	Member = Role("member")
	Admin  = Role("admin")
	Owner  = Role("owner")
	Self   = Role("self")
)

var roleNames = map[Role]string{
	None:   "None",
	Member: "Member",
	Admin:  "Admin",
	Owner:  "Owner",
	Self:   "Self",
}

var roleGrant = map[Role]Roles{ // roles that a particular role can grant
	None:   nil,
	Member: nil,
	Admin:  Roles{Member, Admin},
	Owner:  Roles{Member, Admin, Owner},
	Self:   nil,
}

func ParseRole(s string) (Role, error) {
	c := Role(s)
	_, ok := roleNames[c]
	if ok {
		return c, nil
	} else {
		return "", errInvalidRole
	}
}

func (c Role) String() string {
	return string(c)
}

func (c Role) Name() string {
	n, ok := roleNames[c]
	if ok {
		return n
	} else {
		return "Invalid"
	}
}

func (c Role) CanGrant(r Role) bool {
	for _, e := range roleGrant[c] {
		if e == r {
			return true
		}
	}
	return false
}

func (c Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Role) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	v, err := ParseRole(s)
	if err != nil {
		return err
	}
	*c = v
	return nil
}

func (c Role) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

func (c *Role) UnmarshalText(data []byte) error {
	v, err := ParseRole(string(data))
	if err != nil {
		return err
	}
	*c = v
	return nil
}

func (c Role) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *Role) Scan(src interface{}) error {
	var err error
	var v Role
	switch c := src.(type) {
	case []byte:
		v, err = ParseRole(string(c))
	case string:
		v, err = ParseRole(c)
	default:
		err = fmt.Errorf("Unsupported type: %T", src)
	}
	if err != nil {
		return err
	}
	*c = v
	return nil
}

type Roles []Role

func (s Roles) Len() int {
	return len(s)
}

func (s Roles) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Roles) Less(i, j int) bool {
	return string(s[i]) < string(s[j])
}

func (s Roles) String() string {
	b := &strings.Builder{}
	for i, e := range s {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(e.String())
	}
	return b.String()
}

func (s Roles) CanGrant(r Role) bool {
	for _, e := range s {
		if e.CanGrant(r) {
			return true
		}
	}
	return false
}

func (s Roles) Merge(r Roles) Roles {
	var m Roles
	t := make(map[Role]struct{})
	for _, e := range s {
		t[e] = struct{}{}
	}
	for _, e := range r {
		t[e] = struct{}{}
	}
	for k, _ := range t {
		m = append(m, k)
	}
	return m
}

func (s Roles) Contains(a Role) bool {
	for _, e := range s {
		if e == a {
			return true
		}
	}
	return false
}

func (s Roles) Value() (driver.Value, error) {
	var c = make([]string, len(s))
	for i, e := range s {
		c[i] = e.String()
	}
	return pq.Array(c).Value()
}

func (s *Roles) Scan(src interface{}) error {
	a := pq.StringArray{}
	err := a.Scan(src)
	if err != nil {
		return err
	}
	r := make(Roles, len(a))
	for i, e := range a {
		v, err := ParseRole(e)
		if err != nil {
			return err
		}
		r[i] = v
	}
	*s = r
	return nil
}
