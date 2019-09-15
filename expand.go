// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"
)

func expand(from string) string {

	/*
		if len(from) == 0 || strings.Index(from, "%{") < 0 {
			return from
		}

		if from == "admin_role" {
			print()
		}

	*/

	if i := strings.Index(from, mappingResponsePrefix); i >= 0 {
		prefix := from[:i]
		suffix := ""
		body := ""

		counter := 1
		for j := i + len(mappingResponsePrefix); j < len(from) && counter > 0; j++ {
			switch from[j] {
			case '{':
				counter++
			case '}':
				counter--
				if counter == 0 {
					suffix = from[j+1:]
					body = from[i+len(mappingResponsePrefix) : j]
				}
			}
		}
		if counter != 0 {
			quit("missing closing part of [%s]", from)
		}

		if okay, value := responseValue(expand(body)); okay {
			body = value
		} else {
			quit("failed to resolve [%s]", body)
		}

		result := expand(prefix) + body + expand(suffix)
		return result
	}

	return resolver.Text(from)
}
