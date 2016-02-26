// +build !386,!arm

package gls

var magicNums = []uint64{0xdb8111e56cf1657a, 0x8c325387ac2b1f39}

func With(data interface{}, fn func()) {
    d := &dataCntr{data}
    mark(fn, d, magicNums[0], magicNums[1])
}

func mark(fn func(), cntr *dataCntr, m1, m2 uint64) {
    fn()
}
