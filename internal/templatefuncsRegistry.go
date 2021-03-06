// Copyright 2017-present The Hugo Authors. All rights reserved.
//
// Portions Copyright The Go Authors.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"github.com/flysnow-org/soha/deps"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// TemplateFuncsNamespaceRegistry describes a registry of functions that provide
// namespaces.
var TemplateFuncsNamespaceRegistry []func(d *deps.Deps) *TemplateFuncsNamespace

// AddTemplateFuncsNamespace adds a given function to a registry.
func AddTemplateFuncsNamespace(ns func(d *deps.Deps) *TemplateFuncsNamespace) {
	TemplateFuncsNamespaceRegistry = append(TemplateFuncsNamespaceRegistry, ns)
}

// TemplateFuncsNamespace represents a template function namespace.
type TemplateFuncsNamespace struct {
	// The namespace name, "strings", "lang", etc.
	Name string

	// This is the method receiver.
	Context func(v ...interface{}) interface{}

	// Additional info, aliases and examples, per method name.
	MethodMappings map[string]TemplateFuncMethodMapping
}

// TemplateFuncsNamespaces is a slice of TemplateFuncsNamespace.
type TemplateFuncsNamespaces []*TemplateFuncsNamespace

// AddMethodMapping adds a method to a template function namespace.
func (t *TemplateFuncsNamespace) AddMethodMapping(m interface{}, aliases []string, examples [][2]string) {
	if t.MethodMappings == nil {
		t.MethodMappings = make(map[string]TemplateFuncMethodMapping)
	}

	name := methodToName(m)

	// sanity check
	for _, e := range examples {
		if e[0] == "" {
			panic(t.Name + ": Empty example for " + name)
		}
	}
	for _, a := range aliases {
		if a == "" {
			panic(t.Name + ": Empty alias for " + name)
		}
	}

	t.MethodMappings[name] = TemplateFuncMethodMapping{
		Method:   m,
		Aliases:  aliases,
		Examples: examples,
	}

}

// TemplateFuncMethodMapping represents a mapping of functions to methods for a
// given namespace.
type TemplateFuncMethodMapping struct {
	Method interface{}

	// Any template funcs aliases. This is mainly motivated by keeping
	// backwards compatibility, but some new template funcs may also make
	// sense to give short and snappy aliases.
	// Note that these aliases are global and will be merged, so the last
	// key will win.
	Aliases []string

	// A slice of input/expected examples.
	// We keep it a the namespace level for now, but may find a way to keep track
	// of the single template func, for documentation purposes.
	// Some of these, hopefully just a few, may depend on some test data to run.
	Examples [][2]string
}

func methodToName(m interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(m).Pointer()).Name()
	name = filepath.Ext(name)
	name = strings.TrimPrefix(name, ".")
	name = strings.TrimSuffix(name, "-fm")
	return name
}
