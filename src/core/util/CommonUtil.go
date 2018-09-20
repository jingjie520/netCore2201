package util

import (
	"core/constant"
	"core/entity"
	"errors"
	"strings"
)

//GetCommon 解析指令
// common #分割的命令行
// *Command 返回值
func GetCommon(common string) (entity.Command, error) {
	c := strings.Split(common, constant.CONST_SPLITTER)

	if len(c) == 2 {
		return entity.Command{Action: c[0], Content: c[1]}, nil
	}

	return entity.Command{}, errors.New("指令解释失败！")
}
