package i18n

func init() {
	Message = make(map[MessageType]string)

	//************ 错误提示信息************/
	Message[FAILURE] = "失败"

	Message[ERROR_PASSWD] = "密码错误"
	Message[EMPTY] = "空值"

	Message[NO_SERVER] = "没有找到对应的服务"
	Message[ONlY_REST] = "系统仅仅只支持Rest协议调用"

	//************ 正确提示信息************/
	Message[SUCCESS] = "成功"
}

type AppCnMessage struct {
}

func (cm *AppCnMessage) GetMessage(mesNum MessageType) string {
	message, err := Message[mesNum]
	if err != true {
		panic("没有定义消息类型")
	}

	return message
}
