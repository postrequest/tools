package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// colours
var (
	printGreen = "\033[32;1m"
	printBlue  = "\033[34;1m"
	resetTerm  = "\033[0m"
)

// pre file indent
func indent(depth int) string {
	output := ""
	for i := 0; i < depth; i++ {
		output += "   "
	}
	return output + "|-->"
}

func walk(path string, verbose bool) {
	var temp string
	if _, err := os.Lstat(path); err != nil {
		usage()
	}
	fmt.Println(printGreen + path + resetTerm)
	err := filepath.Walk(path, func(currentPath string, info os.FileInfo, err error) error {
		if path == currentPath {
		} else if info.IsDir() {
			temp = strings.ReplaceAll(currentPath, path, "")
			index := strings.LastIndex(currentPath, "/")
			if index == -1 {
				index = 0
			}
			depth := strings.Count(temp, "/")
			fmt.Println(indent(depth), printGreen+currentPath[index:]+resetTerm)
		}
		if verbose && !info.IsDir() {
			temp = strings.ReplaceAll(currentPath, path, "")
			index := strings.LastIndex(currentPath, "/")
			if index == -1 {
				index = 0
			}
			depth := strings.Count(temp, "/")
			fmt.Println(indent(depth), printBlue+currentPath[index+1:]+resetTerm)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func usage() {
	fmt.Println("Usage:", os.Args[0], "[options] <dir>")
	fmt.Println(" ", os.Args[0], "<dir>")
	fmt.Println(" ", os.Args[0], "-v <dir>")
	fmt.Println("Options:\n  -v  Print files as well")
	os.Exit(1)
}

func main() {
	if len(os.Args) == 1 {
		walk(".", false)
	} else if len(os.Args) == 2 {
		walk(os.Args[1], false)
	} else if len(os.Args) == 3 {
		if os.Args[1] == "-v" {
			walk(os.Args[2], true)
		} else {
			usage()
		}
	} else {
		usage()
	}
}
