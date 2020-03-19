// Copyright 2019 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"strings"
	"unicode"

	gnostic_grpc "github.com/googleapis/gnostic-grpc/generator"
	surface_v1 "github.com/googleapis/gnostic/surface"
)

type GoLanguageModel struct{}

func NewGoLanguageModel() *GoLanguageModel {
	return &GoLanguageModel{}
}

// Prepare sets language-specific properties for all types and methods.
func (language *GoLanguageModel) Prepare(model *surface_v1.Model, inputDocumentType string) {
	for _, t := range model.Types {
		// determine the type used for Go language implementation of the type
		t.TypeName = strings.Title(filteredTypeName(t.Name))

		for _, f := range t.Fields {
			f.FieldName = goFieldName(f.Name, f.Type)
			f.ParameterName = goParameterName(f.Name, f.Type)
			switch f.Type {
			case "boolean":
				f.NativeType = "bool"
			case "number":
				f.NativeType = "int"
			case "integer":
				switch f.Format {
				case "int32":
					f.NativeType = "int32"
				case "int64":
					f.NativeType = "int64"
				default:
					f.NativeType = "int64"
				}
			case "object":
				f.NativeType = "interface{}"
			case "string":
				f.NativeType = "string"
			default:
				f.NativeType = strings.Title(filteredTypeName(f.Type))
			}
		}
	}

	for _, m := range model.Methods {
		m.HandlerName = "Handle" + m.Name
		m.ProcessorName = m.Name
		m.ClientName = m.Name
		m.ResponsesTypeName = strings.Title(filteredTypeName(m.ResponsesTypeName))
	}
	gnostic_grpc.AdjustSurfaceModel(model, inputDocumentType)
}

func goParameterName(originalName string, t string) string {
	name := gnostic_grpc.CleanName(originalName)
	if len(name) == 0 {
		name = gnostic_grpc.CleanName(t)
	}
	// lowercase first letter
	a := []rune(name)
	a[0] = unicode.ToLower(a[0])
	name = string(a)
	// avoid reserved words
	if name == "type" {
		return "myType"
	}
	return name
}

func goFieldName(name string, t string) string {
	name = gnostic_grpc.CleanName(name)
	if len(name) == 0 {
		name = gnostic_grpc.CleanName(t)
	}
	name = snakeCaseToCamelCaseWithCapitalizedFirstLetter(name)
	return name
}

func snakeCaseToCamelCaseWithCapitalizedFirstLetter(snakeCase string) (camelCase string) {
	isToUpper := false
	for _, runeValue := range snakeCase {
		if isToUpper {
			camelCase += strings.ToUpper(string(runeValue))
			isToUpper = false
		} else {
			if runeValue == '_' {
				isToUpper = true
			} else {
				camelCase += string(runeValue)
			}
		}
	}
	camelCase = strings.Title(camelCase)
	return
}

func filteredTypeName(typeName string) (name string) {
	// first take the last path segment
	parts := strings.Split(typeName, "/")
	name = parts[len(parts)-1]
	// then take the last part of a dotted name
	parts = strings.Split(name, ".")
	name = parts[len(parts)-1]
	return name
}
