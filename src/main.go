package main

import (
	"core/conf"
	"core/constant"
	"core/entity"
	"core/service"
	"core/util"
	"fmt"
	"net"
	"os"
)

var userMap map[string]*entity.UserInfo

func main() {

	fmt.Println("Initiating server...")

	//导入配置文件
	configMap := conf.InitConfig("./config.lua")

	userMap = make(map[string]*entity.UserInfo)

	host := configMap["host"]
	port := configMap["port"]

	remote := host + ":" + port
	data := make([]byte, 1024)

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

		go func(con net.Conn) {
			fmt.Println("New connection: ", con.RemoteAddr())

			// Get client's name
			length, err := con.Read(data)
			if err != nil {
				fmt.Printf("断开连接 %v .\n", con.RemoteAddr())
				con.Close()
				return
			}

			receive := string(data[:length])

			command, err := util.GetCommon(receive)
			if err != nil || command.Action != constant.ACTION_UP_LOGIN { //第一道指令必须是LOGIN
				service.SendError(con, "第一道指令必须是LOGIN")
				fmt.Printf("断开连接:指令不对 %v .\n", con.RemoteAddr())
				con.Close()
				return
			}

			name := command.Content

			if userInfo, ok := userMap[name]; ok { //存在
				fmt.Printf("重新登陆: %s .\n", name)
				service.Disconnect(userInfo, userMap)
			}

			userInfo := entity.UserInfo{UserID: name, Status: true, Conn: con}

			//新用户接入
			userMap[name] = &userInfo

			//回传通知
			service.Send(con, entity.Command{Content: "true", Action: command.Action})

			// 开始从客户端接收消息
			for {
				length, err := con.Read(data)
				if err != nil {
					service.Disconnect(&userInfo, userMap)
					return
				}

				res = string(data[:length])

				service.DoAction(&userInfo, res, userMap)
			}
		}(conn)
	}
}
