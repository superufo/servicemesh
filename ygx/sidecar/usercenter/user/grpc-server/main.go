package main

import (
	"errors"
	"fmt"
	"log"
	_ "strings"

	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-chassis/ygx/libs/common"
	pb "github.com/go-chassis/ygx/libs/pb/xfcuser"

	"github.com/go-mesh/openlogging"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	_ "github.com/go-chassis/go-chassis/server/grpc"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/go-chassis/ygx/libs/model/usercenter"

	"github.com/dgrijalva/jwt-go"

	"github.com/go-chassis/foundation/token"
)

var secret []byte
var expireTime int

type JWTUserClaims struct { // token里面添加用户信息，验证token后可能会用到用户信息
	jwt.StandardClaims
	UserId   int    `json:"user_id"`
	Passwd   string `json:"password"`
	Username string `json:"user_name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	RoleId   int32  `json:"role_id"`
}

func init() {
	secret = []byte("adb998230")
	expireTime = 3600
}

// Server if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/grpc/server/
// Server is used to implement helloworld.GreeterServer.
type UserCenterServer struct{}

// SayHello implements helloworld.GreeterServer
func (s *UserCenterServer) GetXfcUserInfo(ctx context.Context, in *pb.UserSelectRequest) (*pb.UserReply, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println(md["x-user"])

	userData := usercenter.XfcUser{}

	user_id := in.GetUserId()
	email := in.GetEmail()
	mobile := in.GetMobile()
	user_name := in.GetUserName()

	if int(user_id) == 0 && email == "" && mobile == "" && user_name == "" {
		openlogging.Error("传入的请求为空.")
		return nil, errors.New("传入的请求为空")
	}

	// 通过UserSelectRequest 查询sql获得结果
	common.DB.Where("user_id = ? or user_name= ？ or mobile=?  or email=? ", user_id, email, mobile, user_name).First(&userData)

	return &pb.UserReply{
		UserId:               userData.UserId,
		UserName:             userData.UserName,
		Passwd:               userData.Passwd,
		Nick:                 userData.Nick,
		Mobile:               userData.Mobile,
		Email:                userData.Email,
		Post:                 userData.Post,
		TeamId:               userData.TeamId,
		TeamName:             userData.TeamName,
		Introduce:            userData.Introduce,
		Department:           userData.Department,
		DepartmentId:         userData.DepartmentId,
		Balance:              userData.Balance,
		FreezeBalance:        userData.FreezeBalance,
		RealBalance:          userData.RealBalance,
		RoleId:               userData.RoleId,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}, nil
}

func (s *UserCenterServer) UserLogin(ctx context.Context, in *pb.AuthRequest) (*pb.UserReplyToken, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println(md["x-user"])

	userData := usercenter.XfcUser{}

	user_name := in.GetUserName()
	passwd := in.GetPasswd()

	if user_name == "" || passwd == "" {
		openlogging.Error("传入的用户名或密码为空.")
		return nil, errors.New("传入的用户名或密码为空")
	}

	log.Println(user_name)
	log.Println(passwd)
	t := token.XfcUserToken{}
	tokenstr, err := t.XfcUserTokenGenenral(user_name, passwd)
	if err != nil {
		openlogging.Error("生成token失败." + fmt.Sprint("%s", err))
		return nil, errors.New("生成token失败." + fmt.Sprint("%s", err))
	}

	userData.UserName = user_name
	userData.Passwd = passwd
	// 通过UserSelectRequest 查询sql获得结果
	common.DB.Where("user_name= ? and passwd=?  ", user_name, passwd).First(&userData)
	if userData.UserId == 0 {
		openlogging.Error("不存在的用户.")
		return nil, errors.New("不存在的用户")
	}

	userReply := pb.UserReply{
		UserId:               userData.UserId,
		UserName:             userData.UserName,
		Passwd:               userData.Passwd,
		Nick:                 userData.Nick,
		Mobile:               userData.Mobile,
		Email:                userData.Email,
		Post:                 userData.Post,
		TeamId:               userData.TeamId,
		TeamName:             userData.TeamName,
		Introduce:            userData.Introduce,
		Department:           userData.Department,
		DepartmentId:         userData.DepartmentId,
		Balance:              userData.Balance,
		FreezeBalance:        userData.FreezeBalance,
		RealBalance:          userData.RealBalance,
		RoleId:               userData.RoleId,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}

	return &pb.UserReplyToken{
		UserReply:            &userReply,
		Token:                tokenstr,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}, nil
}

func (s *UserCenterServer) UserVerify(ctx context.Context, in *pb.TokenInfo) (*pb.UserToken, error) {
	return nil, nil
}

func (s *UserCenterServer) GetXfcUserBySelect(ctx context.Context, in *pb.SearchRequest) (*pb.MultiUserReply, error) {
	return nil, nil
}

func (s *UserCenterServer) UpdateXfcUserInfo(ctx context.Context, in *pb.UserRequest) (*pb.UserReply, error) {
	return nil, nil
}

func (s *UserCenterServer) DeleteXfcUserInfo(ctx context.Context, in *pb.UserDeleteRequest) (*pb.UserReply, error) {
	return nil, nil
}

func main() {
	chassis.RegisterSchema("grpc", &UserCenterServer{}, server.WithRPCServiceDesc(&pb.XfcUser_serviceDesc))
	//chassis.RegisterSchema("grpc", &UserCenterServer{}, server.WithSchemaID("UserService"))
	if err := chassis.Init(); err != nil {
		openlogging.Error("Init failed.")
		return
	}
	chassis.Run()
}
