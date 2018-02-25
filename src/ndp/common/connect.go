package common

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"bufio"
	"strings"
	"ndp/common/model"
)

func Connect(server model.ServerInfo) (*ssh.Client, error) {
	if server.Password != "" {
		return ConnPassword(server.Username, server.Password, server.Host, server.Port)
	}

	if server.PublicKey == "" {
		server.PublicKey = filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa");
	}
	return ConnPublicKey(server.Username, server.PublicKey, server.Host, server.Port)
}

func ConnPassword(user, password, host string, port int) (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}
	// connect to ssh
	addr := fmt.Sprintf("%s:%d", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	return sshClient, nil
}

func ConnPublicKey(user, keyPath, host string, port int) (*ssh.Client, error) {
	//先检查hostKey
	var hostKey = getHostKey(host)
	//读取秘钥
	privateKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		//HostKeyCallback: ssh.InsecureIgnoreHostKey(), //不检查 host key。 不推荐。
		HostKeyCallback: ssh.FixedHostKey(hostKey), //strongly recommend.
	}
	// Connect to the remote server and perform the SSH handshake.
	addr := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	return client, err
}

func getHostKey(host string) (ssh.PublicKey) {
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}
	if hostKey == nil {
		log.Fatalf("no hostkey for %s ,Please login at the terminal first。", host)
	}
	return hostKey;
}
