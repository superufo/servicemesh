package gateway

// tls https://docs.go-chassis.com/user-guides/tls.html?highlight=tls
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/go-mesh/openlogging"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	_ "github.com/go-chassis/go-chassis/bootstrap"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"
	corecommon "github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/endpoint"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
	rf "github.com/go-chassis/go-chassis/server/restful"
	"github.com/go-chassis/ygx/libs/common"
	"github.com/go-chassis/ygx/libs/i18n"

	"github.com/go-chassis/ygx/schemas/gateway/gclient"
)

var (
	reponse      = common.Response{}
	replyCacheCh = make(chan GrpcReply, 10000)
)

// RestFulGateWay  rest 网关
type RestFulGateWay struct {
}

func init() {
	i18n.SetPlatLang("app", "cn")
}

// 前台 app 使用网关通信   grpc做内部的通信
//header service grpc:ygx_usercenter_user_grpc-server/UserLogin:0.0.1
//header token
//header Content-Type  application/json
func (rg *RestFulGateWay) YgxAppGateway(b *rf.Context) {
	//token := b.ReadHeader("UserToken")

	service := b.ReadHeader("service")
	serviceOption := strings.Split(service, ":")
	fmt.Printf("serviceOption:%s \n", serviceOption)
	if serviceOption[0] != "grpc" {
		openlogging.Warn("The app api system only support  grpc  at present")
		//Message.
		b.WriteJSON(1, "The  app api system only support  grpc  at present, protocal is"+serviceOption[0], "application/json")
		return
	}
	callParemeter := strings.Split(serviceOption[1], "/")
	serverName := callParemeter[0]
	method := callParemeter[1]

	version := "lastest"
	if serviceOption[2] != "" {
		version = serviceOption[2]
	}

	body, _ := ioutil.ReadAll(b.ReadRequest().Body)
	// 此处代码也 ok
	reply, err := ChoiceGrpcClient(serverName, method, "", version, body)
	reponseCall := common.Response{}
	if err == nil {
		reponseCall.Status = int(i18n.FAILURE)
		reponseCall.Data = reply
		reponseCall.Message = i18n.LoadMessage(i18n.FAILURE)
	} else {
		reponseCall.Status = int(i18n.SUCCESS)
		reponseCall.Message = fmt.Sprint("%s", err)
		reponseCall.Message = i18n.LoadMessage(i18n.SUCCESS)
	}
	b.WriteJSON(reponseCall, "application/json", "")

	//go grpcClientProvide(serverName, method, "", version, body, b, replyCacheCh)
	//// todo eroor
	////go grpcClientConsumer(replyCacheCh)

	return
}

// 本来使用 go 加 提供者和消费者模式做。bug 出现
//func test() {
//  replyCacheCh := make(chan GrpcReply, 1000)
//	go grpcClientProvide(serverName, method, "", version, body, replyCacheCh)
//	go grpcClientConsumer(b, replyCacheCh)
// replyInfo.b.WriteJSON(reponseCall, "application/json", "")
//	return
//}
type GrpcReply struct {
	reply interface{}
	err   error
	b     *rf.Context
}

func grpcClientProvide(serviceName string, method string, token string, version string, parameter []byte, b *rf.Context, in chan GrpcReply) {
	var (
		c         gclient.ClientInit
		request   interface{}
		reply     interface{}
		schemaId  string
		grpcReply GrpcReply
	)
	grpcReply.b = b

	switch serviceName {
	case "ygx_usercenter_user_grpc-server":
		c = new(gclient.YgxUserCenterUserInit)
		request, reply, schemaId = gclient.NewClientInit(c, method)
		break
	case "":
	default:
		grpcReply.reply = nil
		grpcReply.err = errors.New(i18n.LoadMessage(i18n.NO_SERVER))
		in <- grpcReply
		break
		return
	}

	jerr := json.Unmarshal(parameter, request)
	if jerr != nil {
		grpcReply.reply = nil
		grpcReply.err = jerr
		in <- grpcReply
		return
	}

	//Invoke with microservice name, schema ID and operation ID
	if err := core.NewRPCInvoker().Invoke(context.Background(), serviceName, schemaId, method,
		request, reply, core.WithProtocol("grpc")); err != nil {
		lager.Logger.Error("error" + err.Error())

		grpcReply.reply = nil
		grpcReply.err = err
		in <- grpcReply
		return
	} else {
		grpcReply.reply = reply
		grpcReply.err = nil
	}

	in <- grpcReply
	return
}

func grpcClientConsumer(out chan GrpcReply) {
	reponseCall := common.Response{}

	replyInfo := new(GrpcReply)
	*replyInfo = <-out

	if replyInfo.err == nil {
		reponseCall.Status = 0
		reponseCall.Data = replyInfo.reply
	} else {
		reponseCall.Status = 1
		reponseCall.Message = fmt.Sprint("%s", replyInfo.err)
	}

	//replyInfo.b.WriteJSON(reponseCall, "application/json", "")
	return
}

