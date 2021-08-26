package command

import (
    "fmt"
    "io"
    "os/exec"
    "time"
)

type Runner struct {
    command         *exec.Cmd
    user            string
    password        string
    homeDir         string
    canceledChan    chan bool
}

func newRunner() *Runner {
    return &Runner{}
}

// Cancel cancel running command
func (r *Runner) Cancel() {
    if r.command != nil {
        _ = r.command.Process.Kill()
    }
}

// SetUser set user
func (r *Runner) SetUser(name string) {
    r.user = name
}

// SetPassword set password
func (r *Runner) SetPassword(password string) {
    r.password = password
}

// SetHomeDir set home dir
func (r *Runner) SetHomeDir(homeDir string) {
    r.homeDir = homeDir
}

// SyncRunSimple sync run command, ignore output
func (r *Runner)  SyncRunSimple(commandName string, commandArguments []string, timeOut int) error {
    
    // 1. init command
    r.command = exec.Command(commandName, commandArguments...)
    if err := r.preProcess(); err != nil {
        return err
    }
    // 2.start command
    if err := r.command.Start(); err != nil {
        fmt.Printf("start command error:%s\n", err.Error())
        return ErrCommandStart
    }
    // 3. create goroutine to wait command finish
    finished := make(chan error, 1)
    go func() {
        finished <- r.command.Wait()
    }()
    // 4. wait command finish or timeout
    var err error
    select {
    case err = <-finished:
        fmt.Printf("command execute completed.")
        if err != nil {
            fmt.Printf("command execute error: %s", err.Error())
        }
    case <-time.After(time.Duration(timeOut) * time.Second):
        fmt.Printf("command execute timeout.")
        _ = r.command.Process.Kill()
        err = ErrCommandTimeout
    }
    return err
}

// SyncRun sync run command, write stdout to stdoutWriter, write stderr to stderrWriter
func (r *Runner) SyncRun(
    workingDir string,
    commandName string,
    commandArguments []string,
    stdoutWriter io.Writer,
    stderrWriter io.Writer,
    timeOut int) (exitCode int, status int, err error) {
    
    status = Success
    exitCode = 0
    // 1. init command
    r.command = exec.Command(commandName, commandArguments...)
    r.command.Stdout = stdoutWriter
    r.command.Stderr = stderrWriter
    r.command.Dir = workingDir
    if err := r.preProcess(); err != nil {
        return 0, Fail, err
    }

    // 3. start command
    if err = r.command.Start(); err != nil {
        fmt.Printf("start command fail: %s\n", err)
        exitCode = 1
        return exitCode, Fail, err
    }
    // 4. start goroutine to wait finish
    finished := make(chan WaitProcessResult, 1)
    go func() {
        processState, err := r.command.Process.Wait()
        finished <- WaitProcessResult{
            processState: processState,
            err: err,
        }
    }()
    // 5. wait command execute finish or timeout
    select {
    case waitProcessResult := <-finished:
        fmt.Printf("Command: %s execute completed\n", commandName)
        if waitProcessResult.processState != nil {
            if waitProcessResult.err != nil {
                fmt.Printf("os.Process.Wait() returns error with valid process state\n")
            }
            
            exitCode = waitProcessResult.processState.ExitCode()
            // Sleep 200ms to allow remaining data to be copied back
            time.Sleep(time.Duration(200) * time.Millisecond)
            // Explicitly break select statement in case timer also times out
            break
        } else {
            exitCode = 1
            return exitCode, Fail, waitProcessResult.err
        }
    case <-time.After(time.Duration(timeOut) * time.Second):
        fmt.Printf("command: %s execute timeout", commandName)
        exitCode = 1
        status = Timeout
        err = ErrCommandTimeout
        _ = r.command.Process.Kill()
    }
    
    if r.user != "" {
        _ = r.removeCredential()
    }
    
    return exitCode, status, err
}