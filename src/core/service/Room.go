package service

import (
	"core/constant"
	"core/entity"
	"core/util"
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

	//断开本身
	userInfo.Conn.Close()
	delete(userMap, userInfo.UserID)
	fmt.Printf("用户离线 %s .\n", userInfo.UserID) // 断开连接
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
		fmt.Printf("CUT %s TO %s\n", userInfo.UserID, userInfo.TargetUserID)
		command := entity.Command{Action: constant.ACTION_DOWN_CUT, Content: userInfo.TargetUserID}
		Send(userInfo.Conn, command)

		userInfo.Status = true
		userInfo.TargetUserID = ""
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
		fmt.Printf("DoAction不处理该指令:%s\n", command.Action)

	default:
		fmt.Printf("DoAction不支持该指令:%s\n", command.Action)
		SendError(userInfo.Conn, "不被支持的指令\n")
	}

}

func SendError(conn net.Conn, msg string) {
	command := entity.Command{Action: constant.ACTION_DOWN_ERROR, Content: msg}
	Send(conn, command)
}

func doOff(userInfo *entity.UserInfo, targetUserID string, userMap map[string]*entity.UserInfo) {

	fmt.Printf("OFF %s TO %s\n", userInfo.UserID, targetUserID)

	targetUserInfo, ok := userMap[targetUserID]
	if ok {
		doCut(targetUserInfo, userMap)
	}

	userInfo.Status = true
	userInfo.TargetUserID = ""

	command := entity.Command{Action: constant.ACTION_UP_OFF, Content: "true"}
	Send(userInfo.Conn, command)
}

func doSubscribeCallback(userInfo *entity.UserInfo, res string, userMap map[string]*entity.UserInfo) {

	fmt.Printf("SUBSCRIBE CALLBACK: %s TO %s Sey %s.\n", userInfo.UserID, userInfo.TargetUserID, res)

	targetUserInfo, ok := userMap[userInfo.TargetUserID]
	if ok {
		if res == "true" {
			//CALL 成功
			command := entity.Command{Action: constant.ACTION_UP_CALL, Content: "true"}
			doBilateralSubscribe(userInfo, targetUserInfo)
			fmt.Printf("CALL SUCCEED: %s TO %s.\n", targetUserInfo.UserID, userInfo.UserID)
			Send(targetUserInfo.Conn, command)
			return
		}

		command := entity.Command{Action: constant.ACTION_UP_CALL, Content: "false"}
		Send(targetUserInfo.Conn, command)
		fmt.Printf("CALL RETURN: %s  %s.\n", targetUserInfo.UserID, command.Content)
	}
}

func doCall(userInfo *entity.UserInfo, targetUserID string, userMap map[string]*entity.UserInfo) {

	targetUserInfo, ok := userMap[targetUserID]
	if ok { //用户存在
		if targetUserInfo.Status { //空闲中
			targetUserInfo.TargetUserID = userInfo.UserID
			doSubscribe(targetUserInfo, userInfo.UserID)
			return
		}
	}

	//忙线中 或拒绝
	command := entity.Command{Action: constant.ACTION_UP_CALL, Content: "false"}
	Send(userInfo.Conn, command)
	fmt.Printf("CALL FAIL: %s TO %s.\n", userInfo.UserID, targetUserInfo.UserID)

}

func doSubscribe(userInfo *entity.UserInfo, targetUserID string) {
	command := entity.Command{Action: constant.ACTION_DOWN_SUBSCRIBE, Content: targetUserID}
	Send(userInfo.Conn, command)
	fmt.Printf("SUBSCRIBE : %s TO %s.\n", userInfo.UserID, targetUserID)
}

func doBilateralSubscribe(userInfo *entity.UserInfo, targetUserInfo *entity.UserInfo) {
	userInfo.Status = false
	userInfo.TargetUserID = targetUserInfo.UserID

	targetUserInfo.Status = false
	targetUserInfo.TargetUserID = userInfo.UserID
}
