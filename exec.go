package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

func run(pwd string, filename string, cmd string) string {
	c := exec.Command("bash", "-c", cmd)
	return runCommand(c, pwd, filename)
}

func runstdin(pwd string, filename string, stdin []byte) string {
	c := exec.Command("bash")
	c.Stdin = bytes.NewReader(stdin)
	return runCommand(c, pwd, filename)
}

func runCommand(c *exec.Cmd, pwd string, filename string) string {
	// todo: this is a simple hack to ensure the scripts in the
	// object directory is in the PATH.
	path := os.Getenv("PATH")
	path += ":" + pwd
	os.Setenv("PATH", path)
	os.Setenv("DABBA_DIR", pwd)
	c.Env = append(os.Environ(), "FILE="+filename)
	c.Dir = pwd
	// ignore error as it mostly shows error code when something fails.
	// I want to have what is written on the stderr
	out, err := c.CombinedOutput()
	if err != nil {
		log.Println(err)
	}
	return string(out)
}
