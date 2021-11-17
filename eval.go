package main

import (
	"os"
	"os/exec"
)

func run(pwd string, filename string, cmd string) string {
	c := exec.Command("bash", "-c", cmd)
	// todo: this is a simple hack to ensure the scripts in the
	// object directory is in the PATH.
	path := os.Getenv("PATH")
	path += ":" + pwd
	os.Setenv("PATH", path)
	c.Env = append(os.Environ(), "FILE="+filename)
	c.Dir = pwd
	// ignore error as it mostly shows error code when something fails.
	// I want to have what is written on the stderr
	out, _ := c.CombinedOutput()
	return string(out)
}
