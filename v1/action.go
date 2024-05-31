package acl

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrMethodNotSupported = fmt.Errorf("Method not supported")
)

type Action string

const (
	Read    = Action("read")
	Write   = Action("write")
	Delete  = Action("delete")
	List    = Action("list")
	Approve = Action("approve")
	Notify  = Action("notify")
	Every   = Action("*")
)

var methodToAction = map[string]Action{
	"GET":    Read,
	"POST":   Write,
	"PUT":    Write,
	"PATCH":  Write,
	"DELETE": Delete,
}

func ActionForRequest(req *http.Request) (Action, error) {
	a, ok := methodToAction[strings.ToUpper(req.Method)]
	if ok {
		return a, nil
	} else {
		return "", ErrMethodNotSupported
	}
}

type Actions []Action

func (a Actions) Len() int {
	return len(a)
}

func (a Actions) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Actions) Less(i, j int) bool {
	return string(a[i]) < string(a[j])
}

func (s Actions) Contains(a Action) bool {
	for _, e := range s {
		if e == Every || e == a {
			return true
		}
	}
	return false
}

func (s Actions) String() string {
	var b strings.Builder
	for i, e := range s {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(string(e))
	}
	return b.String()
}
