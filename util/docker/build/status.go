package build

import (
	"fmt"
)

type Log struct {
	str string
}

func LogFromString(fmtStr string, vals ...interface{}) Log {
	return Log{str: fmt.Sprintf(fmtStr, vals...)}
}

func (l Log) Message() string {
	return l.str
}
