package gls

import (
    "crypto/rand"
    "encoding/hex"
    "testing"
)

func makeUniqueStr() string {
    raw := make([]byte, 16)
    rand.Read(raw)
    return hex.EncodeToString(raw)
}

func TestData(t *testing.T) {
    str := makeUniqueStr()
    With(&str, func() {
        if GetSafe() == nil {
            t.Fatal("GetSafe returns nil")
        }
        actual := Get().(*string)
        if actual != &str {
            t.Fatal("Get returns different pointer")
        }
        if *actual != str {
            t.Fatal("Get returns different value")
        }
    })
}

func TestNoData(t *testing.T) {
    actual := GetSafe()
    if actual != nil {
        t.Fatal("GetSafe returns non-nil")
    }
}

func TestNoDataPanic(t *testing.T) {
    recovered := false
    func() {
        defer func() {
            if r := recover(); r != nil {
                recovered = true
            }
        }()
        Get()
    }()
    if !recovered {
        t.Fatal("Get not panic without data")
    }
}

func TestFwGoRoutine(t *testing.T) {
    original := makeUniqueStr()
    With(original, func() {
        ch := make(chan string)
        Go(func() {
            ch <- Get().(string)
        })
        str := <- ch
        if str != original {
            t.Fatal("Get returns a different value")
        }
    })
}
