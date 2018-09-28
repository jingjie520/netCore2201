package entity

import (
	"net"
)

// UserInfo 用户信息
type UserInfo struct {
	UserID       string
	Status       bool
	TargetUserID string
	Conn         net.Conn
}

func (userInfo *UserInfo) SetTargetUserID(targetUserID string) {
	userInfo.TargetUserID = targetUserID
}

func (userInfo *UserInfo) DoCall(targetUserID string) {
	userInfo.Status = false
	userInfo.TargetUserID = targetUserID
}

func (userInfo *UserInfo) DoOff() {
	userInfo.Status = true
	userInfo.TargetUserID = ""
}
