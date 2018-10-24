package constant

/*
常量
*/
const (
	Splitter            = "#"         // 指令分隔符
	ActionUpLogin       = "LOGIN"     // 上行指令：登陆
	ActionUpCall        = "CALL"      // 上行指令：呼叫
	ActionUpOff         = "OFF"       // 上行指令：挂断
	ActionDownSubscribe = "SUBSCRIBE" // 下行指令：订阅
	ActionDownCut       = "CUT"       // 下行指令：挂断
	ActionDownError     = "ERROR"     // 下行指令：报错
)
