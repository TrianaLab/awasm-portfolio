package models

import (
	"fmt"
	"strings"
)

type ResourceBase struct {
	Name       string                 // Unique identifier for the resource
	Namespace  string                 // Namespace (mandatory for namespaced resources)
	Data       map[string]interface{} // Arbitrary key-value data
	Namespaced bool                   // Determines if the resource is namespaced
}

// GetFields returns a map of the resource's fields for display or processing.
func (r ResourceBase) GetFields() map[string]string {
	fields := make(map[string]string)
	fields["Name"] = r.Name
	if r.Namespaced {
		fields["Namespace"] = r.Namespace
	}
	for k, v := range r.Data {
		if strings.HasPrefix(k, "\"") { // Skip quoted keys
			continue
		}
		switch val := v.(type) {
		case ResourceBase:
			fields[k] = val.Name
		case []ResourceBase:
			var resourceNames []string
			for _, resource := range val {
				resourceNames = append(resourceNames, resource.Name)
			}
			fields[k] = strings.Join(resourceNames, ", ")
		default:
			fields[k] = fmt.Sprintf("%v", v)
		}
	}
	return fields
}

// GetScope returns the resource's scope as "namespaced" or "cluster-wide".
func (r ResourceBase) GetScope() string {
	if r.Name == "Namespace" {
		return "cluster-wide"
	}
	if r.Namespaced {
		return "namespaced"
	}
	return "cluster-wide"
}

// GetID returns the unique identifier for the resource (its name).
func (r ResourceBase) GetID() string {
	return r.Name
}
