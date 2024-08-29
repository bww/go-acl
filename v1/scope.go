package acl

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

var (
	errInvalidScope  = errors.New("Invalid scope")
	errInvalidAction = errors.New("Invalid action")
	errEmptyResource = errors.New("Empty resource")
)

type Scope struct {
	Actions  Actions `json:"actions"`
	Resource string  `json:"resource"`
}

func NewScope(r string, a ...Action) Scope {
	return Scope{a, r}
}

func ParseScope(s string) (Scope, error) {
	var err error

	var a Actions
	a, s, err = parseActions(s)
	if err != nil {
		return Scope{}, err
	}
	if s == "" {
		return Scope{}, errEmptyResource
	}

	return Scope{a, s}, nil
}

func (s Scope) Satisfies(r Scope) bool {
	if len(r.Actions) < 1 || r.Resource == "" {
		return false
	}
	if len(s.Actions) < 1 || s.Resource == "" {
		return false
	}
	if s.Resource != r.Resource {
		return false
	}
	for _, e := range r.Actions {
		if !s.Actions.Contains(e) {
			return false
		}
	}
	return true
}

func (s Scope) String() string {
	var b strings.Builder
	for i, a := range s.Actions {
		if a == Every {
			b.Reset()
			b.WriteString(string(Every))
			break
		}
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(string(a))
	}
	if len(s.Actions) > 0 {
		b.WriteString(":")
	}
	b.WriteString(s.Resource)
	return b.String()
}

func (s Scope) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Scope) UnmarshalJSON(data []byte) error {
	var v string
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	x, err := ParseScope(v)
	if err != nil {
		return err
	}
	*s = x
	return nil
}

func (s Scope) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *Scope) UnmarshalText(data []byte) error {
	v, err := ParseScope(string(data))
	if err != nil {
		return err
	}
	*s = v
	return nil
}

func (s Scope) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s *Scope) Scan(src interface{}) error {
	var err error
	var v Scope
	switch c := src.(type) {
	case []byte:
		v, err = ParseScope(string(c))
	case string:
		v, err = ParseScope(c)
	default:
		err = fmt.Errorf("Unsupported type: %T", src)
	}
	if err != nil {
		return err
	}
	*s = v
	return nil
}

type Scopes []Scope

func Union(s ...Scopes) Scopes {
	var m Scopes
	for _, e := range s {
		m = append(m, e...)
	}
	if m == nil {
		return nil
	} else {
		return m.Merged()
	}
}

func (s Scopes) Len() int {
	return len(s)
}

func (s Scopes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Scopes) Less(i, j int) bool {
	return string(s[i].Resource) < string(s[j].Resource)
}

func (s Scopes) String() string {
	b := &strings.Builder{}
	for i, e := range s {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(e.String())
	}
	return b.String()
}

func (s Scopes) Add(a Scopes) Scopes {
	return append(s, a...)
}

func (s Scopes) Merge(a Scopes) Scopes {
	var m Scopes
	m = append(m, s...)
	m = append(m, a...)
	return m.Merged()
}

func (s Scopes) Merged() Scopes {
	m := make(map[string]Actions)
	for _, e := range s {
		r, ok := m[e.Resource]
		if !ok {
			r = make(Actions, 0)
		}
		for _, x := range e.Actions {
			if !r.Contains(x) {
				r = append(r, x)
			}
		}
		m[e.Resource] = r
	}
	r := make(Scopes, len(m))
	i := 0
	for k, v := range m {
		if v.Contains(Every) {
			r[i] = NewScope(k, Every)
		} else {
			r[i] = NewScope(k, v...)
		}
		i++
	}
	return r
}

func (s Scopes) Satisfies(r ...Scope) bool {
outer:
	for _, x := range r {
		for _, e := range s {
			if e.Satisfies(x) {
				continue outer
			}
		}
		return false
	}
	return true
}

func (s Scopes) Value() (driver.Value, error) {
	c := make([]string, len(s))
	for i, e := range s {
		c[i] = e.String()
	}
	return pq.Array(c).Value()
}

func (s *Scopes) Scan(src interface{}) error {
	a := pq.StringArray{}
	err := a.Scan(src)
	if err != nil {
		return err
	}
	x := make(Scopes, len(a))
	for i, e := range a {
		v, err := ParseScope(e)
		if err != nil {
			return err
		}
		x[i] = v
	}
	*s = x
	return nil
}

func parseActions(s string) (Actions, string, error) {
	var a Actions
	var every bool

	for {
		x := strings.IndexAny(s, ",:")
		if x < 0 {
			break
		}
		switch s[:x] {
		case "":
			// empty action, ignore this
		case string(Read):
			a = append(a, Read)
		case string(Write):
			a = append(a, Write)
		case string(Delete):
			a = append(a, Delete)
		case string(List):
			a = append(a, List)
		case string(Every):
			every = true
		default:
			return nil, s, errInvalidAction
		}
		d := s[x]
		s = s[x+1:]
		if d == ':' {
			break
		}
	}

	if every {
		return Actions{Every}, s, nil
	} else {
		return a, s, nil
	}
}
