// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	text, err := ioutil.ReadFile("../doc.go")
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(text), "\n")

	// Drop text up to and including opening comment.
	for i, line := range lines {
		if strings.HasPrefix(line, "/*") {
			lines = lines[i+1:]
			break
		}
	}

	// Drop leading blank line(s).
	for lines[0] == "" {
		lines = lines[1:]
	}

	// Drop text starting with trailing comment.
	for i, line := range lines {
		if strings.HasPrefix(line, "*/") {
			lines = lines[:i]
			break
		}
	}

	// Drop trailing blank line(s).
	for lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	buf := new(bytes.Buffer)
	s := func(s string) {
		buf.WriteString(s)
		buf.WriteByte('\n')
	}

	s(`// Code generated by "go generate robpike.io/ivy/parse"; DO NOT EDIT.`)
	s("")
	s("package parse")
	s("")
	s("var helpLines = []string{")
	for _, line := range lines {
		fmt.Fprintf(buf, "%q,\n", line)
	}
	s("}")

	s("")
	s("type helpIndexPair struct {")
	s("	start, end int")
	s("}")
	s("")

	// Pull out text for all the ops.
	// We know the field widths in runes.
	//	Roll              ?B    ?       One integer selected randomly from the first B integers

	s("var helpUnary = map[string]helpIndexPair{")
	var i int
	for i = 0; lines[i] != "Unary operators"; i++ {
	}
	for i++; i < len(lines); i++ {
		line := lines[i]
		if line == "Binary operators" {
			break
		}
		if len(line) < 33 {
			continue
		}
		if strings.Contains(line, "Name") && strings.Contains(line, "Meaning") { // It's a header.
			continue
		}
		// Find the op.
		runes := []rune(line)
		op := runes[25:]
		for i := 0; i < len(op); i++ {
			if op[i] == ' ' {
				op = op[:i]
				break
			}
		}
		if len(op) == 0 {
			continue
		}
		j := i
		isCircle := false
		if s := string(op); s == "sin" || s == "cos" || s == "tan" {
			isCircle = true
			op = []rune("sin")
			j += 2
		}
		fmt.Fprintf(buf, `%q: {%d, %d},`+"\n", string(op), i, j)
		if isCircle {
			fmt.Fprintf(buf, `%q: {%d, %d},`+"\n", "cos", i, j)
			fmt.Fprintf(buf, `%q: {%d, %d},`+"\n", "tan", i, j)
		}
		i = j
	}

	// Text-converters are all unary.
	for i = 0; lines[i] != "Type-converting operations"; i++ {
	}
	for i++; i < len(lines); i++ {
		line := lines[i]
		if line == "Pre-defined constants" {
			break
		}
		if len(line) < 33 {
			continue
		}
		if strings.Contains(line, "Name") && strings.Contains(line, "Meaning") { // It's a header.
			continue
		}
		// Find the op.
		runes := []rune(line)
		op := runes[25:]
		for i := 0; i < len(op); i++ {
			if op[i] == ' ' {
				op = op[:i]
				break
			}
		}
		if len(op) == 0 {
			continue
		}
		fmt.Fprintf(buf, `%q: {%d, %d},`+"\n", string(op), i, i)
	}

	s("}")
	s("")

	s("var helpBinary = map[string]helpIndexPair{")
	for i = 0; lines[i] != "Binary operators"; i++ {
	}
	for i++; i < len(lines); i++ {
		line := lines[i]
		if line == "Operators and axis indicator" {
			break
		}
		if len(line) < 33 {
			continue
		}
		if strings.Contains(line, "Name") && strings.Contains(line, "Meaning") { // It's a header.
			continue
		}
		// Find the op.
		runes := []rune(line)
		op := runes[29:]
		for i := 0; i < len(op); i++ {
			if op[i] == ' ' {
				op = op[:i]
				break
			}
		}
		if len(op) == 0 {
			continue
		}
		// Circles are unary.
		str := string(op)
		switch str {
		case "sin", "cos", "tan", "asin", "acos", "atan":
			continue
		}
		j := i
		// If the next few lines have no text at the left, they are a continuation. Pull them in.
		for ; j < len(lines); j++ {
			next := lines[j+1]
			if len(next) < 37 || next[1] != ' ' {
				break
			}
		}
		fmt.Fprintf(buf, `%q: {%d, %d},`+"\n", string(op), i, j)
		i = j
	}
	s("}")
	s("")

	for i = 0; lines[i] != "Operators and axis indicator"; i++ {
	}
	s("var helpAxis = map[string]helpIndexPair{")
	for i++; i < len(lines); i++ {
		line := lines[i]
		if line == "Type-converting operations" {
			break
		}
		if len(line) < 33 {
			continue
		}
		if strings.Contains(line, "Name") && strings.Contains(line, "Meaning") { // It's a header.
			continue
		}
		// Find the op.
		runes := []rune(line)
		op := runes[26:]
		for i := 0; i < len(op); i++ {
			if op[i] == ' ' {
				op = op[:i]
				break
			}
		}
		if len(op) == 0 {
			continue
		}
		fmt.Fprintf(buf, `%q: {%d, %d},`+"\n", string(op), i, i)
	}
	s("}")

	var formatted []byte
	if true {
		formatted, err = format.Source(buf.Bytes())
		if err != nil {
			log.Fatal(err)
		}
	} else {
		formatted = buf.Bytes()
	}

	fd, err := os.Create("help.go")
	if err != nil {
		log.Fatal(err)
	}
	_, err = fd.Write(formatted)
	if err != nil {
		log.Fatal("help.go")
	}
}
