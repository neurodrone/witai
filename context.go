package witai

// #cgo LDFLAGS: -L. -lwit -lssl -lcrypto -lsox
// #include <stdio.h>
// #include <stdlib.h>
// #include "wit.h"
//
// extern void asyncCallback(char *);
import "C"

import (
	"encoding/json"
	"syscall"
	"unsafe"
)

var (
	callback func(string)
)

// Context holds the wit context and the access token and the logic
// to make text and voice queries to wit.ai.
type Context struct {
	context      *C.struct_wit_context
	access_token *C.char
}

// NewContext returns a *Context given a device, access token and the level
// of logging verbosity.
func NewContext(dev string, access_token string, v Verbosity) (*Context, error) {
	// The access token C string is freed once the context is closed.
	c_access_token := C.CString(access_token)

	c_device := C.CString(dev)
	defer C.free(unsafe.Pointer(c_device))

	context, err := C.wit_init(c_device, C.uint(v))
	if err != nil {
		switch e, ok := err.(syscall.Errno); {
		case !ok:
			return nil, err

		// There are cases where init returns a non-zero errno when it
		// doesn't find a handler for `alsa`, which is a commonly seen
		// error on OSX. It doesn't stop the implementation from working
		// correctly though, so we specifically ignore the ENOENT error.
		case e == syscall.ENOENT:
			break
		default:
			return nil, err
		}
	}

	return &Context{
		context:      context,
		access_token: c_access_token,
	}, nil
}

// TextQuery allows querying wit.ai service with a text query string.
// It returns the complete result, the first outcome or a non-nil
// error if one occurs.
//
// The result is returned in a blocking manner.
func (c *Context) TextQuery(query string) (*Result, *Outcome, error) {
	c_query := C.CString(query)
	defer C.free(unsafe.Pointer(c_query))

	c_result, err := C.wit_text_query(c.context, c_query, c.access_token)
	if err != nil {
		return nil, nil, err
	}
	defer C.free(unsafe.Pointer(c_result))

	return c.parseResult(C.GoString(c_result))
}

// TextQueryAsync works similar to TextQuery but is not blocking in nature.
// It needs callback to process the response once it arrives. If there is
// an error registering the callback it is returned immediately.
func (c *Context) TextQueryAsync(query string, cb func(string)) error {
	c_query := C.CString(query)
	defer C.free(unsafe.Pointer(c_query))

	callback = cb
	_, err := C.wit_text_query_async(c.context, c_query, c.access_token, (C.wit_resp_callback)(unsafe.Pointer(C.asyncCallback)))
	return err
}

// VoiceQueryAuto is used to query wit.ai service using voice commands
// via the given input device. It returns the complete result, the first
// outcome or a non-nil error if one occurs.
//
// This query is run in a blocking manner.
func (c *Context) VoiceQueryAuto() (*Result, *Outcome, error) {
	c_result, err := C.wit_voice_query_auto(c.context, c.access_token)
	if err != nil {
		return nil, nil, err
	}
	defer C.free(unsafe.Pointer(c_result))

	return c.parseResult(C.GoString(c_result))
}

// VoiceQueryAutoAsync works similar to VoiceQueryAuto but isn't blocking
// in nature. It registers a callback, which is called once the voice command
// query returns.
func (c *Context) VoiceQueryAutoAsync(cb func(string)) error {
	callback = cb
	_, err := C.wit_voice_query_auto_async(c.context, c.access_token, (C.wit_resp_callback)(unsafe.Pointer(C.asyncCallback)))
	return err
}

// VoiceQueryStart starts voice recording and returns immediately. The voice
// recording is stopped using VoiceQueryStop or VoiceQueryStopAsync command.
func (c *Context) VoiceQueryStart() error {
	_, err := C.wit_voice_query_start(c.context, c.access_token)
	return err
}

// VoiceQueryStop stops the voice recording and issues the voice query to
// wit.ai and blockingly waits for the response.
func (c *Context) VoiceQueryStop() (*Result, *Outcome, error) {
	c_result, err := C.wit_voice_query_stop(c.context)
	if err != nil {
		return nil, nil, err
	}
	defer C.free(unsafe.Pointer(c_result))

	return c.parseResult(C.GoString(c_result))
}

// VoiceQueryStopAsync stops the voice recording and issues the voice query to
// wit.ai in a non-blocking manner.
func (c *Context) VoiceQueryStopAsync(cb func(string)) error {
	callback = cb
	_, err := C.wit_voice_query_stop_async(c.context, (C.wit_resp_callback)(unsafe.Pointer(C.asyncCallback)))
	return err
}

// Close closes the context that is used for relaying queries and their
// results back.
func (c *Context) Close() error {
	C.free(unsafe.Pointer(c.access_token))
	_, err := C.wit_close(c.context)
	return err
}

func (c *Context) parseResult(result string) (*Result, *Outcome, error) {
	var r Result
	if err := json.Unmarshal([]byte(result), &r); err != nil {
		return nil, nil, err
	}

	if !r.IsValid() {
		return nil, nil, ErrInvalidResult
	}

	return &r, r.Outcomes[0], nil
}

//export asyncCallback
func asyncCallback(result *C.char) {
	defer C.free(unsafe.Pointer(result))
	if callback != nil {
		callback(C.GoString(result))
		callback = nil
	}
}
