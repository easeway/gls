package gls

import (
    "bytes"
    "strconv"
    "strings"
    "runtime"
    "unsafe"
)

const (
	magicFn = "._6cf1657a_db81_11e5_ac2b_1f398c325387"
)

var (
    StackBufferSize = 4096
)

type context struct {
    userContext interface{}
}

func WithCtx(ctx interface{}, fn func()) {
    containingCtx := &context{ctx}
    _6cf1657a_db81_11e5_ac2b_1f398c325387(fn, containingCtx)
}

func Ctx() interface{} {
    stackBuf := make([]byte, StackBufferSize)
    actualSize := runtime.Stack(stackBuf, false)
    buf := bytes.NewBuffer(stackBuf[0:actualSize])
    var err error = nil
    var line string
    for err == nil {
        line, err = buf.ReadString('\n')
        line = strings.TrimSpace(line)
        if !strings.HasSuffix(line, ")") {
			continue
		}
		pos := strings.LastIndexByte(line, '(')
        start := pos - len(magicFn)
		if start <= 0 {
			continue
		}
		if line[start:pos] == magicFn {
			if ptrs := strings.Split(line[pos+1:len(line)-1], ","); len(ptrs) != 2 {
				continue
			} else if ptr, err := strconv.ParseUint(strings.TrimSpace(ptrs[1]), 0, 64); err != nil {
				continue
			} else {
				return (*context)(unsafe.Pointer(uintptr(ptr))).userContext
			}
		}
    }
    return nil
}

func _6cf1657a_db81_11e5_ac2b_1f398c325387(fn func(), ctx *context) {
    fn()
}
