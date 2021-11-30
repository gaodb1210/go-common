package unit

import (
    "d7y.io/dragonfly/v2/pkg/util/stringutils"
    "encoding/json"
    "fmt"
    "github.com/pkg/errors"
    "gopkg.in/yaml.v3"
    "regexp"
    "strconv"
    "strings"
)

type Bytes int64

const (
    B  Bytes = 1
    KB       = 1024 * B
    MB       = 1024 * KB
    GB       = 1024 * MB
    TB       = 1024 * GB
    PB       = 1024 * TB
    EB       = 1024 * PB
)

var sizeRegexp = regexp.MustCompile(`^([0-9]+)(\.0*)?([MmKkGgTtPpEe])?[iI]?[bB]?$`)

func (f Bytes) ToNumber() int64 {
    return int64(f)
}

func ToBytes(size int64) Bytes {
    return Bytes(size)
}

// Set is used for command flag var
func (f *Bytes) Set(s string) (err error) {
    if stringutils.IsBlank(s) {
        *f = 0
    } else {
        *f, err = parseSize(s)
    }
    
    return
}

func (f Bytes) Type() string {
    return "bytes"
}

func (f Bytes) String() string {
    var (
        symbol string
        unit   Bytes
    )
    
    if f >= PB {
        symbol = "PB"
        unit = PB
    } else if f >= TB {
        symbol = "TB"
        unit = TB
    } else if f >= GB {
        symbol = "GB"
        unit = GB
    } else if f >= MB {
        symbol = "MB"
        unit = MB
    } else if f >= KB {
        symbol = "KB"
        unit = KB
    } else {
        symbol = "B"
        unit = B
    }
    
    return fmt.Sprintf("%.1f%s", float64(f)/float64(unit), symbol)
}

func parseSize(size string) (Bytes, error) {
    size = strings.TrimSpace(size)
    if stringutils.IsBlank(size) {
        return 0, nil
    }
    
    matches := sizeRegexp.FindStringSubmatch(size)
    if len(matches) == 0 {
        return 0, errors.Errorf("parse size %s: invalid format", size)
    }
    
    var unit Bytes
    switch matches[3] {
    case "k", "K":
        unit = KB
    case "m", "M":
        unit = MB
    case "g", "G":
        unit = GB
    case "t", "T":
        unit = TB
    case "p", "P":
        unit = PB
    case "e", "E":
        unit = EB
    default:
        unit = B
    }
    
    num, err := strconv.ParseInt(matches[1], 0, 64)
    if err != nil {
        return 0, errors.Wrapf(err, "failed to parse size: %s", size)
    }
    
    return ToBytes(num) * unit, nil
}

func (f Bytes) MarshalYAML() (interface{}, error) {
    result := f.String()
    return result, nil
}

func (f *Bytes) UnmarshalJSON(b []byte) error {
    return f.unmarshal(json.Unmarshal, b)
}

func (f *Bytes) UnmarshalYAML(node *yaml.Node) error {
    return f.unmarshal(yaml.Unmarshal, []byte(node.Value))
}

func (f *Bytes) unmarshal(unmarshal func(in []byte, out interface{}) (err error), b []byte) error {
    var v interface{}
    if err := unmarshal(b, &v); err != nil {
        return err
    }
    switch value := v.(type) {
    case float64:
        *f = Bytes(int64(value))
        return nil
    case int:
        *f = Bytes(int64(value))
        return nil
    case int64:
        *f = Bytes(value)
        return nil
    case string:
        size, err := parseSize(value)
        if err != nil {
            return errors.WithMessage(err, "invalid byte size")
        }
        *f = size
        return nil
    default:
        return errors.New("invalid byte size")
    }
}

