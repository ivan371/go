package main

import (
	"fmt"
	"io"
	"os"
	"io/ioutil"
	"strconv"
)

func findLastDirIndex(files []os.FileInfo) int {
	var lastDirIndex = 0
	for index, file := range files {
		if (file.IsDir()) {
			lastDirIndex = index
		}
	}

	return lastDirIndex
}

func getIsLastFileObject(files []os.FileInfo, printFiles bool, lastDirIndex int, index int) bool {
	return index == len(files) - 1 || !printFiles && lastDirIndex <= index
}

func printDir(out io.Writer, dir os.FileInfo, path string, printFiles bool, isEnd bool, oldPrint string) {
	if (!dir.IsDir() && !printFiles) {
		return
	}
	
	var endSymbol string

	if (isEnd) {
		endSymbol = "└"
	} else {
		endSymbol = "├"
	}

	var fullString = oldPrint + endSymbol + "───" + dir.Name()
	
	if (!dir.IsDir()) {
		var fileSize = dir.Size()
		var fileName string
		if (fileSize == 0) {
			fileName = " (empty)"
		} else {
			fileName = " (" +  strconv.FormatInt(fileSize, 10) + "b)"
		}
		fmt.Fprintln(out, fullString + fileName)
	} else {
		fmt.Fprintln(out, fullString)
		var newPath = path + "/" + dir.Name()
		var tabSymbol string

		if (isEnd) {
			tabSymbol = ""
		} else {
			tabSymbol = "│"
		}
		
		readDir(out, newPath, printFiles, oldPrint + tabSymbol + "\t")
	}
}

func readDir(out io.Writer, path string, printFiles bool, print string) error {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		panic("Ошибка при чтении переданной директории")
	}

	var lastDirIndex = findLastDirIndex(files)

	for index, file := range files {
		printDir(out, file, path, printFiles, getIsLastFileObject(files, printFiles, lastDirIndex, index), print)
	}

	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return readDir(out, path, printFiles, "")
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
