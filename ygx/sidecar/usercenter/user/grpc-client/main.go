package main

import (
	"context"

	_ "strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/core/lager"
	_ "github.com/go-chassis/go-chassis/server/grpc"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	_ "github.com/go-chassis/go-chassis/client/grpc"
	pb "github.com/go-chassis/ygx/libs/pb/xfcuser"
)

var secret []byte
var expireTime int

type JWTUserClaims struct { // token里面添加用户信息，验证token后可能会用到用户信息
	jwt.StandardClaims
	UserId   int    `json:"user_id"`
	PassWord string `json:"password"`
	Username string `json:"user_name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	RoleId   int32  `json:"role_id"`
}

func init() {
	secret = []byte("adb998230")
	expireTime = 3600
}

func main() {
	//Init framework
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed." + err.Error())
		return
	}
	//declare reply struct

	reply := &pb.UserReplyToken{}
	//Invoke with microservice name, schema ID and operation ID
	if err := core.NewRPCInvoker().Invoke(context.Background(), "ygx_usercenter_user_grpc-server", pb.XfcUser_serviceDesc.ServiceName, "UserLogin",
		&pb.UserRequest{UserName: "米奇", Passwd: "519475228fe35ad067744465c42a19b2"}, reply, core.WithProtocol("grpc")); err != nil {
		lager.Logger.Error("error" + err.Error())
	}
	lager.Logger.Info(reply.Token)
}
