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
const cdbCommand string = "!analyze -v;q"
const version = "v0.0.2"

const ok = 0
const dmpNotFound = 1
const cdbNotFound = 2

func main() {

	var logFolder string
	var cdb string
	var dump string
	var cdbCmd string

	flag.StringVar(&logFolder, "d", "", "log folder contains DMP files")
	flag.StringVar(&cdb, "p", cdbPath, "cdb file path")
	flag.StringVar(&dump, "f", "", "analyze specific dump file, ignore -d if flag set")
	flag.StringVar(&cdbCmd, "c", cdbCommand, "command issued tocdb debugger")
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
		bugCheck := analyze(cdb, dump, cdbCmd)
		fmt.Printf("%s\n\t%s\n\n", dump, bugCheck)
		os.Exit(ok)
	}

	//get bugcheck str for all dump files in folder
	if !FileExist(logFolder) {
		fmt.Fprintf(os.Stderr, "log folder not found.\n\tat location: %s", logFolder)
		os.Exit(dmpNotFound)
	}

	files, err := ioutil.ReadDir(logFolder)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.ToLower(filepath.Ext(file.Name())) == ".dmp" {
			dmpFile := path.Join(logFolder, file.Name())
			bugCheck := analyze(cdb, dmpFile, cdbCmd)
			fmt.Printf("%s\n\t%s\n\n", dmpFile, bugCheck)
		}
	}
}

func printVersion() {
	fmt.Println(version)
}

func executeCdb(cdb string, dmpFile string, cdbCmd string) ([]byte, error) {
	output, err := exec.Command(cdb, "-z", dmpFile, "-c", cdbCmd).Output()
	return output, err
}

func analyze(cdb string, dmpFile string, cdbCmd string) string {
	output, err := executeCdb(cdb, dmpFile, cdbCmd)
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
