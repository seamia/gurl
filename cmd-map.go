// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/url"
	"strings"
)

func processMap(params, options string) {
	comment(echoMapCommand, "MAP command: %s", params)
	key, value := split(expand(params))

	if offline() {
		generate("%s=%s", key, value)
	}

	if len(options) > 0 {
		if lower(options) == "encode" {
			value = url.QueryEscape(value)
		} else if lower(options) == "coalesce" {
			key, choices := split(params)
			key = expand(key)
			for _, candidate := range strings.Split(choices, " ") {
				value := expand(candidate)
				if len(value) > 0 {
					resolver.Add(key, value)
					debug("resolved [%s] to be [%s]", key, value)
					return
				}
			}
			resolver.Add(key, "")
			debug("failed to find non-empty value for key [%s] in [%s]", key, params)
			return
		} else {
			quit("unknown options: %s", options)
		}
	}

	resolver.Add(key, value)
}
