package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

type Command struct {
	Id      string // A unique string to track commands commands from server usually start with an s
	Command string
	Args    []string
	Server  bool
	Pwd     string

	parsed  []string // Stores parsed command and arguments.
	session *Session // the active session that executed the command
}

// Regex for splitting command and arguments.
var cmdRegex = regexp.MustCompile("'.+'|\".+\"|\\S+")

// Regex for valid UNIX file name
var validFile = regexp.MustCompile(`[.a-zA-Z-0-9_/\-]*[a-zA-Z-0-9_/\-]`)

func (c *Command) Exec(s *Session) {
	c.init(s)

	if len(c.parsed) <= 0 {
		return
	}

	switch c.parsed[0] {
	default:
		c.run()

	// Open buffer/file
	case "E":
		c.open()

	case "W":
		c.save()

	case "Cancel":
		c.cancel()

	case "Context":
		c.context()

	case "cd":
		c.cd()

	case "L":
		c.buffers()
	}
}

// init sets the parsed attribute of Command by splitting the Command
// attribute to command, args...
func (c *Command) init(s *Session) {
	c.parsed = cmdRegex.FindAllString(c.Command, -1)

	// Triming the quotes so that while executing a command
	// the quotes are not send as arguments.
	for i, _ := range c.parsed {
		c.parsed[i] = strings.Trim(c.parsed[i], "'\"")
	}
	c.session = s
}

func (c *Command) pushResponse(id string, cmd string, server bool, pwd string, args ...string) {
	response := Command{
		Id:      id,
		Command: cmd,
		Args:    args,
		Server:  server,
		Pwd:     pwd,
	}
	c.session.output <- &response
}

func (c *Command) push(id string, cmd string, args ...string) {
	c.pushResponse(id, cmd, false, c.session.cwd, args...)
}

func (c *Command) pushCwd(id string, cmd string, cwd string, args ...string) {
	c.pushResponse(id, cmd, false, cwd, args...)
}

// open sends a valid file for editing.
func (c *Command) open() {
	if len(c.parsed) <= 1 {
		return
	}
	c.openFile(c.parsed[1:])
	c.done()
}

// context handles right-clicks on the client.
func (c *Command) context() {
	for _, arg := range c.Args {
		matches := validFile.FindAllString(arg, -1)
		c.openFile(matches)
	}
}

func (c *Command) openFile(files []string) {
	pwd := c.session.cwd
	if c.Pwd != "" {
		pwd = c.Pwd
	}
	for _, filename := range files {
		if pwd != "" && !filepath.IsAbs(filename) {
			filename = filepath.Join(pwd, filename)
		}
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			log.Println(filename, "does not exist")
			continue
		}
		filename, err := filepath.Abs(filename)
		if err != nil {
			c.push(c.Id, "console", err.Error())
			log.Println(err)
			continue
		}

		filestat, err := os.Stat(filename)
		if err != nil {
			c.push(c.Id, "console", err.Error())
			log.Println(err)
			continue
		}
		if filestat.IsDir() {
			c.openDir(filename)
			continue
		}

		c.session.current = filename
		if _, ok := c.session.buffers[filename]; ok {
			c.push(c.Id, "openFile", filename)
			continue
		}

		file, err := ioutil.ReadFile(filename)
		if err != nil {
			c.push(c.Id, "console", err.Error())
			log.Println(err)
			continue
		}
		c.session.buffers[filename] = c.session.bufferCounter
		c.session.bufferCounter++

		possExt := strings.Split(filename, ".")
		mode := possExt[len(possExt)-1]
		c.push(c.Id, "newFile", filename, string(file), MODES[mode])
	}
}

// openDir initiates a server command to list the files in a given directory.
func (c *Command) openDir(dir string) {
	files, _ := ioutil.ReadDir(dir)
	response := ""
	for _, f := range files {
		response = response + f.Name() + "\n"
	}
	relDir := relativePath(c.session.cwd, dir)

	c.pushCwd(c.Id, "openDir", dir, relDir, response)
	c.done()
}

// save writes changes of the buffer to the disk.
func (c *Command) save() {
	if len(c.Args) == 2 {
		ioutil.WriteFile(c.Args[0], []byte(c.Args[1]), os.ModePerm)
	}
	c.done()
}

// run executes the appropriate os commands with arguments and sends response
// back to the client to be displayed in the console.
func (c *Command) run() {
	defer c.done()
	log.Println("Executing ", c.Command)
	var oscmd *exec.Cmd

	if len(c.parsed) > 1 {
		oscmd = exec.Command(c.parsed[0], c.parsed[1:]...)
	} else {
		oscmd = exec.Command(c.parsed[0])
	}
	if c.session.cwd != "" {
		oscmd.Dir = c.session.cwd
	}

	stdout, err := oscmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		c.push(c.Id, "console", err.Error())
		return
	}
	stderr, err := oscmd.StderrPipe()
	if err != nil {
		log.Println(err)
		c.push(c.Id, "console", err.Error())
		return
	}

	err = oscmd.Start()
	if err != nil {
		c.push(c.Id, "console", err.Error())
		log.Println(err)
		return
	}
	c.session.processes[c.Id] = oscmd.Process.Pid

	reader := bufio.NewReader(stdout)
	readerErr := bufio.NewReader(stderr)
	go c.readAndPush(readerErr)
	c.readAndPush(reader)

	oscmd.Wait()
}

func (c *Command) readAndPush(reader *bufio.Reader) {
	for {
		line, err := reader.ReadString('\n')
		c.push(c.Id, "console", line)
		if err != nil {
			break
		}
	}
}

func (c *Command) cancel() {
	log.Println("Cancelling ", c.Id)
	pid, ok := c.session.processes[c.Id]
	if !ok {
		log.Println("Can't find command, ", c.Id)
		c.push(c.Id, "console", "Command exited.")
		return
	}
	if pid == 0 {
		log.Println("Command seems to have exited.")
		c.push(c.Id, "console", "Command exited.")
		return
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		log.Println("Command seems to have exited.")
		c.push(c.Id, "console", err.Error())
		return
	}
	err = process.Kill()
	if err != nil {
		log.Println(err)
		c.push(c.Id, "console", err.Error())
		return
	}
	c.push(c.Id, "console", "Process Killed.")
}

// cd changes the current working directory of the session.
func (c *Command) cd() {
	dirname := ""
	if len(c.parsed) < 2 {
		usr, err := user.Current()
		if err != nil {
			c.push(c.Id, "console", "Could not find default home directory.")
			return
		}
		dirname = usr.HomeDir
	} else {
		dirname = c.parsed[1]
	}
	if c.session.cwd != "" {
		if !filepath.IsAbs(dirname) {
			dirname = filepath.Join(c.session.cwd, dirname)
		}
	}
	fileinfo, err := os.Stat(dirname)
	if err != nil {
		c.push(c.Id, "console", err.Error())
		log.Println(err)
		return
	}
	if !fileinfo.IsDir() {
		c.push(c.Id, "console", dirname+" is not a directory.")
		log.Println(err)
		return
	}
	c.session.cwd = dirname
	c.done()
}

func (c *Command) buffers() {
	for file, id := range c.session.buffers {
		format := "%d %s"
		if file == c.session.current {
			format = "*%d %s"
		}
		r := fmt.Sprintf(format, id, relativePath(c.session.cwd, file))
		c.push(c.Id, "console", r)
	}
	c.done()
}

// done sends a command complete to the client.
func (c *Command) done() {
	c.push(c.Id, "done")
}
