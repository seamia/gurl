// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"

	"github.com/seamia/libs/resolve"
)

var (
	baseUrl     = "https://gurl.seamia.net/test"
	curlOptions = "-i"

	headers              = map[string]string{}
	printResponseHeaders = printResponseHeadersDefault
	generateCurlCommands = generateCurlCommandsDefault
	collectTimingInfo    = collectTimingInfoDefault
	resolveExternalFiles = resolveExternalFilesDefault

	echoSilent         = false
	echoDebug          = false
	echoVerbose        = false
	echoProgress       = echoDefault
	echoMapCommand     = echoDefault
	echoSetCommand     = echoDefault
	echoGetCommand     = echoDefault
	echoPostCommand    = echoDefault
	echoPatchCommand   = echoDefault
	echoDeleteCommand  = echoDefault
	echoHeaderCommand  = echoDefault
	echoIncludeCommand = echoDefault
	echoEchoCommand    = true
	echoRequireCommand = echoDefault
	echoLoadCommand    = echoDefault
	echoSectionCommand = true

	resolver = resolve.New()

	currentFile       = ""
	currentLineNumber = 0
	currentCommand    = ""
	currentSection    = ""

	responsePrettyPrintBody = responsePrettyPrintBodyDefault

	incrementalCounter int64 // used by ${increment}
)

const (
	echoPrefix = "echo."
)

var (
	dials = map[string]*bool{
		"print.response.headers": &printResponseHeaders,
		//	"generate.curl.commands": &generateCurlCommands,
		"collect.timing.info":    &collectTimingInfo,
		"resolve.external.files": &resolveExternalFiles,
		"pretty.print.body":      &responsePrettyPrintBody,

		echoPrefix + "map":      &echoMapCommand,
		echoPrefix + "set":      &echoSetCommand,
		echoPrefix + "get":      &echoGetCommand,
		echoPrefix + "post":     &echoPostCommand,
		echoPrefix + "patch":    &echoPatchCommand,
		echoPrefix + "delete":   &echoDeleteCommand,
		echoPrefix + "header":   &echoHeaderCommand,
		echoPrefix + "progress": &echoProgress,
		echoPrefix + "echo":     &echoEchoCommand,
		echoPrefix + "require":  &echoRequireCommand,
		echoPrefix + "load":     &echoLoadCommand,
	}
)

var (
	savedRequest           *http.Request
	savedResponse          *http.Response
	savedResponseBody      []byte
	persistRequestResponse = false
)

func offline() bool {
	return generateCurlCommands
}
