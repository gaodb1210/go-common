package fileutil

import (
    "fmt"
    "github.com/gaodb1210/go-common/util/unit"
    "github.com/pkg/errors"
    "io"
    "os"
    "path/filepath"
    "syscall"
)

const (
    // PrivateDirMode means read and execute access for everyone and also write access for the owner of the directory.
    PrivateDirMode = 0755
)

// MkdirAll creates a directory named path with 0755 perm.
func MakeDirAll(dir string) error {
    return os.MkdirAll(dir, PrivateDirMode)
}

// DeleteFile deletes a regular file not a directory.
func DeleteFile(path string) error {
    if PathExist(path) {
        if IsDir(path) {
            return errors.Errorf("delete file %s: not a regular file", path)
        }
        return os.Remove(path)
    }
    
    return nil
}

// OpenFile opens a file. If the parent directory of the file isn't exist,
// it will create the directory with 0755 perm.
func OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
    if !PathExist(path) && (flag&syscall.O_CREAT != 0) {
        if err := MakeDirAll(filepath.Dir(path)); err != nil {
            return nil, errors.Wrapf(err, "open file %s", path)
        }
    }
    
    return os.OpenFile(path, flag, perm)
}

// Link creates a hard link pointing to oldName named newName for a file.
func Link(oldName string, newName string) error {
    if PathExist(newName) {
        if IsDir(newName) {
            return errors.Errorf("link %s to %s: src already exists and is a directory", newName, oldName)
        }
        if err := DeleteFile(newName); err != nil {
            return errors.Wrapf(err, "to link %s to %s", newName, oldName)
        }
    } else if err := MakeDirAll(filepath.Dir(newName)); err != nil {
        return errors.Wrapf(err, "link %s to %s", newName, oldName)
    }
    
    return os.Link(oldName, newName)
}

func IsSymbolicLink(name string) bool {
    f, e := os.Lstat(name)
    if e != nil {
        return false
    }
    return f.Mode()&os.ModeSymlink != 0
}

// SymbolicLink creates newname as a symbolic link to oldname.
func SymbolicLink(oldname string, newname string) error {
    if PathExist(newname) && IsSymbolicLink(newname) && PathExist(oldname) {
        dstPath, err := os.Readlink(newname)
        if err != nil {
            return err
        }
        if dstPath == oldname {
            return nil
        }
    }
    if PathExist(newname) {
        if err := os.Remove(newname); err != nil {
            return fmt.Errorf("symlink %s to %s when deleting target file: %v", newname, oldname, err)
        }
    }
    
    if err := MakeDirAll(filepath.Dir(newname)); err != nil {
        return errors.Wrapf(err, "symlink %s to %s", newname, oldname)
    }
    return os.Symlink(oldname, newname)
}

// PathExist reports whether the path is exist.
// Any error, from os.Lstat, will return false.
func PathExist(path string) bool {
    _, err := os.Lstat(path)
    return err == nil
}

func IsDir(path string) bool {
    f, err := stat(path)
    if err != nil {
        return false
    }
    
    return f.IsDir()
}

func IsRegular(path string) bool {
    f, err := stat(path)
    if err != nil {
        return false
    }
    
    return f.Mode().IsRegular()
}

func IsEmptyDir(path string) (bool, error) {
    f, err := os.Open(path)
    if err != nil {
        return false, err
    }
    defer f.Close()
    
    if _, err = f.Readdirnames(1); err == io.EOF {
        return true, nil
    }
    
    return false, err
}

func stat(path string) (os.FileInfo, error) {
    f, err := os.Stat(path)
    return f, err
}

// GetFreeSpace gets the free disk space of the path.
func GetFreeSpace(path string) (unit.Bytes, error) {
    fs := syscall.Statfs_t{}
    if err := syscall.Statfs(path, &fs); err != nil {
        return 0, err
    }
    
    return unit.Bytes(fs.Bavail * uint64(fs.Bsize)), nil
}

// GetTotalSpace
func GetTotalSpace(path string) (unit.Bytes, error) {
    fs := syscall.Statfs_t{}
    if err := syscall.Statfs(path, &fs); err != nil {
        return 0, err
    }
    
    return unit.Bytes(fs.Blocks * uint64(fs.Bsize)), nil
}

func GetTotalAndFreeSpace(path string) (unit.Bytes, unit.Bytes, error) {
    fs := syscall.Statfs_t{}
    if err := syscall.Statfs(path, &fs); err != nil {
        return 0, 0, err
    }
    total := unit.Bytes(fs.Blocks * uint64(fs.Bsize))
    free := unit.Bytes(fs.Bavail * uint64(fs.Bsize))
    return total, free, nil
}

// GetUsedSpace
func GetUsedSpace(path string) (unit.Bytes, error) {
    fs := syscall.Statfs_t{}
    if err := syscall.Statfs(path, &fs); err != nil {
        return 0, err
    }
    return unit.Bytes((fs.Blocks - fs.Bavail) * uint64(fs.Bsize)), nil
}

// MoveFile moves the file src to dst.
func MoveFile(src string, dst string) error {
    if !IsRegular(src) {
        return fmt.Errorf("move %s to %s: src is not a regular file", src, dst)
    }
    if PathExist(dst) && !IsDir(dst) {
        if err := DeleteFile(dst); err != nil {
            return fmt.Errorf("move %s to %s when deleting dst file: %v", src, dst, err)
        }
    }
    return os.Rename(src, dst)
}
