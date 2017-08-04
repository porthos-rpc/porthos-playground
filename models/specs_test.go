package models

import (
	"testing"
)

func TestBodySpecSimple(t *testing.T) {
	specs, err := UnmarshalSpecs([]byte(`{"service":"UserService","specs":{"doSomethingElse":{"contentType":"application/json","body":{"complex":{"Type":"struct","Description":"Required","body":{"value":{"Type":"bool","Description":""}}},"value":{"Type":"float32","Description":"Required"}}},"doSomethingThatReturnsValue":{"contentType":"application/json","body":{"someComplex":{"Type":"complex","Description":"Optional","body":{"value":{"Type":"int","Description":"Required"}}},"value":{"Type":"int","Description":"Required"}}}}}`))

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
