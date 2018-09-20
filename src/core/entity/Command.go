package entity

import "fmt"

// Command 指令
// Action 指令名称
// Content 指令内容
type Command struct {
	Action  string
	Content string
}

//ToString 指令转字符串
//command
//return string
func (command Command) ToString() string {
	return fmt.Sprintf("%s#%s", command.Action, command.Content)
}

//ToByteArry 指令转字节数组
//command
//return []byte
func (command Command) ToByteArry() []byte {
	return []byte(command.ToString())
}
