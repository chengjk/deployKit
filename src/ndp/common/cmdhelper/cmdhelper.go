package cmdhelper

import (
	"golang.org/x/crypto/ssh"
	"log"
	"os"
)

func ExecRemote(sshClient *ssh.Client, basePath string, cmds []string) {
	// create session
	var session *ssh.Session
	var err error
	if session, err = sshClient.NewSession(); err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	cmdStr := "cd " + basePath
	for _, cmd := range cmds {
		cmdStr = cmdStr + ";" + cmd
	}
	log.Println("execute cmd :" + cmdStr)
	session.Run(cmdStr)
	log.Println("execute suffix cmd succeed!")
}
