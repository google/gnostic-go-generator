// +build ignore
//
// This file is omitted when getting with `go get github.com/googleapis/gnostic/...`
//
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
	"context"
	"fmt"
	"sort"

	"github.com/googleapis/gnostic-go-generator/examples/v2.0/apis_guru/apis_guru"
)

func main() {
	c := apis_guru.NewClient("http://api.apis.guru/v2")

	metrics, err := c.GetMetrics()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", metrics)

	ctx := context.Background()
	apis, err := c.ListAPIs(ctx)
	if err != nil {
		panic(err)
	}

	keys := make([]string, 0)
	for key, _ := range *apis.OK {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		api := (*apis.OK)[key]
		versions := make([]string, 0)
		for key, _ := range api.Versions {
			versions = append(versions, key)
		}
		sort.Strings(versions)
		fmt.Printf("[%s]:%+v\n", key, versions)
	}

	api := (*apis.OK)["xkcd.com"].Versions["1.0.0"]
	fmt.Printf("%+v\n", api.SwaggerUrl)
}
