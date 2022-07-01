package single

// package single provides a mechanism to ensure, that only one instance of a program is running
import (
    "errors"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path"
    "path/filepath"
    "strconv"
    "syscall"
)

var (
    // ErrAlreadyRunning -- the instance is already running
    ErrAlreadyRunning = errors.New("the program is already running")
)

// Single represents the name and the open file descriptor
type Single struct {
    name string
    lockFile string
    pidFile string
    file *os.File
}

// New creates a Single instance
func New(name string, lockFile, pidFile string) *Single {
    return &Single{
        name: name,
        lockFile: lockFile,
        pidFile: pidFile,
    }
}

// Lock tries to obtain an exclude lock on a lockfile and exits the program if an error occurs
func (s *Single) Lock() {
    if err := s.CheckLock(); err != nil {
        log.Fatal(err)
    }
}

// Unlock releases the lock, closes and removes the lockfile. All errors will be reported directly.
func (s *Single) Unlock() {
    if err := s.TryUnlock(); err != nil {
        log.Print(err)
    }
}

// CheckLock tries to obtain an exclude lock on a lockfile and returns an error if one occurs
func (s *Single) CheckLock() error {

    // open/create lock file
    f, err := os.OpenFile(s.fileName(), os.O_RDWR|os.O_CREATE, 0600)
    if err != nil {
        return err
    }
    s.file = f
    // set the lock type to F_WRLCK, therefore the file has to be opened writable
    flock := syscall.Flock_t{
        Type: syscall.F_WRLCK,
        Pid:  int32(os.Getpid()),
    }
    // try to obtain an exclusive lock - FcntlFlock seems to be the portable *ix way
    if err := syscall.FcntlFlock(s.file.Fd(), syscall.F_SETLK, &flock); err != nil {
        return ErrAlreadyRunning
    }
    // write pid file
    if len(s.pidFile) == 0 {
        s.pidFile = path.Join(os.TempDir(), fmt.Sprintf("%s.pid", s.name))
    }
    var d1 = []byte(strconv.Itoa(os.Getpid()))
    _ = ioutil.WriteFile(s.pidFile, d1, 0666)
    return nil
}

// TryUnlock unlocks, closes and removes the lockfile
func (s *Single) TryUnlock() error {
    // set the lock type to F_UNLCK
    flock := syscall.Flock_t{
        Type: syscall.F_UNLCK,
        Pid:  int32(os.Getpid()),
    }
    if err := syscall.FcntlFlock(s.file.Fd(), syscall.F_SETLK, &flock); err != nil {
        return fmt.Errorf("failed to unlock the lock file: %v", err)
    }
    if err := s.file.Close(); err != nil {
        return fmt.Errorf("failed to close the lock file: %v", err)
    }
    if err := os.Remove(s.fileName()); err != nil {
        return fmt.Errorf("failed to remove the lock file: %v", err)
    }
    if s.pidFile != "" {
        _ = os.Remove(s.pidFile)
    }
    return nil
}

// fileName returns an absolute filename, appropriate for the operating system
func (s *Single) fileName() string {
    if len(s.lockFile) > 0 {
        return s.lockFile
    }
    return filepath.Join(os.TempDir(), fmt.Sprintf("%s.lock", s.name))
}
