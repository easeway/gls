package gls

import (
    "crypto/rand"
    "encoding/hex"
    "testing"
)

func TestContext(t *testing.T) {
    raw := make([]byte, 16)
    rand.Read(raw)
    str := hex.EncodeToString(raw)
    WithCtx(&str, func() {
        actual := Ctx().(*string)
        if actual != &str {
            t.Fail()
        }
        if *actual != str {
            t.Fail()
        }
    })
}

func TestNoContext(t *testing.T) {
    actual := Ctx()
    if actual != nil {
        t.Fail()
    }
}
