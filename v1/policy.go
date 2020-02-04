package acl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bww/go-util/uuid"
)

var ErrUnsupportedPolicyType = fmt.Errorf("Unsupported policy type")

// Implemented by ACL policies
type Policy interface {
	Id() uuid.UUID
	WithId(uuid.UUID) Policy
	Eval(*http.Request) (Effect, error)
}

// An ACL policy persistent representation
type PersistentPolicy struct {
	Id      uuid.UUID       `json:"id" db:"id,pk"`
	Type    string          `json:"type" db:"type"`
	Data    json.RawMessage `json:"data" db:"data"`
	Created time.Time       `json:"created_at" db:"created_at"`
}

const resourcePolicyType = "resource"

func MarshalPolicy(p Policy) (uuid.UUID, string, []byte, error) {
	switch c := p.(type) {
	case ResourcePolicy:
		return marshalResourcePolicy(c)
	default:
		return uuid.Zero, "", nil, ErrUnsupportedPolicyType
	}
}

func UnmarshalPolicy(id uuid.UUID, t string, d json.RawMessage) (Policy, error) {
	switch t {
	case resourcePolicyType:
		return unmarshalResourcePolicy(id, d)
	default:
		return nil, ErrUnsupportedPolicyType
	}
}

type policy struct {
	id uuid.UUID `json:"-"`
}

type ResourcePolicy struct {
	policy
	Actions ActionSet `json:"actions"`
	Paths   PathSet   `json:"paths"`
	Effect  Effect    `json:"effect"`
}

func marshalResourcePolicy(p ResourcePolicy) (uuid.UUID, string, []byte, error) {
	d, err := json.Marshal(p)
	return p.id, resourcePolicyType, d, err
}

func unmarshalResourcePolicy(id uuid.UUID, d json.RawMessage) (ResourcePolicy, error) {
	var p ResourcePolicy
	err := json.Unmarshal(d, &p)
	if err != nil {
		return ResourcePolicy{}, err
	}
	p.id = id
	return p, nil
}

func AllowResource(a Action, p Path) ResourcePolicy {
	return ResourcePolicy{policy{}, ActionSet{a}, PathSet{p}, Allow}
}

func DenyResource(a Action, p Path) ResourcePolicy {
	return ResourcePolicy{policy{}, ActionSet{a}, PathSet{p}, Deny}
}

func (p ResourcePolicy) Id() uuid.UUID {
	return p.id
}

func (p ResourcePolicy) WithId(id uuid.UUID) Policy {
	return ResourcePolicy{
		policy: policy{
			id: id,
		},
		Actions: p.Actions,
		Paths:   p.Paths,
		Effect:  p.Effect,
	}
}

func (p ResourcePolicy) Eval(req *http.Request) (Effect, error) {
	a, ok := methodToAction[strings.ToUpper(req.Method)]
	if !ok {
		return Deny, ErrMethodNotSupported
	}
	if !p.Actions.Contains(a) {
		return Pass, nil
	}
	if !p.Paths.MatchesString(req.URL.Path) {
		return Pass, nil
	}
	return p.Effect, nil
}

func (p ResourcePolicy) String() string {
	return fmt.Sprintf("%v: [%v] %v", p.Effect, p.Actions, p.Paths)
}
