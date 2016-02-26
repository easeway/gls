package gls

import (
    "strconv"
    "regexp"
    "runtime"
    "unsafe"
)

var StackBufferSize = 4096

var magicRex = func() *regexp.Regexp {
    rexStr := `\(0x[[:xdigit:]]+, 0x([[:xdigit:]]+)`
    for _, n := range magicNums {
        rexStr += ", 0x" + strconv.FormatUint(uint64(n), 16)
    }
    return regexp.MustCompile(rexStr + `\)`)
}()

// this is the containing structure for
// user data associated with context
type dataCntr struct {
    data interface{}
}

func findCntr() *dataCntr {
    stack := make([]byte, StackBufferSize)
    size := runtime.Stack(stack, false)
    n := magicRex.FindSubmatchIndex(stack[0:size])
    if len(n) < 4 {
        return nil
    }
    ptrStr := string(stack[n[2]:n[3]])
    if ptr, err := strconv.ParseUint(ptrStr, 16, 64); err != nil {
        return nil
    } else {
        return (*dataCntr)(unsafe.Pointer(uintptr(ptr)))
    }
}

func Get() interface{} {
    cntr := findCntr()
    if cntr == nil {
        panic("Unable to find GLS data, make sure gls.Go is called")
    }
    return cntr.data
}

func GetSafe() interface{} {
    cntr := findCntr()
    if cntr == nil {
        return nil
    }
    return cntr.data
}
