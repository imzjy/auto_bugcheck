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
	"regexp"
	"strings"
)

const cdbPath string = `C:\Program Files (x86)\Windows Kits\10\Debuggers\x64\cdb.exe`
const cdbCommand string = "!analyze -v;q"
const exactRegex = "^BUGCHECK_STR"
const version = "v0.0.3"

const ok = 0
const dmpNotFound = 1
const cdbNotFound = 2

func main() {

	var logFolder string
	var cdb string
	var dump string
	var cdbCmd string
	var regPattern string

	flag.StringVar(&logFolder, "d", "", "folder contains DMP files")
	flag.StringVar(&cdb, "p", cdbPath, "cdb file path")
	flag.StringVar(&dump, "f", "", "analyze specific dump file, ignore -d if flag set")
	flag.StringVar(&cdbCmd, "c", cdbCommand, "command issued tocdb debugger")
	flag.StringVar(&regPattern, "regex", exactRegex, "regular express to exact from cdb output")
	ver := flag.Bool("version", false, "print version")
	flag.Parse()

	if *ver {
		printVersion()
		os.Exit(ok)
	}

	if (logFolder == "" && dump == "") || cdb == "" || cdbCmd == "" {
		fmt.Fprintf(os.Stderr, "%s %s \n", filepath.Base(os.Args[0]), version)
		flag.PrintDefaults()
		os.Exit(dmpNotFound)
	}

	if !FileExist(cdb) {
		fmt.Fprintf(os.Stderr, "cdb not found.\n\tat location: %s", cdb)
		os.Exit(cdbNotFound)
	}

	//specific dump analyze
	if FileExist(dump) {
		prettyPrintMatched(dump, analyze(cdb, dump, cdbCmd, regPattern))
		os.Exit(ok)
	}

	//get bugcheck str for all dump files in folder
	if !FileExist(logFolder) {
		fmt.Fprintf(os.Stderr, "log folder not found.\n\tat location: %s", logFolder)
		os.Exit(dmpNotFound)
	}
	processDumpFolder(logFolder, cdb, cdbCmd, regPattern)

}

func printVersion() {
	fmt.Println(version)
}

func executeCdb(cdb string, dmpFile string, cdbCmd string) ([]byte, error) {
	output, err := exec.Command(cdb, "-z", dmpFile, "-c", cdbCmd).Output()
	return output, err
}

func analyze(cdb string, dmpFile string, cdbCmd string, regPattern string) string {
	output, err := executeCdb(cdb, dmpFile, cdbCmd)
	if err != nil {
		fmt.Println(err.Error())
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {

		line := scanner.Text()
		matched, err := regexp.MatchString(regPattern, line)

		if err != nil {
			fmt.Fprintf(os.Stderr, "regex syntax error: %s", err.Error())
			return "ERROR"
		}

		if matched {
			return line
		}
	}

	return "NOT FOUND"
}

func prettyPrintMatched(dumpFile string, matchedStr string) {
	dumpFullPath, err := filepath.Abs(dumpFile)
	if err != nil {
		dumpFullPath = dumpFile
	}

	fmt.Printf("%s\n\t%s\n\n", dumpFullPath, matchedStr)
}

func processDumpFolder(logFolder string, cdb string, cdbCmd string, regPattern string) {

	files, err := ioutil.ReadDir(logFolder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.ToLower(filepath.Ext(file.Name())) == ".dmp" {
			dumpFile := path.Join(logFolder, file.Name())
			prettyPrintMatched(dumpFile, analyze(cdb, dumpFile, cdbCmd, regPattern))
		}
	}
}
