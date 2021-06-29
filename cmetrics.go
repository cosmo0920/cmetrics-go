package cmetrics

/*
#cgo LDFLAGS: -L/usr/local/lib -lcmetrics -lmpack -lxxhash
#cgo CFLAGS: -I/usr/local/include/ -w

#include <cmetrics/cmetrics.h>
#include <cmetrics/cmt_gauge.h>
#include <cmetrics/cmt_encode_prometheus.h>
#include <cmetrics/cmt_encode_msgpack.h>
#include <cmetrics/cmt_encode_text.h>
#include <cmetrics/cmt_counter.h>
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

type CMTContext struct {
	context *C.struct_cmt
}

type CMTGauge struct {
	gauge *C.struct_cmt_gauge
}

type CMTCounter struct {
	counter *C.struct_cmt_counter
}

func GoStringArrayToCptr(arr []string) **C.char {
	size := C.size_t(unsafe.Sizeof((*C.char)(nil)))
	length := C.size_t(len(arr))
	ptr := C.malloc(length * size)

	for i := 0; i < len(arr); i++ {
		element := (**C.char)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof((*C.char)(nil))))
		*element = C.CString(arr[i])
	}
	return (**C.char)(ptr)
}

func (g *CMTGauge) Add(ts time.Time, value float64, labels []string) error {
	ret := C.cmt_gauge_add(g.gauge, C.ulong(ts.UnixNano()), C.double(value), C.int(len(labels)), GoStringArrayToCptr(labels))
	if ret != 0 {
		return fmt.Errorf("cannot substract gauge value")
	}
	return nil
}

func (g *CMTGauge) Inc(ts time.Time, labels []string) error {
	ret := C.cmt_gauge_inc(g.gauge, C.ulong(ts.UnixNano()), C.int(len(labels)), GoStringArrayToCptr(labels))
	if ret != 0 {
		return fmt.Errorf("cannot increment gauge value")
	}
	return nil
}

func (g *CMTGauge) Dec(ts time.Time, labels []string) error {
	ret := C.cmt_gauge_dec(g.gauge, C.ulong(ts.UnixNano()), C.int(len(labels)), GoStringArrayToCptr(labels))
	if ret != 0 {
		return fmt.Errorf("cannot decrement gauge value")
	}
	return nil
}

func (g *CMTGauge) Sub(ts time.Time, value float64, labels []string) error {
	ret := C.cmt_gauge_sub(g.gauge, C.ulong(ts.UnixNano()), C.double(value), C.int(len(labels)), GoStringArrayToCptr(labels))
	if ret != 0 {
		return fmt.Errorf("cannot substract gauge value")
	}
	return nil
}

func (g *CMTGauge) GetVal(labels []string) (float64, error) {
	var value C.double
	ret := C.cmt_gauge_get_val(
		g.gauge,
		C.int(len(labels)),
		GoStringArrayToCptr(labels),
		&value)

	if ret != 0 {
		return -1, fmt.Errorf("cannot get value for gauge")
	}
	return float64(value), nil
}

func (g *CMTGauge) Set(ts time.Time, value float64, labels []string) error {
	ret := C.cmt_gauge_set(g.gauge, C.ulong(ts.UnixNano()), C.double(value), C.int(len(labels)), GoStringArrayToCptr(labels))
	if ret != 0 {
		return fmt.Errorf("cannot set gauge value")
	}
	return nil
}

func (ctx *CMTContext) EncodePrometheus() (string, error) {
	ret := C.cmt_encode_prometheus_create(ctx.context, 1)
	if ret == nil {
		return "", fmt.Errorf("error encoding to prometheus format")
	}
	return C.GoString(ret), nil
}

func (ctx *CMTContext) EncodeText() (string, error) {
	buffer := C.cmt_encode_text_create(ctx.context)
	if buffer == nil {
		return "", fmt.Errorf("error encoding to text format")
	}
	var text string = C.GoString(buffer)
	C.cmt_sds_destroy(buffer)
	return text, nil
}

func (ctx *CMTContext) EncodeMsgPack() (string, error) {
	var buffer string
	var cBuffer = C.CString(buffer)
	var size = C.size_t(len(buffer))
	ret := C.cmt_encode_msgpack(ctx.context, &cBuffer, (*C.size_t)(unsafe.Pointer(&size)))
	if ret != 0 {
		return "", fmt.Errorf("error encoding to msgpack format")
	}
	return C.GoString(cBuffer), nil
}

func (ctx *CMTContext) GaugeCreate(namespace, subsystem, name, help string, labelKeys []string) (*CMTGauge, error) {
	gauge := C.cmt_gauge_create(ctx.context,
		C.CString(namespace),
		C.CString(subsystem),
		C.CString(name),
		C.CString(help),
		C.int(len(labelKeys)),
		GoStringArrayToCptr(labelKeys),
	)
	if gauge == nil {
		return nil, fmt.Errorf("cannot create gauge")
	}
	return &CMTGauge{gauge}, nil
}

func (g *CMTCounter) Add(ts time.Time, value float64, labels []string) error {
	ret := C.cmt_counter_add(g.counter, C.ulong(ts.UnixNano()), C.double(value), C.int(len(labels)), GoStringArrayToCptr(labels))
	if ret != 0 {
		return fmt.Errorf("cannot substract counter value")
	}
	return nil
}

func (g *CMTCounter) Inc(ts time.Time, labels []string) error {
	ret := C.cmt_counter_inc(g.counter, C.ulong(ts.UnixNano()), C.int(len(labels)), GoStringArrayToCptr(labels))
	if ret != 0 {
		return fmt.Errorf("cannot Inc counter value")
	}
	return nil
}

func (g *CMTCounter) GetVal(labels []string) (float64, error) {
	var value C.double
	ret := C.cmt_counter_get_val(
		g.counter,
		C.int(len(labels)),
		GoStringArrayToCptr(labels),
		&value)

	if ret != 0 {
		return -1, fmt.Errorf("cannot get value for counter")
	}
	return float64(value), nil
}

func (g *CMTCounter) Set(ts time.Time, value float64, labels []string) error {
	ret := C.cmt_counter_set(g.counter, C.ulong(ts.UnixNano()), C.double(value), C.int(len(labels)), GoStringArrayToCptr(labels))
	if ret != 0 {
		return fmt.Errorf("cannot set counter value")
	}
	return nil
}

func (ctx *CMTContext) CounterCreate(namespace, subsystem, name, help string, labelKeys []string) (*CMTCounter, error) {
	counter := C.cmt_counter_create(ctx.context,
		C.CString(namespace),
		C.CString(subsystem),
		C.CString(name),
		C.CString(help),
		C.int(len(labelKeys)),
		GoStringArrayToCptr(labelKeys),
	)
	if counter == nil {
		return nil, fmt.Errorf("cannot create counter")
	}
	return &CMTCounter{counter}, nil
}

func (ctx *CMTContext) Destroy() {
	C.cmt_destroy(ctx.context)
}

func NewCMTContext() (*CMTContext, error) {
	cmt := C.cmt_create()
	if cmt == nil {
		return nil, fmt.Errorf("cannot create cmt context")
	}
	return &CMTContext{context: cmt}, nil
}