// service格式    rest:post:ygx_admin_user/updateUser:0.01   通讯协议:微服务名称:method
func (rg *RestFulGateWay) YgxAdminGateway(b *rf.Context) {
	token := b.ReadHeader("Token")

	service := b.ReadHeader("service")
	serviceOption := strings.Split(service, ":")
	fmt.Printf("serviceOption:%s \n", serviceOption)
	if serviceOption[0] != "rest" {
		openlogging.Warn("The system only support  rest  at present")
		//Message.
		b.WriteJSON(1, "The system only support  rest  at present, protocal is"+serviceOption[0], "application/json")
		return
	}

	httpMethod := serviceOption[1]
	serviceName := serviceOption[2]
	version := serviceOption[3]
	//version := b.ReadHeader("version")
	if version == "" {
		version = "lastest"
	}
	fmt.Printf("httpMethod:%s  serviceName:%s  \n", httpMethod, serviceName)

	errorMes := verifyHeader(httpMethod, serviceName)
	if errorMes != "" {
		b.WriteJSON(reponse.GetErrorResponse(1, errorMes), "application/json", "")
		return
	}

	urlOption := strings.Split(serviceName, "/")
	_, err := endpoint.GetEndpoint("default", urlOption[0], version)
	if err != nil {
		openlogging.Warn("failed to find config center endpoints, err: " + err.Error())
		b.WriteJSON(reponse.GetErrorResponse(1, "没有找到配置实例,请先启动"+serviceOption[0]+"版本:"+version+" server 错误:"+err.Error()), "application/json", "")
		return
	}

	var returnRes = make(chan *http.Request)
	var wg sync.WaitGroup
	wg.Add(1)

	reponseCall := common.Response{}
	go func() {
		// 此处一定要大写 跨域 否则header不允许通过
		httpMethod = strings.ToUpper(strings.TrimSpace(httpMethod))
		url := "http://" + serviceName
		body, _ := ioutil.ReadAll(b.ReadRequest().Body)

		req, err := rest.NewRequest(string(httpMethod), url, nil)
		httputil.SetBody(req, body)

		if err != nil {
			lager.Logger.Error("new request failed.")
			panic("ERROR:call " + url + " failed.")
		} else {
			lager.Logger.Error("new request ok.")

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Token", token)
			returnRes <- req
		}

		lager.Logger.Error(fmt.Sprintf("req:%+v \n", req))

		defer wg.Done()
	}()

	tmp := <-returnRes
	// lager.Logger.Error(fmt.Sprintf("tmp:%+v \n", tmp))
	ctx := context.WithValue(context.TODO(), corecommon.ContextHeaderKey{}, map[string]string{})

	resp, err := core.NewRestInvoker().ContextDo(ctx, tmp)
	if err != nil {
		lager.Logger.Error("do request failed." + err.Error())
		return
	}

	// fmt.Printf("resp:%s n", resp)
	if resp.StatusCode != http.StatusOK {
		lager.Logger.Error(fmt.Sprintf("http服务器错误 http状态码: %d.\n", resp.StatusCode))
		reponseCall.Status = 1
		reponseCall.Message = fmt.Sprintf("服务器错误,错误码为:http  %d.\n", resp.StatusCode)
	} else {
		json.Unmarshal(readResponseBody(resp), &reponseCall)
	}

	b.WriteJSON(reponseCall, "application/json", "")
	defer resp.Body.Close()

	wg.Wait()
}

func verifyHeader(method string, service string) (errorMes string) {
	var isExist = false
	for _, m := range common.HttpMethodSet {
		if strings.ToLower(method) == m {
			isExist = true
			break
		}
	}

	if isExist == false {
		errorMes = "请求的方法不能够为空"
	}

	// todo 从服务中心获取 注册的服务名称  暂时做空的判断
	if service == "" {
		errorMes = fmt.Sprintf("%s \n 或者%s", errorMes, "请求的微服务名称不能够为空")
	}

	return errorMes
}

func readResponseBody(resp *http.Response) []byte {
	if resp != nil && resp.Body != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			lager.Logger.Error(fmt.Sprintf("read body failed: %s", err.Error()))
			return nil
		}
		return body
	}
	lager.Logger.Error("response body or response is nil")
	return nil
}

// URLPatterns helps to respond for corresponding API calls
// ygx app   表示壹共享app网关
// ygx admin 表示壹共享后台网关
func (rg *RestFulGateWay) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodPost, Path: "/ygx/app", ResourceFunc: rg.YgxAppGateway,
			Metadata: map[string]interface{}{
				"tags": []string{"token", "method", "service", "func_name"},
			},
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodPost, Path: "/ygx/admin", ResourceFunc: rg.YgxAdminGateway,
			Metadata: map[string]interface{}{
				"tags": []string{"token", "method", "service", "func_name"},
			},
			Returns: []*rf.Returns{{Code: 200}}},
	}
}
