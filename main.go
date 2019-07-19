package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const cdbPath string = `C:\Program Files (x86)\Windows Kits\10\Debuggers\x64\cdb.exe`

func main() {

	var logFolder string
	var cdb string
	flag.StringVar(&logFolder, "d", "", "log folder contains DMP files")
	flag.StringVar(&cdb, "p", cdbPath, "cdb file path")
	flag.Parse()

	if logFolder == "" || cdb == "" {
		fmt.Fprintf(os.Stderr, "usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	if fileNotExist(cdb) {
		fmt.Fprintf(os.Stderr, "cdb not found.\n\tat location: %s", cdb)
		return
	}

	if fileNotExist(logFolder) {
		fmt.Fprintf(os.Stderr, "log folder not found.\n\tat location: %s", logFolder)
		return
	}

	files, err := ioutil.ReadDir(logFolder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.ToLower(filepath.Ext(file.Name())) == ".dmp" {
			dmpFile := path.Join(logFolder, file.Name())
			bugCheck := getBugCheckStr(cdb, dmpFile)
			fmt.Printf("%s\n\t%s\n\n", dmpFile, bugCheck)
		}
	}
}

func getBugCheckStr(cdb string, dmpFile string) string {
	output, err := exec.Command(cdb, "-z", dmpFile, `-c`, `!analyze -v;q`).Output()
	if err != nil {
		fmt.Println(err.Error())
	}

	result := string(output)
	scanner := bufio.NewScanner(strings.NewReader(result))
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "BUGCHECK_STR") {
			return scanner.Text()
		}
	}

	return "NOT FOUND"
}

func fileNotExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return os.IsNotExist(err)
}
