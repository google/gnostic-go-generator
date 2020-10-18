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
	"log"

	docopt "github.com/docopt/docopt-go"
	"github.com/googleapis/gnostic-go-generator/examples/googleauth"
	"github.com/googleapis/gnostic-go-generator/examples/v3.0/urlshortener/urlshortener"
)

func main() {
	usage := `
Usage:
	urlshortener get <url>
	urlshortener list
	urlshortener insert <url>
	`
	arguments, err := docopt.Parse(usage, nil, false, "URL Shortener 1.0", false)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	path := "https://www.googleapis.com/urlshortener/v1" // this should be generated

	client, err := googleauth.NewOAuth2Client("https://www.googleapis.com/auth/urlshortener")
	if err != nil {
		log.Fatalf("Error building OAuth client: %v", err)
	}
	c := urlshortener.NewClient(path, client)

	// get
	if arguments["get"].(bool) {
		ctx := context.Background()
		response, err := c.Urlshortener_Url_Get(ctx, "FULL", arguments["<url>"].(string))
		if err != nil {
			log.Fatalf("%+v", err)
		}
		fmt.Println(response.Default.LongUrl)
	}

	// list
	if arguments["list"].(bool) {
		ctx := context.Background()
		response, err := c.Urlshortener_Url_List(ctx, "", "")
		if err != nil {
			log.Fatalf("%+v", err)
		}
		for _, item := range response.Default.Items {
			fmt.Printf("%-40s %s\n", item.Id, item.LongUrl)
		}
	}

	// insert
	if arguments["insert"].(bool) {
		ctx := context.Background()
		var url urlshortener.Url
		url.LongUrl = arguments["<url>"].(string)
		response, err := c.Urlshortener_Url_Insert(ctx, url)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		fmt.Printf("%+v\n", response.Default.Id)
	}
}
