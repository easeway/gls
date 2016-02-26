package gls

import (
    "crypto/rand"
    "encoding/hex"
    "testing"
)

func TestData(t *testing.T) {
    raw := make([]byte, 16)
    rand.Read(raw)
    str := hex.EncodeToString(raw)
    Go(&str, func() {
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
