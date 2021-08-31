package logging

import (
    "bytes"
    "fmt"
    "github.com/sirupsen/logrus"
    "runtime"
    "sort"
    "strings"
)

const defaultTimestampFormat = "2006-01-02 15:04:05"

// Formatter - implements logrus.Formatter
type Formatter struct {
    // FieldsOrder - default: fields sorted alphabetically
    FieldsOrder []string
    
    // TimestampFormat - default: time.StampMilli = "Jan _2 15:04:05.000"
    TimestampFormat string
    
    // Separator - default: space,
    Separator string
    
    // ShowFullLevel - show a full level [WARNING] instead of [WARN]
    ShowFullLevel bool
    
    // NoUppercaseLevel - no upper case for level value
    NoUppercaseLevel bool
    
    // TrimMessages - trim whitespaces on messages
    TrimMessages bool
    
    // NoFunction - do not print function
    NoFunction bool

    // CustomCallerFormatter - set custom formatter for caller info
    CustomCallerFormatter func(*runtime.Frame) string
}

// Format an log entry
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
    
    // output buffer
    b := &bytes.Buffer{}

    if f.Separator == "" {
        f.Separator = "  "
    }

    // write timestamp
    f.writeTimestamp(b, entry)
    // write level
    f.writeLevel(b, entry)
    // write caller
    f.writeCaller(b, entry)
    
    // write message
    if f.TrimMessages {
        b.WriteString(strings.TrimSpace(entry.Message))
    } else {
        b.WriteString(entry.Message)
    }
    b.WriteString(f.Separator)

    // write fields
    if f.FieldsOrder == nil {
        f.writeFields(b, entry)
    } else {
        f.writeOrderedFields(b, entry)
    }

    b.WriteByte('\n')
    
    return b.Bytes(), nil
}

func (f *Formatter) writeTimestamp(b *bytes.Buffer, entry *logrus.Entry) {
    timestampFormat := f.TimestampFormat
    if timestampFormat == "" {
        timestampFormat = defaultTimestampFormat
    }
    // write timestamp
    b.WriteString(entry.Time.Format(timestampFormat))
    b.WriteString(f.Separator)
}

func (f *Formatter) writeLevel(b *bytes.Buffer, entry *logrus.Entry)  {
    // write level
    var level string
    if f.NoUppercaseLevel {
        level = entry.Level.String()
    } else {
        level = strings.ToUpper(entry.Level.String())
    }

    if f.ShowFullLevel {
        b.WriteString(level)
    } else {
        b.WriteString(level[:4])
    }

    b.WriteString(f.Separator)
}

func (f *Formatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
    if entry.HasCaller() {
        if f.CustomCallerFormatter != nil {
            _, _ = fmt.Fprintf(b, f.CustomCallerFormatter(entry.Caller))
        } else {
            if f.NoFunction {
                _, _ = fmt.Fprintf(
                    b,
                    " %s:%d",
                    entry.Caller.File,
                    entry.Caller.Line,
                )
            } else {
                _, _ = fmt.Fprintf(
                    b,
                    " %s:%d %s",
                    entry.Caller.File,
                    entry.Caller.Line,
                    entry.Caller.Function,
                )
            }
            
        }
        b.WriteString(f.Separator)
    }
}

func (f *Formatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
    if len(entry.Data) != 0 {
        fields := make([]string, 0, len(entry.Data))
        for field := range entry.Data {
            fields = append(fields, field)
        }
        
        sort.Strings(fields)
        
        for _, field := range fields {
            f.writeField(b, entry, field)
        }
    }
}

func (f *Formatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
    length := len(entry.Data)
    foundFieldsMap := map[string]bool{}
    for _, field := range f.FieldsOrder {
        if _, ok := entry.Data[field]; ok {
            foundFieldsMap[field] = true
            length--
            f.writeField(b, entry, field)
        }
    }
    
    if length > 0 {
        notFoundFields := make([]string, 0, length)
        for field := range entry.Data {
            if foundFieldsMap[field] == false {
                notFoundFields = append(notFoundFields, field)
            }
        }
        
        sort.Strings(notFoundFields)
        
        for _, field := range notFoundFields {
            f.writeField(b, entry, field)
        }
    }
}

func (f *Formatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
    _, _ = fmt.Fprintf(b, "[%s=%v]", field, entry.Data[field])
    b.WriteString(" ")
}
