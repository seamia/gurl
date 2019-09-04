// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"strconv"
	"strings"
)

var (
	quittingCodes = map[int]bool{}
)

func quitOn(value string) {
	// possible inputs:
	// "" --> remove all
	// "201 400"

	value = trim(value)
	if len(value) == 0 {
		debug("emptying the store of status.codes of when to quit")
		quittingCodes = map[int]bool{}
		return
	}

	for len(value) > 0 {
		value = trim(value)
		split := strings.IndexAny(value, " /t,;")
		remains := ""
		if split > 0 {
			remains = value[split+1:]
			value = value[:split]
		}

		value = trim(value)
		status, err := strconv.ParseInt(value, 10, 64)
		quitOnError(err, "cannot convert [%v] to a legit status code", value)

		if status < 0 {
			status = -status
			debug("removing [%v] as quitting code", status)
			quittingCodes[int(status)] = false
		} else {
			debug("adding [%v] as quitting code", status)
			quittingCodes[int(status)] = true
		}

		value = remains
	}

}

func quitIfRequired(resp *http.Response) {
	status := resp.StatusCode

	if exit, present := quittingCodes[status]; present && exit {
		debug("quitting because we hit one of the quitting codes [%v]", status)
		if echoVerbose {
			displayResponse(resp)
		}
		showRequestResponseInfo()

		quit("Quitting [status: %v]", status)
	}
}

func showRequestResponseInfo() {
	debug("Request/Response Info:")
	if savedRequest != nil {
		debug("Verb [%s], URL: [%s]", savedRequest.Method, savedRequest.URL.String())
	}
	if savedResponse != nil {
		debug("Verb [%s], URL: [%s]", savedResponse.Status, savedRequest.URL.String())

		if savedResponse.Header != nil {
			for key, _ := range savedResponse.Header {
				if attentionNeeded(key) {
					debug("Attention [%s]", savedResponse.Header.Get(key))
				}
			}
		}
	}
}
