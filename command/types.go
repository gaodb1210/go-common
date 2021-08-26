package command

import (
    "errors"
    "os"
)

const (
    Success int = iota
    Fail
    Timeout
    groupsIdentifier = "groups="
)

var (
    ErrCommandStart = errors.New("error occurred starting the command")
    ErrCommandTimeout = errors.New("command execute timeout")
)

type WaitProcessResult struct {
    processState *os.ProcessState
    err error
}