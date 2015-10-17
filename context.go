package witai

// #cgo LDFLAGS: -L. -lwit -lssl -lcrypto -lsox
// #include <stdio.h>
// #include <stdlib.h>
// #include "wit.h"
import "C"

import (
	"encoding/json"
	"syscall"
	"unsafe"
)

type Context struct {
	context      *C.struct_wit_context
	access_token *C.char
}

func NewContext(dev string, access_token string, v Verbosity) (*Context, error) {
	c_access_token := C.CString(access_token)
	c_device := C.CString(dev)
	defer C.free(unsafe.Pointer(c_device))

	context, err := C.wit_init(nil, C.uint(v))
	if err != nil {
		switch e, ok := err.(syscall.Errno); {
		case !ok:
			return nil, err
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

func (c *Context) VoiceQueryAuto() (*Result, *Outcome, error) {
	c_result, err := C.wit_voice_query_auto(c.context, c.access_token)
	if err != nil {
		return nil, nil, err
	}
	defer C.free(unsafe.Pointer(c_result))

	return c.parseResult(C.GoString(c_result))
}

func (c *Context) VoiceQueryStart() error {
	_, err := C.wit_voice_query_start(c.context, c.access_token)
	return err
}

func (c *Context) VoiceQueryStop() (*Result, *Outcome, error) {
	c_result, err := C.wit_voice_query_stop(c.context)
	if err != nil {
		return nil, nil, err
	}
	defer C.free(unsafe.Pointer(c_result))

	return c.parseResult(C.GoString(c_result))
}

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
