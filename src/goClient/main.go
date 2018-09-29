package main

import (
	"bufio"
	"core/util/logUtil"
	"fmt"
	"net"
	"os"
)

var writeStr, readStr = make([]byte, 1024), make([]byte, 1024)

func main() {
	var (
		host   = "127.0.0.1"
		port   = "8000"
		remote = host + ":" + port
		reader = bufio.NewReader(os.Stdin)
	)

	con, err := net.Dial("tcp", remote)

	if err != nil {
		logUtil.LOG_INFO("无法连接服务器:%s.", remote)
		os.Exit(-1)
	}

	defer con.Close()

	logUtil.LOG_INFO("已连接服务器：%s", remote)

	fmt.Printf("请输入您的昵称: ")
	fmt.Scanf("%s", &writeStr)
	in, err := con.Write([]byte(writeStr))
	if err != nil {
		fmt.Printf("Error when send to server: %d\n", in)
		os.Exit(0)
	}

	fmt.Println("Now begin to talk!")
	go read(con)

	for {
		writeStr, _, _ = reader.ReadLine()
		if string(writeStr) == "quit" {
			fmt.Println("Communication terminated.")
			os.Exit(1)
		}

		in, err := con.Write([]byte(writeStr))
		if err != nil {
			fmt.Printf("Error when send to server: %d\n", in)
			os.Exit(0)
		}

	}
}

func read(conn net.Conn) {
	for {
		length, err := conn.Read(readStr)
		if err != nil {
			fmt.Printf("Error when read from server. Error:%s\n", err)
			os.Exit(0)
		}
		fmt.Println(string(readStr[:length]))
	}
}
