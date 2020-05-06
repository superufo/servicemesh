package gclient

import (
	"github.com/golang/protobuf/proto"

	pbXfcUser "github.com/go-chassis/ygx/libs/pb/xfcuser"
)

//proto "github.com/golang/protobuf/proto"
type GrpcRequest proto.Message
type GrpcReply proto.Message

type ClientInit interface {
	Inits(method string) (request GrpcRequest, reply GrpcReply, schemaId string)
}

func NewClientInit(c ClientInit, method string) (request GrpcRequest, reply GrpcReply, schemaId string) {
	return c.Inits(method)
}

//xfcuser
type YgxUserCenterUserInit struct {
}

func (yui *YgxUserCenterUserInit) Inits(method string) (request GrpcRequest, reply GrpcReply, schemaId string) {
	switch method {
	case "UserLogin":
		request = &pbXfcUser.UserRequest{}
		reply = &pbXfcUser.UserReplyToken{}
		schemaId = pbXfcUser.XfcUser_serviceDesc.ServiceName
		break
	case "":
	default:
		return nil, nil, ""
		break
	}

	return request, reply, schemaId
}
