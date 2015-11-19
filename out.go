package main

import (
	"fmt"
	"os"

	"github.com/aybabtme/rgbterm"
)

const statusPrefix = "--->"

func format(fmtStr string, args []interface{}) string {
	if len(args) == 0 {
		return fmtStr
	}
	return fmt.Sprintf(fmtStr, args)
}

func errfAndExit(code int, fmtStr string, args ...interface{}) {
	s := format(fmtStr, args)
	s = rgbterm.String(s, 255, 0, 0, 0, 0, 0)
	fmt.Println(s)
	os.Exit(code)
}

func errf(fmtStr string, args ...interface{}) {
	s := fmt.Sprintf("%s %s", statusPrefix, format(fmtStr, args))
	s = rgbterm.String(s, 255, 0, 0, 0, 0, 0)
	fmt.Println(s)
}

func successf(fmtStr string, args ...interface{}) {
	s := format(fmtStr, args)
	s = rgbterm.String(s, 0, 255, 0, 0, 0, 0)
	fmt.Println(s)
}

func statusf(fmtStr string, args ...interface{}) {
	s := format(fmtStr, args)
	fmt.Printf("%s %s\n", statusPrefix, s)
}
