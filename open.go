package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func main() {
	var err error
	var editor string = os.Getenv("EDITOR")
	var editorPath, openPath string
	var files []string
	var fileSel int

	// Passed to os.Walk. Needs to be a Function Literal to have access
	// to any useful state
	callback := func(path string, info os.FileInfo, err error) error {
		if info.Name() == os.Args[1] {
			files = append(files, path)
		}

		return nil
	}

	if len(os.Args) < 2 {
		println("usage: ", os.Args[0], " <filename")
		os.Exit(1)
	}

	if editor == "" {
		print("No EDITOR env var set")
		os.Exit(1)
	}

	// The editor exists, right?
	if editorPath, err = exec.LookPath(editor); err != nil {
		println(err.Error())
		os.Exit(1)
	}

	// Find the source file, or ask for a choice if more than one exists
	filepath.Walk(".", callback)
	switch {
	case len(files) == 0:
		println(os.Args[1], "not found")
		os.Exit(1)
	case len(files) == 1:
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

	if err = syscall.Exec(editorPath, []string{editorPath, openPath}, os.Environ()); err != nil {
		println("Couldn't exec", err.Error())
		os.Exit(1)
	}

}

/* vim: set noexpandtab:ts=4:sw=4:sts=4 */
