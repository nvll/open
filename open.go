/*
 * Copyright (c) 2013, Chris Anderson
 * All rights reserved.
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 * list of conditions and the following disclaimer.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func main() {
	var editor, pathArg, nameArg, editorPath, openPath string
	var files []string
	var fileSel int

	if len(os.Args) < 2 {
		println("usage: ", os.Args[0], " <filename")
		os.Exit(1)
	}

	nameArg = os.Args[1]
	if len(os.Args) == 3 {
		pathArg = os.Args[2]
	} else {
		pathArg = "."
	}

	editor = os.Getenv("EDITOR")
	if editor == "" {
		println("No EDITOR env var set")
		os.Exit(1)
	}

	editorPath, err := exec.LookPath(editor);
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	// Find the source file, or ask for a choice if more than one exists
	callback := func(path string, info os.FileInfo, err error) error {
		if info.Name() == nameArg {
			files = append(files, path)
		} else if match, err := filepath.Match(nameArg, info.Name()); match && err == nil {
			files = append(files, path)
		}

		return nil
	}

	filepath.Walk(pathArg, callback)
	switch len(files) {
	case 0:
		println(nameArg, "not found")
		os.Exit(1)
	case 1:
		openPath = files[0]
	default:
		for i, fpath := range files {
			fmt.Printf("[%d]\t%s\n", i, fpath)
		}

	selection:
		print("? ")
		fmt.Scanf("%d", &fileSel)
		if fileSel > len(files) {
			println("invalid selection:", fileSel)
			goto selection
		}

		openPath = files[fileSel]
	}

	if err := syscall.Exec(editorPath, []string{editorPath, openPath}, os.Environ()); err != nil {
		println("Couldn't exec", err.Error())
		os.Exit(1)
	}

}

// vim: noet ts=4 sw=4 sts=4 tw=0
