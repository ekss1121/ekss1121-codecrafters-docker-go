package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	// Uncomment this block to pass the first stage!
	// "os"
	// "os/exec"
)

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage!
	//
	command := os.Args[3]
	args := os.Args[4:len(os.Args)]

	// create a temporary directory to chroot into
	chrootDir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(chrootDir)

	// copy the binary to run
	targetPath := filepath.Join(chrootDir, command)
	// fmt.Println("target path is ", targetPath)
	err = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm)

	if err != nil {
		panic(err)
	}

	bytesRead, err := ioutil.ReadFile(command)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(targetPath, bytesRead, os.ModePerm)

	err = syscall.Chroot(chrootDir)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			os.Exit(ee.ExitCode())
		} else {
			fmt.Printf("Err: %v\n", err)
			os.Exit(1)
		}
	}
}
