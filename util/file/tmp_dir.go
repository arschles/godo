package file

import (
	"fmt"
	"io/ioutil"
	"os"
)

type TmpDirCreator func(prefixFmt string, vals ...interface{}) (string, error)

func DefaultTmpDirCreator() TmpDirCreator {
	return TmpDirCreator(func(prefixFmt string, vals ...interface{}) (string, error) {
		prefix := fmt.Sprintf(prefixFmt, vals...)
		return ioutil.TempDir("", prefix)
	})
}

func LocalTmpDirCreator() TmpDirCreator {
	return TmpDirCreator(func(prefixFmt string, vals ...interface{}) (string, error) {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		prefix := fmt.Sprintf(prefixFmt, vals...)
		return ioutil.TempDir(wd, prefix)
	})
}
