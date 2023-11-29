package main

import (
	"github.com/casbin/casbin/v2"
	"testing"
)

func TestEnforce(t *testing.T) {
	type Request struct {
		Sub, Obj, Act string
		Want          bool
	}

	listTests := []Request{
		{Sub: "admin", Obj: "/", Act: "GET", Want: true},
		{Sub: "admin", Obj: "/clgt", Act: "GET", Want: true},
		{Sub: "admin", Obj: "/login", Act: "GET", Want: true},
		{Sub: "admin", Obj: "/login", Act: "POST", Want: true},
		{Sub: "admin", Obj: "/register", Act: "GET", Want: true},
		{Sub: "admin", Obj: "/register", Act: "POST", Want: true},

		{Sub: "guest", Obj: "/register", Act: "GET", Want: true},
		{Sub: "guest", Obj: "/register", Act: "POST", Want: true},
		{Sub: "guest", Obj: "/inventory", Act: "GET", Want: false},
		{Sub: "guest", Obj: "/inventory", Act: "POST", Want: false},

		{Sub: "sale", Obj: "/payment", Act: "GET", Want: true},
		{Sub: "sale", Obj: "/payment", Act: "POST", Want: true},
		{Sub: "sale", Obj: "/payment/add", Act: "GET", Want: true},
		{Sub: "sale", Obj: "/payment/add", Act: "POST", Want: true},
		{Sub: "sale", Obj: "/inventory", Act: "GET", Want: true},
		{Sub: "sale", Obj: "/inventory", Act: "POST", Want: true},
		{Sub: "sale", Obj: "/inventory/add", Act: "GET", Want: true},
		{Sub: "sale", Obj: "/inventory/add", Act: "POST", Want: true},

		{Sub: "member", Obj: "/gam?tab=1", Act: "POST", Want: true},
		{Sub: "member", Obj: "/gam?tab=1&id=2", Act: "POST", Want: true},

		{Sub: "accountant", Obj: "/payment", Act: "POST", Want: true},
		{Sub: "accountant", Obj: "/payment/edit", Act: "POST", Want: true},
		{Sub: "accountant", Obj: "/payment/?id=5", Act: "POST", Want: true},
	}

	e, _ := casbin.NewEnforcer("./auth_model.conf", "./policy.csv")
	for _, req := range listTests {
		got, _ := e.Enforce(req.Sub, req.Obj, req.Act)
		if got != req.Want {
			t.Errorf("[`%s`<-->`%s`<-->`%s`] >>> got: `%t`, wanted: `%t`", req.Sub, req.Obj, req.Act, got, req.Want)
		}
	}

}
