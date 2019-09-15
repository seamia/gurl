// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"strings"
)

const (
	blank = "    "
)

func produceCurlCommand(fullUrl, verb, data string) {
	printer := generate

	printer("# %s %s", verb, fullUrl)
	printer("curl \\")
	if len(curlOptions) > 0 {
		printer(blank+"%s \\", curlOptions)
	}
	printer(blank+"--request %s \\", strings.ToUpper(verb))
	printer(blank+"--url %s \\", fullUrl)

	for key, value := range headers {
		if len(key) > 0 && len(value) > 0 {
			//   --header 'origin: ${value}'   \
			printer(blank+"--header '%s: %s'   \\", key, value)
		}
	}

	if len(data) > 0 {
		external, filename := dataPointsToExternalFile(data)
		if external {
			printer(blank+"--data-binary \"@%s\"", filename)
		} else {
			printer(blank+"--data '%s'", data)
		}
	}
	printer("")
}
