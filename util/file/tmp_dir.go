package file

import (
	"fmt"
	"io/ioutil"
)

type TmpDirCreator func(loc, prefixFmt string, vals ...interface{}) (string, error)

func DefaultTmpDirCreator() TmpDirCreator {
	return TmpDirCreator(func(location, prefixFmt string, vals ...interface{}) (string, error) {
		prefix := fmt.Sprintf(prefixFmt, vals...)
		return ioutil.TempDir(location, prefix)
	})
}
