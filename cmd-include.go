// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

var (
	includeCounter int // this is a simple quard against endless recursion
)

func processInclude(params, options string) {
	comment(echoIncludeCommand, "INCLUDE command: %s", params)

	includeCounter++
	if includeCounter > maxIncludesAllowed {
		quit("there were too many INCLUDE calls... there is a good chance you you have an include 'loop' ...")
	}

	processFile(expand(params))
}
