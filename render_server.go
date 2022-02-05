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
	"fmt"
	"strings"

	surface "github.com/googleapis/gnostic/surface"
)

func (renderer *Renderer) RenderServer() ([]byte, error) {
	f := NewLineWriter()
	f.WriteLine("// GENERATED FILE: DO NOT EDIT!")
	f.WriteLine(``)
	f.WriteLine("package " + renderer.Package)
	f.WriteLine(``)
	imports := []string{
		"github.com/gorilla/mux",
		"net/http",
	}
	f.WriteLine(``)
	f.WriteLine(`import (`)
	for _, imp := range imports {
		f.WriteLine(`"` + imp + `"`)
	}
	f.WriteLine(`)`)

	f.WriteLine(`func int64Value(s string) (int64, error) {`)
	f.WriteLine(`	return strconv.ParseInt(s, 10, 64)`)
	f.WriteLine(`}`)
	f.WriteLine(``)
	f.WriteLine(`func intValue(s string) (int, error) {`)
	f.WriteLine(`	return strconv.Atoi(s)`)
	f.WriteLine(`}`)
	f.WriteLine(``)
	f.WriteLine(`// This package-global variable holds the user-written Provider for API services.`)
	f.WriteLine(`// See the Provider interface for details.`)
	f.WriteLine(`var provider Provider`)
	f.WriteLine(``)
	f.WriteLine(`// These handlers serve API methods.`)
	f.WriteLine(``)

	for _, method := range renderer.Model.Methods {
		parametersType := renderer.Model.TypeWithTypeName(method.ParametersTypeName)
		responsesType := renderer.Model.TypeWithTypeName(method.ResponsesTypeName)

		f.WriteLine(`// Handler`)
		f.WriteLine(commentForText(method.Description))
		f.WriteLine(`func ` + method.HandlerName + `(w http.ResponseWriter, r *http.Request) {`)
		f.WriteLine(`  var err error`)
		if parametersType != nil {
			f.WriteLine(`// instantiate the parameters structure`)
			f.WriteLine(`parameters := &` + parametersType.Name + `{}`)
			if method.Method == "POST" && parametersType.FieldWithPosition(surface.Position_BODY) != nil {
				f.WriteLine(`// deserialize request from post data`)
				f.WriteLine(`decoder := json.NewDecoder(r.Body)`)
				f.WriteLine(`err = decoder.Decode(&parameters.` +
					parametersType.FieldWithPosition(surface.Position_BODY).FieldName + `)`)
				f.WriteLine(`if err != nil {`)
				f.WriteLine(`	w.WriteHeader(http.StatusBadRequest)`)
				f.WriteLine(`	w.Write([]byte(err.Error() + "\n"))`)
				f.WriteLine(`	return`)
				f.WriteLine(`}`)
			}
			f.WriteLine(`// get request fields in path and query parameters`)
			if parametersType.HasFieldWithPosition(surface.Position_PATH) {
				f.WriteLine(`vars := mux.Vars(r)`)
			}
			if parametersType.HasFieldWithPosition(surface.Position_FORMDATA) {
				f.WriteLine(`r.ParseForm()`)
			}
			for _, field := range parametersType.Fields {
				if field.Position == surface.Position_PATH {
					if field.Type == "string" {
						f.WriteLine(fmt.Sprintf("// %+v", field))
						f.WriteLine(`if value, ok := vars["` + field.Name + `"]; ok {`)
						f.WriteLine(`	parameters.` + field.FieldName + ` = value`)
						f.WriteLine(`}`)
					} else if field.Type == "number" {
						f.WriteLine(fmt.Sprintf("// %+v", field))
						f.WriteLine(`if value, ok := vars["` + field.Name + `"]; ok {`)
						f.WriteLine(`	parameters.` + field.FieldName + `, _ = intValue(value)`)
						f.WriteLine(`}`)
					} else { // field.Type == "integer", "boolean", "array"
						f.WriteLine(fmt.Sprintf("// %+v", field))
						f.WriteLine(`if value, ok := vars["` + field.Name + `"]; ok {`)
						f.WriteLine(`	parameters.` + field.FieldName + `, _ = int64Value(value)`)
						f.WriteLine(`}`)
					}
				} else if field.Position == surface.Position_FORMDATA {
					f.WriteLine(`if len(r.Form["` + field.Name + `"]) > 0 {`)
					f.WriteLine(`	parameters.` + field.FieldName + `, _ = int64Value(r.Form["` + field.Name + `"][0])`)
					f.WriteLine(`}`)
				}
			}
		}

		responseName := ""
		if responsesType != nil {
			responseName = strings.ToLower(method.ResponsesTypeName)
			f.WriteLine(`// instantiate the response structure`)
			f.WriteLine(responseName + ` := &` + method.ResponsesTypeName + `{}`)
		}
		f.WriteLine(`// call the service provider`)
		callLine := `err = provider.` + method.ProcessorName
		if parametersType != nil {
			if responsesType != nil {
				callLine += fmt.Sprintf(`(r.Context(), parameters, %s)`, responseName)
			} else {
				callLine += `(r.Context(), parameters)`
			}
		} else {
			if responsesType != nil {
				callLine += fmt.Sprintf(`(r.Context(), %s)`, responseName)
			} else {
				callLine += `(r.Context())`
			}
		}
		f.WriteLine(callLine)
		f.WriteLine(`if err == nil {`)
		if responsesType != nil {
			f.WriteLine(`  // write the normal response`)
			f.WriteLine(`  encoder := json.NewEncoder(w)`)
			f.WriteLine(fmt.Sprintf(`  encoder.Encode(%s)`, responseName))
			f.WriteLine(`  return`)
		}
		f.WriteLine(`} else {`)
		f.WriteLine(`  w.WriteHeader(http.StatusInternalServerError)`)
		f.WriteLine(`  w.Write([]byte(err.Error() + "\n"))`)
		f.WriteLine(`  return`)
		f.WriteLine(`}`)
		f.WriteLine(`}`)
		f.WriteLine(``)
	}
	f.WriteLine(`// Initialize the API service.`)
	f.WriteLine(`func Initialize(p Provider) *mux.Router {`)
	f.WriteLine(`  provider = p`)
	f.WriteLine(`  var router = mux.NewRouter()`)
	for _, method := range renderer.Model.Methods {
		f.WriteLine(`router.HandleFunc("` + method.Path + `", ` + method.HandlerName + `).Methods("` + method.Method + `")`)
	}
	f.WriteLine(`  http.Handle("/", router)`)
	f.WriteLine(`  return router`)
	f.WriteLine(`}`)
	f.WriteLine(``)
	f.WriteLine(`// Provide the API service over HTTP.`)
	f.WriteLine(`func ServeHTTP(address string) error {`)
	f.WriteLine(`  if provider == nil {`)
	f.WriteLine(`    return errors.New("Use ` + renderer.Package + `.Initialize() to set a service provider.")`)
	f.WriteLine(`  }`)
	f.WriteLine(`  return http.ListenAndServe(address, nil)`)
	f.WriteLine(`}`)
	return f.Bytes(), nil
}
