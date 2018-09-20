package main

import (
	"fmt"
	"net"
	"os"
)

var clients []net.Conn

func main() {
	var (
		host   = "localhost"
		port   = "8000"
		remote = host + ":" + port
		data   = make([]byte, 1024)
	)
	fmt.Println("Initiating server...")

	lis, err := net.Listen("tcp", remote)
	defer lis.Close()

	if err != nil {
		fmt.Printf("Error when listen: %s, Err: %s\n", remote, err)
		os.Exit(-1)
	}

	for {
		var res string
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println("Error accepting client: ", err.Error())
			os.Exit(0)
		}
		clients = append(clients, conn)

		go func(con net.Conn) {
			fmt.Println("New connection: ", con.RemoteAddr())

			// Get client's name
			length, err := con.Read(data)
			if err != nil {
				fmt.Printf("Client %v quit.\n", con.RemoteAddr())
				con.Close()
				disconnect(con, con.RemoteAddr().String())
				return
			}
			name := string(data[:length])
			comeStr := name + " entered the room."
			notify(con, comeStr)

			// Begin recieve message from client
			for {
				length, err := con.Read(data)
				if err != nil {
					fmt.Printf("Client %s quit.\n", name)
					con.Close()
					disconnect(con, name)
					return
				}
				res = string(data[:length])
				sprdMsg := name + " said: " + res
				fmt.Println(sprdMsg)
				res = "You said:" + res
				con.Write([]byte(res))
				notify(con, sprdMsg)
			}
		}(conn)
	}
}

//广播
func notify(conn net.Conn, msg string) {
	for _, con := range clients {
		if con.RemoteAddr() != conn.RemoteAddr() {
			con.Write([]byte(msg))
		}
	}
}

func sendToClient(conn net.Conn, msg string) {

}

func disconnect(conn net.Conn, name string) {
	for index, con := range clients {
		if con.RemoteAddr() == conn.RemoteAddr() {
			disMsg := name + " has left the room."
			fmt.Println(disMsg)
			clients = append(clients[:index], clients[index+1:]...)
			notify(conn, disMsg)
		}
	}
}

func CheckErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
