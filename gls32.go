// +build 386 arm

package gls

var magicNums = []uint32{0x6cf1657a, 0xdb8111e5, 0xac2b1f39, 0x8c325387}

func Go(data interface{}, fn func()) {
    d := &dataCntr{data}
    mark(fn, d, magicNums[0], magicNums[1], magicNums[2], magicNums[3])
}

func mark(fn func(), cntr *dataCntr, m1, m2, m3, m4 uint32) {
    fn()
}
