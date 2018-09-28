package service

import (
	"core/constant"
	"core/entity"
	"core/util"
	"core/util/logUtil"
	"fmt"
	"net"
)

//Disconnect 断开连接
//userInfo 当前用户
//userMap userMap
func Disconnect(userInfo *entity.UserInfo, userMap map[string]*entity.UserInfo) {
	//检查是否正在通话中
	if !userInfo.Status {
		//通知对方断开
		doCutByuserID(userInfo.TargetUserID, userMap)
	}

	logUtil.LOG_INFO("User Offline:%s.", userInfo.UserID)
	// 断开连接
	userInfo.Conn.Close()
	delete(userMap, userInfo.UserID)
}

//doCutByuserID 通话被挂断
//userID 目标用户
//userMap userMap
func doCutByuserID(userID string, userMap map[string]*entity.UserInfo) {
	value, ok := userMap[userID]
	if ok == true {
		doCut(value, userMap)
	}
}

func doCut(userInfo *entity.UserInfo, userMap map[string]*entity.UserInfo) {
	if !userInfo.Status {
		logUtil.LOG_INFO("CUT %s TO %s\n", userInfo.UserID, userInfo.TargetUserID)
		command := entity.Command{Action: constant.ACTION_DOWN_CUT, Content: userInfo.TargetUserID}
		Send(userInfo.Conn, command)

		userInfo.DoOff()
	}
}

//Send 下发指令
//conn
//command
func Send(conn net.Conn, command entity.Command) {
	conn.Write(command.ToByteArry())
}

//DoAction 指令处理
//userInfo 当前用户
//res 指令
//userMap userMap
func DoAction(userInfo *entity.UserInfo, res string, userMap map[string]*entity.UserInfo) {
	//解释指令
	command, err := util.GetCommon(res)
	if err != nil {
		SendError(userInfo.Conn, fmt.Sprintf("DoAction无法理解该指令：%s", res))
	}

	//遍历指令
	switch command.Action {
	case constant.ACTION_UP_CALL:
		doCall(userInfo, command.Content, userMap)

	case constant.ACTION_UP_OFF:
		doOff(userInfo, command.Content, userMap)

	case constant.ACTION_DOWN_SUBSCRIBE:
		doSubscribeCallback(userInfo, command.Content, userMap)

	case constant.ACTION_DOWN_CUT:
		logUtil.LOG_INFO("DoAction不处理该指令:%s\n", command.Action)

	default:
		logUtil.LOG_INFO("DoAction不支持该指令:%s\n", command.Action)
		SendError(userInfo.Conn, "不被支持的指令\n")
	}

}

//SendError 下发错误信息
//conn
//msg
func SendError(conn net.Conn, msg string) {
	command := entity.Command{Action: constant.ACTION_DOWN_ERROR, Content: msg}
	Send(conn, command)
}

func doOff(userInfo *entity.UserInfo, targetUserID string, userMap map[string]*entity.UserInfo) {

	logUtil.LOG_INFO("OFF %s TO %s\n", userInfo.UserID, targetUserID)

	targetUserInfo, ok := userMap[targetUserID]
	if ok {
		doCut(targetUserInfo, userMap)
	}

	userInfo.DoOff()

	command := entity.Command{Action: constant.ACTION_UP_OFF, Content: "true"}
	Send(userInfo.Conn, command)
}

func doSubscribeCallback(userInfo *entity.UserInfo, res string, userMap map[string]*entity.UserInfo) {

	logUtil.LOG_INFO("SUBSCRIBE CALLBACK: %s TO %s Sey %s.\n", userInfo.UserID, userInfo.TargetUserID, res)

	targetUserInfo, ok := userMap[userInfo.TargetUserID]
	if ok {
		if res == "true" {
			//CALL 成功
			command := entity.Command{Action: constant.ACTION_UP_CALL, Content: "true"}

			userInfo.DoCall(targetUserInfo.UserID)
			targetUserInfo.DoCall(userInfo.UserID)

			logUtil.LOG_INFO("CALL SUCCEED: %s TO %s.\n", targetUserInfo.UserID, userInfo.UserID)
			Send(targetUserInfo.Conn, command)
			return
		}

		command := entity.Command{Action: constant.ACTION_UP_CALL, Content: "false"}
		Send(targetUserInfo.Conn, command)
		logUtil.LOG_INFO("CALL RETURN: %s  %s.\n", targetUserInfo.UserID, command.Content)
	}
}

func doCall(userInfo *entity.UserInfo, targetUserID string, userMap map[string]*entity.UserInfo) {

	targetUserInfo, ok := userMap[targetUserID]
	if ok { //用户存在
		if targetUserInfo.Status { //空闲中
			targetUserInfo.SetTargetUserID(userInfo.UserID)

			doSubscribe(targetUserInfo, userInfo.UserID)
			return
		}

		logUtil.LOG_INFO("CALL FAIL: %s Busy.\n", targetUserID)

	} else {
		logUtil.LOG_INFO("CALL FAIL: %s Offline.\n", targetUserID)
	}

	//忙线中 或拒绝
	command := entity.Command{Action: constant.ACTION_UP_CALL, Content: "false"}
	Send(userInfo.Conn, command)
	logUtil.LOG_INFO("CALL FAIL: %s TO %s.\n", userInfo.UserID, targetUserID)

}

func doSubscribe(userInfo *entity.UserInfo, targetUserID string) {
	command := entity.Command{Action: constant.ACTION_DOWN_SUBSCRIBE, Content: targetUserID}
	Send(userInfo.Conn, command)
	logUtil.LOG_INFO("SUBSCRIBE : %s TO %s.\n", userInfo.UserID, targetUserID)
}
