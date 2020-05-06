package i18n

import "fmt"

type MessageType int

var Lang string
var Plat string
var Message map[MessageType]string

// 成功
const (
	//************ 错误提示信息************/
	FAILURE MessageType = 1

	//系统级别的错误
	NO_SERVER MessageType = -100
	ONlY_REST MessageType = -101

	//密码错误
	ERROR_PASSWD MessageType = -1001

	//空值
	EMPTY         MessageType = -1002
	PARAMETEREOOR MessageType = -1003

	//************ 正确提示信息************/
	SUCCESS MessageType = 0
)

func SetLang(lang string) {
	Lang = lang
}

func GetLang() string {
	return Lang
}

func SetPlat(plat string) {
	Plat = plat
}

func GetPlat() string {
	return Plat
}

func SetPlatLang(plat string, lang string) {
	Plat = plat
	Lang = lang
}

func GetPlatLang() (plat string, lang string) {
	return Plat, Lang
}

type LangMessage interface {
	GetMessage(num MessageType) string
}

func LoadMessage(mess MessageType) string {
	platLang := fmt.Sprintf("%s_%s")
	var lang LangMessage

	switch platLang {
	case "admin_cn":
		lang = new(AdminCnMessage)
	case "app_en":
		break
	case "app_cn":
	default:
		lang = new(AppCnMessage)
		break
	}

	return lang.GetMessage(mess)
}
