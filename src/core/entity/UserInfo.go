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
