package main

import (
	"flag"
	"fmt"
	"github.com/goccy/go-yaml"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
)

type login struct {
	Username string `yaml:"username"`
	// Having password here is temp and needs to be probably encrypted coming in that gets decrypted throughout
	Password string `yaml:"password"`
	Msg      string `yaml:"msg"`
}

func check(errormsg error) {
	if errormsg != nil {
		log.Fatal(errormsg)
	}
}

func main() {
	// Grab flags
	var fileArg string
	var hostArg string
	var portArg string
	// Temp default value of the nokia file and 192.168.0.1
	flag.StringVar(&fileArg, "file", "nokia.yml", "File that contains login and a hello msg")
	flag.StringVar(&hostArg, "host", "192.168.0.1", "Host you want to target")
	flag.StringVar(&portArg, "port", "22", "Port on your host you want to target")
	// Parse flags
	flag.Parse()

	// Open and check file contents
	fileContentsPtr, err := os.ReadFile(fileArg)
	check(err)

	// Read contents into YAML parser -
	var yamlContents login
	if err := yaml.Unmarshal([]byte(fileContentsPtr), &yamlContents); err != nil {
		log.Fatal(err)
	}
	// Print out vars
	//fmt.Println("Username to use is: ", yamlContents.Username)
	//fmt.Println("Password to use is: ", yamlContents.Password)
	//fmt.Println("Connecting to", yamlContents.Username+"@"+hostArg, "on port", portArg)
	fmt.Println("Msg to use is: '", yamlContents.Msg, "'")

	// Create a session config then create a session
	sessionConfigPtr := &ssh.ClientConfig{
		User: yamlContents.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(yamlContents.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	fmt.Println("Trying to connect to", yamlContents.Username+"@"+hostArg, "on port", portArg+"...")
	conn, err := ssh.Dial("tcp", hostArg+":"+portArg, sessionConfigPtr)
	check(err)
	defer func() {
		conn.Close()
		fmt.Println("Closing connection")
	}()
	fmt.Println("Connected to", yamlContents.Username+"@"+hostArg, "on port", portArg+"...")

	// Create a new tty session
	clientSession, err := conn.NewSession()
	check(err)
	defer func() {
		clientSession.Close()
		fmt.Println("Closing session")
	}()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 115200,
		ssh.TTY_OP_OSPEED: 115200,
	}

	if err := clientSession.RequestPty("xterm", 40, 80, modes); err != nil {
		log.Fatal(err)
	}

	// Establish stdin and stdout, write msg from file into stdin, run command, and take output from stdout print to screen
	stdin, err := clientSession.StdinPipe()
	check(err)
	stdout, err := clientSession.StdoutPipe()
	check(err)
	if err := clientSession.Shell(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Sending msg...")
	// - Turn the message into bytes for the Write method, also add at end of message actually leaving the host (should probably be defined in a file but this is good for now)
	msg := []byte(yamlContents.Msg + "\n" + "quit\n")
	stdin.Write(msg)
	defer func() {
		fmt.Println("Command run, output below:")
		io.Copy(log.Writer(), stdout)
	}()
}
