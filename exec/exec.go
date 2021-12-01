package exec

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Run command under the working directory.
// Sets FILENAME env variable to filename and adds pwd to PATH
func Run(pwd string, filename string, cmd string) string {
	c := exec.Command("sh", "-c", cmd)
	return runCommand(c, pwd, filename)
}

// RunStdin runs command under the working directory passing in stdin
// Sets FILENAME env variable to filename and adds pwd to PATH
func RunStdin(pwd string, filename string, stdin []byte) string {
	c := exec.Command("sh")
	c.Stdin = bytes.NewReader(stdin)
	return runCommand(c, pwd, filename)
}

func runCommand(c *exec.Cmd, pwd string, filename string) string {
	// todo: this is a simple hack to ensure the scripts in the
	// object directory is in the PATH.
	os.Setenv("PATH", addPath(os.Getenv("PATH"), pwd))
	os.Setenv("EDI_DIR", pwd)
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

func addPath(path, dir string) string {
	dirs := strings.Split(path, ":")
	for _, other := range dirs {
		if other == dir {
			return path
		}
	}
	dirs = append(dirs, dir)
	return strings.Join(dirs, ":")
}
