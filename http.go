// Copyright 2019 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/seamia/libs/printer"
)

func call(relativeUrl, verb, data string) {

	// wipe the slate
	savedRequest = nil
	savedResponse = nil
	savedResponseBody = nil

	one := expand(baseUrl) + expand(relativeUrl)
	u, err := url.Parse(one)
	quitOnError(err, "Parsing url [%s]", one)

	/*
		u, err := url.Parse(expand(baseUrl))
		quitOnError(err, "Parsing url [%s]", baseUrl)
		u.Path = path.Join(u.Path, expand(relativeUrl))
	*/

	fullUrl := u.String()

	if generateCurlCommands {
		produceCurlCommand(fullUrl, verb, data)
	} else {
		data = loadExternalFile(data)
		var payload io.Reader
		if len(data) > 0 {
			payload = bytes.NewReader([]byte(data))
		}

		client := &http.Client{}
		request, err := http.NewRequest(strings.ToUpper(verb), fullUrl, payload)

		quitOnError(err, "...")

		// add the headers
		for key, value := range headers {
			if len(key) > 0 && len(value) > 0 {
				request.Header.Set(key, expand(value))
			}
		}
		request.Header.Set("User-Agent", userAgent)

		savedRequest = request
		start := time.Now()
		resp, err := client.Do(request)
		if collectTimingInfo {
			duration := time.Now().Sub(start)
			response("the request took %s", duration.String())
		}

		savedResponse = resp
		if resp != nil && resp.Body != nil {
			data, err := ioutil.ReadAll(resp.Body)
			quitOnError(err, "Ingesting response body")
			savedResponseBody = data
		}

		persistIfNecessary()

		quitOnError(err, "......")
		quitIfRequired(resp)
		displayResponse(resp)
	}
}

func displayResponse(resp *http.Response) {
	if resp == nil {
		response("got an empty response")
	}

	print := responseFailure
	if resp.StatusCode < http.StatusBadRequest {
		print = responseSuccess
	}

	// colorPrint(colorResponse, format, a...)

	print("Status: %s", resp.Status)
	displayHeaders(resp, print)

	if savedResponseBody != nil {

		switch getContentType(resp) {
		case contentTypeJson:
			displayJsonBody(savedResponseBody, print)
		default:
			displayPlainBody(savedResponseBody, print)
		}
	}
	// saveResponse(resp)
}

func displayHeaders(resp *http.Response, print printer.Printer) {
	const format = "\tHeader: [%s] = [%s]"
	if printResponseHeaders && len(resp.Header) > 0 {
		for key := range resp.Header {
			value := resp.Header.Get(key)
			if attentionNeeded(key) {
				responseAttention(format, key, value)
			} else {
				print(format, key, value)
			}
		}
	}
}

func getContentType(resp *http.Response) string {
	return lower(resp.Header.Get(headerContentType))
}

func displayPlainBody(data []byte, print printer.Printer) {
	if len(data) == 0 {
		print("Body is empty.")
	} else {
		print("Body: %s", string(data))
	}
}

func displayJsonBody(data []byte, print printer.Printer) {
	if !responsePrettyPrintBody || len(data) == 0 {
		displayPlainBody(data, print)
		return
	}

	var blank interface{}
	blank = slice{}
	if err := json.Unmarshal(data, &blank); err == nil {
		if pretty, err := json.MarshalIndent(blank, marshalPrefix, marshalIndent); err == nil {
			displayPlainBody(pretty, print)
			return
		} else {
			reportError(err, "marshalling")
		}
	} else {
		reportError(err, "unmarshalling")
	}

	blank = msi{}
	if err := json.Unmarshal(data, &blank); err == nil {
		if pretty, err := json.MarshalIndent(blank, marshalPrefix, marshalIndent); err == nil {
			displayPlainBody(pretty, print)
			return
		} else {
			reportError(err, "marshalling")
		}
	} else {
		reportError(err, "unmarshalling")
	}
	displayPlainBody(data, print)
}

func attentionNeeded(key string) bool {
	if strings.HasSuffix(lower(key), headerAttentionSuffix) {
		return true
	}
	return false
}

type SavedResponse struct {
	Response struct {
		Status     string      `json:"status"`
		StatusCode int         `json:"status-code"`
		Header     http.Header `json:"headers"`
	} `json:""`
	Request struct {
		Url    string `json:"url"`
		Method string `json:"method"`
	} `json:""`
}

func saveResponse(resp *http.Response) {
	if resp == nil {
		return
	}

	saved := SavedResponse{}
	saved.Response.Status = resp.Status
	saved.Response.StatusCode = resp.StatusCode
	saved.Response.Header = resp.Header

	saved.Request.Method = resp.Request.Method
	saved.Request.Url = resp.Request.URL.String()

	if data, err := json.MarshalIndent(&saved, marshalPrefix, marshalIndent); err == nil {
		fmt.Println(string(data))
	} else {
		fmt.Println(err)
	}

}

var persistenceCounter int64

func persistIfNecessary() {
	if !persistRequestResponse {
		return
	}

	current := atomic.AddInt64(&persistenceCounter, 1)
	filename := fmt.Sprintf("state_%v.txt", current)
	file, err := os.Create(filename)
	if err != nil {
		reportError(err, "failed to open file [%s]", filename)
		return
	}
	defer file.Close()
	out := func(format string, args ...interface{}) {
		_, _ = fmt.Fprintf(file, format+"\n", args...)
	}
	saveRequestResponseInfo(out)
}

func saveRequestResponseInfo(out printer.Printer) {
	out("Section: [%s]", currentSection)
	out("Script:  [%s]", currentFile)
	out("Line:    [%v]", currentLineNumber)
	out("Command: [%s]", currentCommand)
	out("Time:    [%v]", time.Now())
	out("")

	out("Request Info:")
	if savedRequest != nil {
		out("Verb [%s]", savedRequest.Method)
		out("URL: [%s]", savedRequest.URL.String())
		printHeaders(savedRequest.Header, out)
	}

	out("")
	out("Response Info:")
	if savedResponse != nil {
		out("Status [%s]", savedResponse.Status)
		printHeaders(savedResponse.Header, out)

		switch getContentType(savedResponse) {
		case contentTypeJson:
			displayJsonBody(savedResponseBody, out)
		default:
			displayPlainBody(savedResponseBody, out)
		}
	}
}

func printHeaders(head http.Header, out printer.Printer) {
	all := make([]string, 0, len(head))
	for key, _ := range head {
		all = append(all, key)
	}

	sort.Strings(all)

	for _, key := range all {
		out("\t[%s]\t [%s]", key, head.Get(key))
	}
}
