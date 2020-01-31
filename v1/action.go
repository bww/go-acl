package acl

import (
	"fmt"
	"strings"
)

var (
	ErrMethodNotSupported = fmt.Errorf("Method not supported")
)

type Action string

const (
	Read   = Action("read")
	Write  = Action("write")
	Delete = Action("delete")
	Every  = Action("*")
)

var methodToAction = map[string]Action{
	"GET":    Read,
	"POST":   Write,
	"PUT":    Write,
	"PATCH":  Write,
	"DELETE": Delete,
}

type ActionSet []Action

func (s ActionSet) Contains(a Action) bool {
	for _, e := range s {
		if e == Every || e == a {
			return true
		}
	}
	return false
}

func (s ActionSet) String() string {
	var b strings.Builder
	for i, e := range s {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(string(e))
	}
	return b.String()
}
