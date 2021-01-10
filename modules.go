package cxgo

import "bytes"

type SrcFunc struct {
	File        string `json:"file,omitempty"`
	Line        int    `json:"line,omitempty"`
	OffsetStart int    `json:"offset_start,omitempty"`
	OffsetEnd   int    `json:"offset_end,omitempty"`
	Func        string `json:"func"`
	Src         string `json:"src,omitempty"`
	Proto       string `json:"proto,omitempty"`
}

func SourceFunc(fname string, src []byte, fd *CFuncDecl) SrcFunc {
	fsrc := ""
	psrc := ""
	var (
		istart int
		iend   int
		line   int
	)
	if rng := fd.Range; rng != nil && rng.Start > 0 {
		line = rng.StartLine
		start := rng.Start
		end := rng.End
		if e, f := findEnd(src[start:]); e > 0 {
			e += start
			f += start
			end = e
			psrc = string(bytes.TrimSpace(src[start:f])) + ";\n"
		}
		if s := findCommentStart(src[:start]); s > 0 {
			start = s
		}
		if end > 0 {
			istart = start
			iend = end
			fsrc = string(src[start:end])
		}
	}
	return SrcFunc{
		File:        fname,
		Func:        fd.Name.Name,
		Src:         fsrc,
		Proto:       psrc,
		Line:        line,
		OffsetStart: istart,
		OffsetEnd:   iend,
	}
}

func findCommentStart(src []byte) int {
	if len(src) == 0 || src[len(src)-1] != '\n' {
		return -1
	}
	src = src[:len(src)-1]
	i := bytes.LastIndex(src, []byte("//"))
	if i < 0 || bytes.ContainsAny(src[i:], "\n\r") {
		return -1
	}
	return i
}

func findEnd(src []byte) (int, int) {
	started := false
	level := 0
	first := -1
	i := 0
	for ; (!started || level > 0) && i < len(src); i++ {
		switch src[i] {
		case '{':
			if !started {
				started = true
				first = i
			}
			level++
		case '}':
			level--
			if level == 0 {
				break
			}
		}
	}
	if i >= len(src) {
		return -1, -1
	}
	for j := 0; j < 2; j++ {
		if i+1 < len(src) && src[i+1] == '\n' {
			i++
		}
	}
	return i, first
}
