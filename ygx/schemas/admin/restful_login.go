package admin

import (
	"fmt"
	"net/http"

	"github.com/go-chassis/foundation/token"

	"github.com/go-chassis/ygx/libs/common"

	"github.com/go-mesh/openlogging"

	rf "github.com/go-chassis/go-chassis/server/restful"
)

type RestFulAdminLogin struct {
}

type Data struct {
	Token string `json:"token"`
}

// Root Sayhello is a method used to reply user with hello
func (r *RestFulAdminLogin) Root(b *rf.Context) {
	b.Write([]byte(fmt.Sprintf("x-forwarded-host %s", b.ReadRequest().Host)))
}

func (r *RestFulAdminLogin) Login(b *rf.Context) {
	reslut := struct {
		UserName string
		Password string
	}{}

	reponse := common.Response{
		Status:  0,
		Data:    Data{},
		Message: "",
	}

	data := Data{Token: ""}

	err := b.ReadEntity(&reslut)
	if err != nil {
		b.WriteHeaderAndJSON(http.StatusInternalServerError, reslut, "application/json")
		return
	}

	t := token.Token{}
	tokenstr, err := t.TokenGenenral(reslut.UserName, reslut.Password)

	if err != nil {
		reponse.Status = 1
		reponse.Message = fmt.Sprintf("生成Token失败,错误:%s", err)
		openlogging.Error(fmt.Sprintf("生成Token失败,错误:%s", err))
	} else {
		data.Token = tokenstr
		reponse.Data = data
		reponse.Status = 0
		openlogging.Info("生成Token成功")
	}

	b.WriteJSON(reponse, "application/json", "")
	return
}

func (r *RestFulAdminLogin) VerifyToken(b *rf.Context) {
	reponse := common.Response{
		Status:  0,
		Data:    Data{""},
		Message: "",
	}

	tokenstr := b.ReadHeader("token")
	t := token.Token{tokenstr}

	_, err := t.TokenVerify()
	if err != nil {
		reponse.Status = 1
		reponse.Message = "验证Token失败"
		openlogging.Error("验证Token失败")
	} else {
		reponse.Status = 0
		reponse.Message = "验证Token成功"
		openlogging.Info("验证Token成功")
	}

	b.WriteJSON(reponse, "application/json", "")
	return
}

// URLPatterns helps to respond for corresponding API calls
func (r *RestFulAdminLogin) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/", ResourceFunc: r.Root,
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodPost, Path: "/login",
			ResourceFunc: r.Login,
			Metadata: map[string]interface{}{
				"tags": []string{"admin", "test"},
			},
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodPost, Path: "/verifyToken",
			ResourceFunc: r.VerifyToken,
			Metadata: map[string]interface{}{
				"tags": []string{"users", "test"},
			},
			Returns: []*rf.Returns{{Code: 200}}},
	}
}
