package acl

import (
	"fmt"
	"net/url"
	"strings"
)

var errInvalidRealm = fmt.Errorf("Invalid realm")

// A Realm describes the context in which access is granted. All scopes are
// considered in the context of a relam. Realms are expressed as a path of
// typed components.
//
//	workspace:1/project:2/resource:3
type Realm []Element

func ParseRealm(s string) (Realm, error) {
	var d Realm
	err := d.UnmarshalText([]byte(s))
	return d, err
}

func (r Realm) Len() int {
	return len(r)
}

func (r Realm) Shift() (Element, Realm) {
	if len(r) > 0 {
		return r[0], Realm(r[1:])
	} else {
		return Element{}, Realm{}
	}
}

func (r Realm) Contains(v Realm) bool {
	if len(v) < len(r) {
		return false // param realm has fewer components; receiver cannot contain it
	}
	for i, e := range r {
		if !e.Equals(v[i]) {
			return false
		}
	}
	return true
}

func (r Realm) String() string {
	t, err := r.MarshalText()
	if err != nil {
		panic(err) // this should never happen
	}
	return string(t)
}

func (d Realm) MarshalText() ([]byte, error) {
	sb := strings.Builder{}
	for i, e := range d {
		if i > 0 {
			sb.WriteString("/")
		}
		s, err := e.MarshalText()
		if err != nil {
			return nil, err
		}
		sb.WriteString(string(s))
	}
	return []byte(sb.String()), nil
}

func (d *Realm) UnmarshalText(text []byte) error {
	var p []Element
	s := string(text)
	for len(s) > 0 {
		var v string
		if x := strings.Index(s, "/"); x < 0 {
			v, s = s, ""
		} else {
			v, s = s[:x], s[x+1:]
		}
		var c Element
		err := c.UnmarshalText([]byte(v))
		if err != nil {
			return fmt.Errorf("%w: in %s", err, string(text))
		}
		p = append(p, c)
	}
	*d = p
	return nil
}

type Element struct {
	Type string
	Name string
}

func (c Element) Equals(v Element) bool {
	return c.Type == v.Type && c.Name == v.Name
}

func (c Element) MarshalText() ([]byte, error) {
	sb := strings.Builder{}
	sb.WriteString(url.PathEscape(c.Type))
	if c.Name != "" {
		sb.WriteString(":")
		sb.WriteString(url.PathEscape(c.Name))
	}
	return []byte(sb.String()), nil
}

func (c *Element) UnmarshalText(text []byte) error {
	s := string(text)
	var err error

	var l, r string
	if x := strings.Index(s, ":"); x >= 0 {
		l, r = s[:x], s[x+1:]
	} else {
		l = s
	}

	var t string
	t, err = url.PathUnescape(l)
	if err != nil {
		return fmt.Errorf("%w: invalid type in: %s", errInvalidRealm, l)
	}

	var n string
	if r != "" {
		n, err = url.PathUnescape(r)
		if err != nil {
			return fmt.Errorf("%w: invalid type in: %s", errInvalidRealm, r)
		}
	}

	*c = Element{
		Type: t,
		Name: n,
	}
	return nil
}
