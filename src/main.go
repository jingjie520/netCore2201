package main

import (
	"core/conf"
	"core/constant"
	"core/entity"
	"core/service"
	"core/util"
	"core/util/logUtil"
	"net"
	"os"
)

var userMap map[string]*entity.UserInfo

func main() {
	logUtil.LOG_INFO("Service Initiating.")

	//导入配置文件
	configMap := conf.InitConfig("./config.lua")

	userMap = make(map[string]*entity.UserInfo)

	host := configMap["host"]
	port := configMap["port"]

	remote := host + ":" + port

	data := make([]byte, 1024)

	lis, err := net.Listen("tcp", remote)
	defer lis.Close()

	logUtil.LOG_INFO("Service Run AT %s.", remote)

	if err != nil {
		logUtil.LOG_ERROR("Error when listen: %s, Err: %s\n", remote, err)
		os.Exit(-1)
	}

	for {
		var res string
		conn, err := lis.Accept()
		if err != nil {
			logUtil.LOG_ERROR("Error accepting client: ", err.Error())
			os.Exit(0)
		}

		go func(con net.Conn) {
			logUtil.LOG_INFO("New connection: %v.", con.RemoteAddr())

			// Get client's name
			length, err := con.Read(data)
			if err != nil {
				logUtil.LOG_INFO("Offline : %v.", con.RemoteAddr())
				con.Close()
				return
			}

			receive := string(data[:length])

			command, err := util.GetCommon(receive)
			if err != nil || command.Action != constant.ACTION_UP_LOGIN { //第一道指令必须是LOGIN
				logUtil.LOG_INFO("Close connection : %v.", con.RemoteAddr())
				service.SendError(con, "第一道指令必须是LOGIN")
				con.Close()
				return
			}

			name := command.Content

			if userInfo, ok := userMap[name]; ok { //存在
				logUtil.LOG_INFO("Rest Login : %s.", name)
				service.Disconnect(userInfo, userMap)
			}

			userInfo := entity.UserInfo{UserID: name, Status: true, Conn: con}

			//新用户接入
			userMap[name] = &userInfo

			//回传通知
			service.Send(con, entity.Command{Content: "true", Action: command.Action})

			logUtil.LOG_INFO("User Login done : %s.", name)

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
