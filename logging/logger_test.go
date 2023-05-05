package logging

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestInitLogger 使用zap实现的logger
func TestInitLogger(t *testing.T) {
	err := InitLogger("./test.log", "debug", 1, 7, false)
	assert.Equal(t, nil, err)
	Debugf("this is a test, level = %s", "DEBUG")
	Infof("this is a test, level = %s", "INFO")
	Warnf("this is a test, level = %s", "WARNING")
	//log.Errorf("this is a test, level = %s", "ERROR")
}
