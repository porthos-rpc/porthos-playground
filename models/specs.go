package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
)

// FieldSpec represents a spec of a body field.
type FieldSpec struct {
	Type        string      `json:type`
	Description string      `json:description`
	Body        BodySpecMap `json:"body,omitempty"`
}

// BodySpecMap represents a body spec.
type BodySpecMap map[string]FieldSpec

// Spec to a remote procedure.
type Spec struct {
	ContentType string      `json:"contentType"`
	Body        BodySpecMap `json:"body"`
}

// SpecsMap map of string and Spec
type SpecsMap map[string]Spec

// ServiceSpecs is the model sent by a porthos service with all the registered specs.
type ServiceSpecs struct {
	Service string   `db:"service" json:"service"`
	Specs   SpecsMap `db:"specs" json:"specs"`
}

// Value prepares the value for the database.
func (p SpecsMap) Value() (driver.Value, error) {
	if len(p) == 0 {
		return nil, nil
	}

	return json.Marshal(p)
}

// Scan prepares the value from the db to the go code.
func (p SpecsMap) Scan(src interface{}) error {
	v := reflect.ValueOf(src)
	if !v.IsValid() || v.IsNil() {
		return nil
	}
	if data, ok := src.([]byte); ok {
		return json.Unmarshal(data, &p)
	}
	return fmt.Errorf("Could not not decode type %T -> %T", src, p)
}

// Value prepares the value for the database.
func (p BodySpecMap) Value() (driver.Value, error) {
	if len(p) == 0 {
		return nil, nil
	}

	return json.Marshal(p)
}

// Scan prepares the value from the db to the go code.
func (p BodySpecMap) Scan(src interface{}) error {
	v := reflect.ValueOf(src)
	if !v.IsValid() || v.IsNil() {
		return nil
	}
	if data, ok := src.([]byte); ok {
		return json.Unmarshal(data, &p)
	}
	return fmt.Errorf("Could not not decode type %T -> %T", src, p)
}

// UnmarshalSpecs takes an array of bytes (json encoded), allocates and returns service specs.
func UnmarshalSpecs(bytes []byte) (*ServiceSpecs, error) {
	var specs *ServiceSpecs
	err := json.Unmarshal(bytes, &specs)

	return specs, err
}
