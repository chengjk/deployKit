package cmder

import (
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"strings"
)

func ExecRemote(sshClient *ssh.Client,  cmds []string) {
	// create session
	var session *ssh.Session
	var err error
	if session, err = sshClient.NewSession(); err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	cmdStr:=strings.Join(cmds,";")
	log.Println("execute cmd: " + cmdStr)
	session.Run(cmdStr)
}
