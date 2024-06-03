package acl

import (
	"fmt"
	"net/url"
	"strings"
)

var errInvalidDomain = fmt.Errorf("Invalid domain")

// A Domain describes the context in which access is granted. All scopes are
// considered in the context of a domain. Domains are expressed as a path
// of typed components.
//
//	workspace:1/project:2/resource:3
type Domain []Component

func ParseDomain(s string) (Domain, error) {
	var d Domain
	err := d.UnmarshalText([]byte(s))
	return d, err
}

func (d Domain) MarshalText() ([]byte, error) {
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

func (d *Domain) UnmarshalText(text []byte) error {
	var p []Component
	s := string(text)
	for len(s) > 0 {
		var v string
		if x := strings.Index(s, "/"); x < 0 {
			v, s = s, ""
		} else {
			v, s = s[:x], s[x+1:]
		}
		var c Component
		err := c.UnmarshalText([]byte(v))
		if err != nil {
			return fmt.Errorf("%w: in %s", err, string(text))
		}
		p = append(p, c)
	}
	*d = p
	return nil
}

type Component struct {
	Type string
	Name string
}

func (c Component) MarshalText() ([]byte, error) {
	sb := strings.Builder{}
	sb.WriteString(url.PathEscape(c.Type))
	sb.WriteString(":")
	sb.WriteString(url.PathEscape(c.Name))
	return []byte(sb.String()), nil
}

func (c *Component) UnmarshalText(text []byte) error {
	s := string(text)
	x := strings.Index(s, ":")
	if x < 0 {
		return fmt.Errorf("%w: no component delimiter in: %s", errInvalidDomain, s)
	}
	p := s[:x]
	t, err := url.PathUnescape(p)
	if err != nil {
		return fmt.Errorf("%w: invalid type in: %s", errInvalidDomain, p)
	}
	p = s[x+1:]
	n, err := url.PathUnescape(p)
	if err != nil {
		return fmt.Errorf("%w: invalid type in: %s", errInvalidDomain, p)
	}
	*c = Component{
		Type: t,
		Name: n,
	}
	return nil
}
