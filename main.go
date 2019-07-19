package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func main() {

	var logFolder string
	flag.StringVar(&logFolder, "d", "", "log folder contains DMP files")
	flag.Parse()

	if logFolder == "" {
		fmt.Println("Usage:")
		return
	}

	files, err := ioutil.ReadDir(logFolder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.ToLower(filepath.Ext(file.Name())) == ".dmp" {
			dmpFile := path.Join(logFolder, file.Name())
			bugCheck := getBugCheckStr(dmpFile)
			fmt.Printf("%s\n\t%s\n\n", dmpFile, bugCheck)
		}
	}
}

func getBugCheckStr(dmpFile string) string {
	output, err := exec.Command("cdb", "-z", dmpFile, `-c`, `!analyze -v;q`).Output()
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

	return ""
}
