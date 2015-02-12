package main

import (
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var (
	OptVerbose bool
	OptDir     string
	Args       []string
)

// Certain git commands are marked as "silent" and will only display "OK" if they suceed.
// This is to stop the user from being overwhelmed with information when running basic commands
// This behavior can be overriden by passing the "-v" switch
var silentCommends []string = []string{
	"add",
	"checkout",
	"commit",
	"push",
	"fetch",
	"merge",
	"mv",
	"pull",
	"push",
	"rebase",
	"reset",
	"rm",
	"tag",
}

func Usage() {
	fmt.Println("multigit - run git commands against many git repositories at once")
	os.Exit(0)
}

func main() {
	spew.Config.DisableMethods = true
	flag.BoolVar(&OptVerbose, "verbose", false, "be verbose")
	flag.StringVar(&OptDir, "dir", "", "directory that contains your git repositories")
	flag.Usage = Usage
	flag.Parse()

	Args := flag.Args()
	if len(Args) == 0 {
		Usage()
		os.Exit(0)
	}

	if OptDir != "" {
		err := os.Chdir(OptDir)
		if err != nil {
			log.Fatal(err)
		}
	}

	// List all directories and determine which ones are a git directory
	basedir, err := os.Open(".")
	if err != nil {
		log.Fatal(err)
	}
	subdirs, err := basedir.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}

	gitdirs := make([]string, 0)
	maxnamelen := 0
	for _, fi := range subdirs {
		if fi.IsDir() {
			_, err := os.Stat(fi.Name() + "/.git")
			if err == nil {
				gitdirs = append(gitdirs, fi.Name())
				if len(fi.Name()) > maxnamelen {
					maxnamelen = len(fi.Name())
				}
			}
		}
	}

	// Run the commands
	var wg sync.WaitGroup
	results := make([]string, len(gitdirs))
	errmap := make(map[string]error)
	numerror := 0
	for i, dir := range gitdirs {
		wg.Add(1)
		go func(i int, dir string) {
			cmd := exec.Command("git", Args...)
			cmd.Dir = dir
			output, err := cmd.CombinedOutput()
			if errmap[dir] != nil {
				numerror++
			}
			errmap[dir] = err
			results[i] = string(output)
			wg.Done()
		}(i, dir)
	}
	wg.Wait()

	// Report the results
	for i, dir := range gitdirs {
		fmt.Print(dir + " ")
		if errmap[dir] == nil {
			fmt.Println("OK")
		} else {
			fmt.Println("ERROR")
		}

		out := strings.Split(strings.Trim(results[i], "\n"), "\n")
		for _, line := range out {
			fmt.Println("    " + line)
		}
		fmt.Println("")
	}

	// Exit with correct code
	if numerror == len(gitdirs) {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
