// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/seamia/libs/printer"
)

func main() {
	if len(os.Args) < 2 || help(filename()) {
		usage()
	}
	debug("args: %v, %v", len(os.Args), os.Args)

	loadDefaults()

	printer.Set(debug) // todo: revisit this
	handlers["include"] = processInclude

	processFile(filename())
	/*
		data, err := ioutil.ReadFile(filename())
		quitOnError(err, "Opening file %s", filename())

		comment(echoProgress, "Processing file %s", filename())
		generate("# generating curls commands from %s", filename())
		currentFile = filename()
		processScript(string(data))
	*/
	os.Exit(exitCodeOnSuccess)
}

func processFile(file string) {
	data, err := ioutil.ReadFile(file)
	quitOnError(err, "Opening file %s", file)

	comment(echoProgress, "Processing file %s", file)
	generate("# generating curls commands from %s", file)
	processChildScript(file, string(data))
}

func processScript(script string) {
	lines := strings.Split(script, lineSeparator)

	// skip shebang
	offset := 1
	if strings.HasPrefix(lines[0], shebang) {
		lines = lines[1:]
		offset++
	}

	command := []string{}
	insideCommentBlock := false

	for lineNumber, line := range lines {
		currentLineNumber = lineNumber + offset

		// ignore whitespace
		line = strings.TrimLeft(line, leadingWhiteSpace)
		line = strings.TrimRight(line, trailingWhiteSpace)

		if insideCommentBlock {
			if strings.HasSuffix(line, "*/") {
				insideCommentBlock = false
			}
			continue
		}

		if strings.HasPrefix(line, "/*") {
			if !strings.HasSuffix(line, "*/") {
				insideCommentBlock = true
			}
			continue
		}

		// ignore comments
		if strings.HasPrefix(line, commentPrefix) {
			continue
		}
		if pound := strings.Index(line, commentPrefix); pound > 0 {
			line = trim(line[:pound])
		}

		if len(line) == 0 {
			processCommand(strings.Join(command, " "))
			command = []string{}
		} else {
			if len(command) == 0 {
				cmd, _ := split(line)
				if !multiLineCommand(cmd) {
					processCommand(line)
					continue
				}
			}
			command = append(command, line)
		}
	}

	// deal with the remains ...
	if len(command) > 0 {
		processCommand(strings.Join(command, " "))
		command = []string{}
	}
}

func processChildScript(name, script string) {

	tmpFile, tmpLineNumber, tmpCommand := currentFile, currentLineNumber, currentCommand

	debug("starting processing [%s]", name)
	currentFile = name
	processScript(script)
	debug("done processing [%s]", name)

	currentFile, currentLineNumber, currentCommand = tmpFile, tmpLineNumber, tmpCommand

}

func processCommand(command string) []string {
	command = trim(command)
	if len(command) == 0 {
		return nil
	}

	currentCommand = command
	fullcmd, payload := split(command)
	cmd, options := splitBy(fullcmd, ":")

	if handler, found := handlers[lower(cmd)]; found {
		handler(payload, options)
	} else {
		/*
			if lower(cmd) == "include" {
				return processInclude(payload,options)
			}
		*/
		quit("Unknown command [%s]", fullcmd)
	}

	return nil
}

type cmdHandler func(params, options string)

var handlers = map[string]cmdHandler{
	"set":    processSet,
	"map":    processMap,
	"header": processHeader,
	// 	"include": processInclude,

	"get":    processGet,
	"patch":  processPatch,
	"post":   processPost,
	"delete": processDelete,

	"echo":    processEcho,
	"require": processRequire,
	"load":    processLoad,
	"section": processSection,
}
