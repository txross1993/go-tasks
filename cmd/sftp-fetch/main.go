package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type sftpConfig struct {
	Username string
	Password string
	Host     string
	Port     string
}

func (s sftpConfig) ConnectionString() string {
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	return addr
}

func ConnectSFTP(cfg sftpConfig) (*sftp.Client, error) {
	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(cfg.Password))
	sshCfg := &ssh.ClientConfig{
		User:    cfg.Username,
		Auth:    auth,
		Timeout: time.Duration(600 * time.Second),
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	sshConn, err := ssh.Dial("tcp", cfg.ConnectionString(), sshCfg)
	if err != nil {
		return nil, err
	}
	return sftp.NewClient(sshConn)
}

func main() {
	cfg := sftpConfig{}
	var file string

	flag.StringVar(&cfg.Username, "user", "", "User to auth with sftp")
	flag.StringVar(&cfg.Password, "pass", "", "Password to auth with sftp")
	flag.StringVar(&cfg.Host, "host", "", "Sftp Host")
	flag.StringVar(&cfg.Port, "port", "", "Sftp Port")
	flag.StringVar(&file, "file", "", "The fully qualified path of the file to fetch")
	flag.Parse()

	sftpClient, err := ConnectSFTP(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer sftpClient.Close()

	start := time.Now()
	f, err := sftpClient.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(f.Name(), b, os.ModePerm)
	elapsed := time.Since(start)
	fmt.Println("reading file %s took %d", f.Name(), elapsed)

}
