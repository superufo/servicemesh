package gateway

import (
	"context"
	"encoding/json"
	"errors"
	_ "strings"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	_ "github.com/go-chassis/go-chassis/client/grpc"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/core/lager"
	_ "github.com/go-chassis/go-chassis/server/grpc"
	"github.com/go-chassis/ygx/schemas/gateway/gclient"

	pbXfcUser "github.com/go-chassis/ygx/libs/pb/xfcuser"
)

// ChoiceGrpcClient1  ChoiceGrpcClient两种不同方法都可以
func ChoiceGrpcClient(serviceName string, method string, token string, version string, parameter []byte) (response interface{}, err error) {
	var (
		c        gclient.ClientInit
		request  interface{}
		reply    interface{}
		schemaId string
	)

	switch serviceName {
	case "ygx_usercenter_user_grpc-server":
		c = new(gclient.YgxUserCenterUserInit)
		request, reply, schemaId = gclient.NewClientInit(c, method)
		break
	case "":
	default:
		return nil, errors.New("没有找到对应的服务")
		break
	}

	jerr := json.Unmarshal(parameter, request)
	if jerr != nil {
		return nil, jerr
	}

	//Invoke with microservice name, schema ID and operation ID
	if err := core.NewRPCInvoker().Invoke(context.Background(), serviceName, schemaId, method,
		request, reply, core.WithProtocol("grpc")); err != nil {
		lager.Logger.Error("error" + err.Error())
		return nil, err
	} else {
		response = reply
	}

	return response, nil
}

// 下面的调用也是ok的
type CallGrpcClient func() (response interface{}, err error)

func execCall(serviceName string, method string, token string, version string, parameter []byte) (CallGrpcClient, error) {
	var (
		request  interface{}
		reply    interface{}
		schemaId string
	)

	if serviceName == "ygx_usercenter_user_grpc-server" {
		switch method {
		case "UserLogin":
			request = &pbXfcUser.UserRequest{}
			reply = &pbXfcUser.UserReplyToken{}
			schemaId = pbXfcUser.XfcUser_serviceDesc.ServiceName
			break
		case "":
			break
		default:
			break
		}

	}

	if request == nil || reply == nil || schemaId == "" {
		return nil, errors.New("初始化参数失败")
	}

	jerr := json.Unmarshal(parameter, request)
	if jerr != nil {
		return nil, jerr
	}

	return func() (response interface{}, err error) {
		//Invoke with microservice name, schema ID and operation ID
		if err := core.NewRPCInvoker().Invoke(context.Background(), serviceName, schemaId, method,
			request, reply, core.WithProtocol("grpc")); err != nil {
			lager.Logger.Error("error" + err.Error())
			return nil, err
		} else {
			response = reply
		}

		return response, nil
	}, nil
}

// ChoiceGrpcClient1  ChoiceGrpcClient两种不同方法都可以
func ChoiceGrpcClientBak1(serviceName string, method string, token string, version string, parameter []byte) (response interface{}, err error) {
	f, error := execCall(serviceName, method, token, version, parameter)
	if error != nil {
		return nil, error
	}
	return f()
}
