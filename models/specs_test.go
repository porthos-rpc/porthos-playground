package models

import (
	"testing"
)

func TestBodySpecSimple(t *testing.T) {
	specs, err := UnmarshalSpecs([]byte(`{ "service": "UserService", "specs": { "doSomething": { "contentType": "application/json", "body":{ "arg1": "float32", "arg2":" string", "arg3":" int" } }, "doSomethingElse": { "contentType": "application/json", "body":{ "arg1": "bool", "arg2":" int" } } } } `))

	if err != nil {
		t.Fatal("UnmarshalSpecs failed.", err)
	}

	if specs.Service != "UserService" {
		t.Errorf("Got an unexpected service name, %s", specs.Service)
	}

	if len(specs.Specs) != 2 {
		t.Errorf("Got an unexpected specs count, %s", len(specs.Specs))
	}
}
