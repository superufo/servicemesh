package cn

var Message map[int]string

// 成功
const (
	SUCCESS int16 = 0
	FAILURE int16 = -1

	//密码错误
	ERROR_PASSWD int16 = -1001

	//空值
	EMPTY int16 = -1002
)

func init() {
	Message = make(map[int]string)

	Message[SUCCESS] = "成功"
	Message[FAILURE] = "失败"

	Message[ERROR_PASSWD] = "密码错误"
	Message[EMPTY] = "空值"
}

func GetMessage(num int) string {
	message, error := Message[num]
	if error != null {
		panic("没有定义消息类型")
	}

	return message
}
