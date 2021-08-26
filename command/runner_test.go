package command

import (
    "bytes"
    "testing"
)

func TestRunner_SyncRunSimple(t *testing.T) {
    r := newRunner()
    err := r.SyncRunSimple("sh", []string{"-c", "mkdir test_runner"}, 2)
    if err != nil {
        t.Error("command execute error:", err)
    }
}

func TestRunner_SyncRun(t *testing.T) {
    r := newRunner()
    r.SetUser("gaodb")
    r.SetPassword("123")
    output := bytes.NewBufferString("")
    r.SyncRun("", "sh", []string{"-c", "ls -al ./*"},
    output, output, 2)
    println(string(output.Bytes()))
}