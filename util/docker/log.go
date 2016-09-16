package docker

import (
	"fmt"
)

// Log represents a single line from a running docker container
type Log struct {
	fmt.Stringer
	str string
}

func LogFromString(fmtStr string, vals ...interface{}) Log {
	return Log{str: fmt.Sprintf(fmtStr, vals...)}
}

func (l Log) Message() string {
	return l.str
}

func (l Log) String() string {
	return l.str
}
