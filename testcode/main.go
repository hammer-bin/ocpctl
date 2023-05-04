package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"ocpctl/testcode/common"
	"ocpctl/testcode/service"
	"os"

	"golang.org/x/crypto/ssh"
)

func main() {

	//conn := service.ConnectServer()

	host := common.ConfInfo["master.ip"]
	user := common.ConfInfo["username"]
	keyPath := common.ConfInfo["key.path"]
	host = fmt.Sprint(host, ":22")

	// SSH 연결 설정
	config, err := service.SSHConfig(keyPath, user)
	if err != nil {
		log.Fatalf("unable to create ssh config: %v", err)
	}

	// 원격 서버에 연결
	conn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer conn.Close()

	// SSH 세션 생성
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("unable to create session: %v", err)
	}
	defer session.Close()

	var stdin io.WriteCloser
	var stdout, stderr io.Reader

	stdin, err = session.StdinPipe()
	if err != nil {
		fmt.Println(err.Error())
	}

	stdout, err = session.StdoutPipe()
	if err != nil {
		fmt.Println(err.Error())
	}

	stderr, err = session.StderrPipe()
	if err != nil {
		fmt.Println(err.Error())
	}

	wr := make(chan []byte, 10)

	go func() {
		for {
			select {
			case d := <-wr:
				_, err := stdin.Write(d)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for {
			if tkn := scanner.Scan(); tkn {
				rcv := scanner.Bytes()

				raw := make([]byte, len(rcv))
				copy(raw, rcv)

				fmt.Println(string(raw))
			} else {
				if scanner.Err() != nil {
					fmt.Println(scanner.Err())
				} else {
					fmt.Println("io.EOF")
				}
				return
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)

		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	session.Shell()

	for {
		fmt.Println("$")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		text := scanner.Text()

		wr <- []byte(text + "\n")
	}
}
