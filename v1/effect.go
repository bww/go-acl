package acl

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Effect int

const (
	InvalidEffect Effect = iota
	Allow
	Deny
	Pass
)

var effectToName = map[Effect]string{
	Allow: "allow",
	Deny:  "deny",
	Pass:  "pass",
}

var nameToEffect = map[string]Effect{
	"allow": Allow,
	"deny":  Deny,
	"pass":  Pass,
}

func ParseEffect(s string) (Effect, error) {
	e, ok := nameToEffect[s]
	if !ok {
		return InvalidEffect, fmt.Errorf("Invalid effect: %q", s)
	}
	return e, nil
}

func (e Effect) Inverse() Effect {
	if e == Deny {
		return Allow
	} else {
		return Deny
	}
}

func (e Effect) String() string {
	return effectToName[e]
}

func (e Effect) Value() (driver.Value, error) {
	return e.String(), nil
}

func (e *Effect) Scan(src interface{}) error {
	var err error
	switch c := src.(type) {
	case []byte:
		*e, err = ParseEffect(string(c))
	case string:
		*e, err = ParseEffect(c)
	default:
		err = fmt.Errorf("Unsupported type: %T", src)
	}
	return err
}

func (e Effect) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *Effect) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*e, err = ParseEffect(s)
	if err != nil {
		return err
	}
	return nil
}
