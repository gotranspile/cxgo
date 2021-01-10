package libc

import (
	"runtime"
	"strconv"
	"strings"
)

func Atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func Atof(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func FuncName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "<unknown func>"
	}
	fnc := runtime.FuncForPC(pc)
	if fnc == nil {
		return "<unknown func>"
	}
	name := fnc.Name()
	if i := strings.LastIndex(name, "."); i > 0 {
		name = name[i+1:]
	}
	return name
}
