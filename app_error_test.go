package main

import (
	"bytes"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	er := errorResponse{
		Error: "This is an error",
		Code:  400,
	}
	if er.String() != `{"error":"This is an error"}` {
		t.Fatal("Wrong JSON output")
	}
}

func testErrorHelper(t *testing.T, callWriteError bool) {
	er := errorResponse{
		Error: "This is an error",
		Code:  400,
	}
	recorder := httptest.NewRecorder()
	if callWriteError {
		writeError(recorder, er)
	} else {
		er.Write(recorder)
	}

	// First check the header.
	header := recorder.Header()
	if len(header) != 1 {
		t.Fatal("There should be only one error!")
	}
	if val, ok := header["Content-Type"]; !ok || len(val) != 1 ||
		val[0] != "application/json" {

		t.Fatal("Wrong Content-Type value")
	}

	// Check the code.
	if recorder.Code != 400 {
		t.Fatalf("Response code should be 400, not %v\n", recorder.Code)
	}

	// The body should be the the JSON error itself.
	if recorder.Body.String() != `{"error":"This is an error"}` {
		t.Fatal("Wrong body value")
	}
}

func TestWrite(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})
	log.SetOutput(buffer)

	testErrorHelper(t, false)

	// Test the logger.
	if !strings.Contains(buffer.String(), "ERROR: This is an error") {
		t.Fatal("Wrong logged message!")
	}
}

func TestWriteError(t *testing.T) {
	testErrorHelper(t, true)
}
