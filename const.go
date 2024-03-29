// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/fatih/color"

const (
	exitCodeOnError   = 7
	exitCodeOnUsage   = 3
	exitCodeOnSuccess = 1

	lineSeparator  = "\n"
	wordSeparator  = " \t"
	shebang        = "#!/"
	itemsSeparator = "/"

	leadingWhiteSpace  = " \t"
	trailingWhiteSpace = " \t\r\n"

	commentPrefix             = "#"
	userAgent                 = "seamia/gurl"
	envDefaultsLocation       = "GURL_DEFAULT_SETTINGS"
	configurationHeaderPrefix = "header:"
	externalFilePrefix        = "@"

	printResponseHeadersDefault = true
	generateCurlCommandsDefault = false
	collectTimingInfoDefault    = false
	resolveExternalFilesDefault = true

	colorComment           = color.FgGreen
	colorError             = color.FgRed
	colorUsage             = color.FgHiCyan
	colorResponse          = color.FgHiMagenta
	colorDebug             = color.FgMagenta
	colorVerbose           = color.FgHiBlack
	colorResponseSuccess   = color.FgHiGreen
	colorResponseFailure   = color.FgHiRed
	colorResponseAttention = color.FgYellow

	headerContentType     = "Content-Type"
	contentTypeJson       = "application/json"
	headerAttentionSuffix = "-error"

	responsePrettyPrintBodyDefault = true
	fallbackForUnknowBinaryState   = false

	mapSessionKeyName    = "session"
	mapVesionKeyName     = "version"
	mapScripFileName     = "script"
	mapScripFullFileName = "script.full"

	includeAllKey = "*"

	marshalPrefix = ""
	marshalIndent = "    "

	mappingResponseValues = "response:"
	mappingResponsePrefix = "${" + mappingResponseValues

	echoDefault  = false
	indexInvalid = -1

	maxIncludesAllowed = 100
)

var (
	colorSection = []color.Attribute{color.FgHiWhite, color.BgBlack}
)
