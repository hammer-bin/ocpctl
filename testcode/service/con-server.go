package service

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"ocpctl/testcode/common"
	"os"
)

// SSH 연결 설정 함수
func SSHConfig(keyPath string, user string) (*ssh.ClientConfig, error) {

	// 키 파일 읽기
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	// SSH 키 파싱
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	// SSH 연결 설정
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return config, nil
}

func ConnectServer() ssh.Client {
	var conn *ssh.Client

	host := common.ConfInfo["master.ip"]
	user := common.ConfInfo["username"]
	pwd := common.ConfInfo["password"]

	fmt.Println(host, user, pwd)
	host = fmt.Sprint(host, ":22")

	//pKey := []byte("os-new-cp-common-key.pem")

	// 키 파일 열기
	key, err := os.ReadFile("os-new-cp-common-key.pem")
	if err != nil {
		log.Fatal("unable to read private key:", err)
	}

	// SSH 키 파싱
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("unable to parse private key:", err)
	}

	//signer, err = ssh.ParsePrivateKey(pKey)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}

	//var hostkeyCallback ssh.HostKeyCallback
	//hostkeyCallback, err = knownhosts.New("~/.ssh/known_hosts")
	//if err != nil {
	//	fmt.Println(err.Error())
	//}

	conf := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(pwd),
			ssh.PublicKeys(signer),
		},
	}

	conn, err = ssh.Dial("tcp", host, conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close()

	return *conn
}
