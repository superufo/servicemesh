package admin

import (

	// client
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/go-chassis/foundation/token"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
	rf "github.com/go-chassis/go-chassis/server/restful"
	"github.com/go-chassis/ygx/libs/common"
	"github.com/go-chassis/ygx/libs/common/mapstructure"
	"github.com/go-chassis/ygx/libs/common/redis"
	"github.com/go-chassis/ygx/libs/i18n"

	"github.com/go-mesh/openlogging"

	AdminModel "github.com/go-chassis/ygx/libs/model/admin"

	_ "github.com/go-chassis/go-chassis/bootstrap"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"
	corecommon "github.com/go-chassis/go-chassis/core/common"

	"github.com/go-chassis/go-chassis/core/lager"
)

func init() {
	i18n.SetPlatLang("admin", "cn")
}

type RestFulAdmin struct {
}

// Root Sayhello is a method used to reply user with hello
func (r *RestFulAdmin) Root(b *rf.Context) {
	b.Write([]byte(fmt.Sprintf("x-forwarded-host %s", b.ReadRequest().Host)))
}

// GetUserInfo 可以使用json存储   jsonStr, _ := json.Marshal(hashMap) redis.HashSet("XfcAdmin", strconv.Itoa(user_id), string(jsonStr))
func (r *RestFulAdmin) GetUserInfo(b *rf.Context) {
	var userMap map[string]interface{}
	reponse := common.Response{}
	tokenstr := b.ReadHeader("Token")

	if tokenstr == "" {
		reponse.Status = 1
		reponse.Message = "Token 不能够为空"
		openlogging.Error("Token 不能够为空")

		b.WriteJSON(reponse, "application/json", "")
		return
	}

	t := token.Token{tokenstr}

	claims, err := t.TokenVerify()
	user_id := claims.UserId

	userData := AdminModel.XfcAdmin{
		AdminId:     0,
		UserName:    "",
		DepId:       0,
		PosId:       0,
		Mobile:      "",
		Email:       "",
		Password:    "",
		Ecsalt:      "",
		AddTime:     0,
		Lastlogin:   0,
		Lastip:      "",
		NavList:     "",
		LangType:    "",
		AgencyId:    0,
		SuppliersId: 0,
		TodoList:    "",
		RoleId:      0,
		IsLock:      0,
		IsSales:     0,
		IsGovsales:  0,
	}

	hkey := "XfcAdmin:" + strconv.Itoa(user_id)
	hashMap, err := redis.Client.HGetAll(hkey).Result()
	// redis 存在hash值
	if err != nil {
		for i, v := range hashMap {
			userMap[i] = v
		}
		mapstructure.Decode(userMap, &userData)
	} else {
		userData.AdminId = user_id
		common.DB.Where(&userData).Find(&userData)

		userMap := common.StructToMap(userData)
		redis.BatchHashSet(hkey, userMap)
	}

	UserInfo := AdminModel.UserInfo{}
	UserInfo.User = userData

	if err != nil {
		reponse.Status = 1
		reponse.Message = "获取用户信息失败"
		openlogging.Error("获取用户信息失败")
	} else {
		reponse.Status = 0
		reponse.Data = UserInfo
		reponse.Message = "获取用户信息成功"
		openlogging.Info("获取用户信息成功")
	}

	b.WriteJSON(reponse, "application/json", "")
	return
}

// UpdateUserInfo  更新用户
func (r *RestFulAdmin) UpdateUserInfo(b *rf.Context) {
	userId := b.ReadHeader("userid")
	UserInfo := AdminModel.UserInfo{}
	userData := AdminModel.XfcAdmin{}
	reponse := common.Response{}

	fmt.Printf("userId %s", userId)
	// 更新用户信息
	b.ReadEntity(&userData)
	userMap := common.StructToMapNoEmpty(userData)
	fmt.Printf("userMap %+v", userData)
	fmt.Printf("userMap %+v", userMap)
	common.DB.Model(&userData).Where("admin_id = ?", userId).Updates(userMap)
	common.DB.Where("admin_id = ?", userId).First(&userData)
	UserInfo.User = userData

	reponse.Data = UserInfo
	reponse.Status = 0
	reponse.Message = "更新用户信息成功"

	b.WriteJSON(reponse, "application/json", "")
	return
}

type SayHello struct {
	AdminId  int    `json:"userid"`
	UserName string `json:"user_name"`
}

type Test struct {
	Test string `json:"test"`
}

func (r *RestFulAdmin) Call(b *rf.Context) {
	token := b.ReadHeader("Token")
	userName := b.ReadHeader("username")

	var returnRes = make(chan *http.Request)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		body, _ := ioutil.ReadAll(b.ReadRequest().Body)

		//body := common.StrToBytes("{\"userid\":\"admin@ioa365.com\",\"user_name\":\"aaaaaaaaaaaa\"}")
		lager.Logger.Error(fmt.Sprintf("body:%s \n", string(body)))
		req, err := rest.NewRequest("POST", "http://ygx_example_demo_rest_server/sayhello/"+userName, nil)

		httputil.SetBody(req, body)
		if err != nil {
			lager.Logger.Error("new request failed.")
			returnRes <- nil
			panic(" call sayhello request failed. ")
		} else {
			lager.Logger.Error("new request ok.")

			req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			req.Header.Set("Token", token)

			returnRes <- req
		}
		wg.Done()
	}()

	tmp := <-returnRes
	lager.Logger.Error(fmt.Sprintf("tmp:%+v \n", tmp))
	ctx := context.WithValue(context.TODO(), corecommon.ContextHeaderKey{}, map[string]string{
		"user": userName,
	})
	resp, err := core.NewRestInvoker().ContextDo(ctx, tmp)

	if err != nil {
		lager.Logger.Error("do request failed.")
		return
	}
	reponseCall := common.Response{}
	json.Unmarshal(httputil.ReadBody(resp), &reponseCall)

	fmt.Printf("reponseCall.Data %+V", reponseCall.Data)
	//var test, _ = reponseCall.Data.(Test)

	lager.Logger.Info("REST Server sayhello[GET]: " + string(httputil.ReadBody(resp)))
	reponse := common.Response{}
	if reponseCall.Status == int(i18n.SUCCESS) {
		reponse.Data = reponseCall.Data
		reponse.Status = int(i18n.SUCCESS)
		reponse.Message = "调用成功"
	} else {
		reponse.Data = nil
		reponse.Status = int(i18n.FAILURE)
		reponse.Message = "调用失败"
	}

	b.WriteJSON(reponse, "application/json", "")
	defer resp.Body.Close()

	wg.Wait()

	return
}

// Call  调用服务 RESTServer
func (r *RestFulAdmin) CallPhp(b *rf.Context) {
	var returnRes = make(chan *http.Request)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		req, err := rest.NewRequest("GET", "http://ygx_example_rest/callphp", nil)
		if err != nil {
			lager.Logger.Error("new request failed.")
			returnRes <- nil
			panic(" call sayhello request failed. ")
		} else {
			lager.Logger.Error("new request ok.")

			req.Header.Set("Content-Type", "application/json")
			returnRes <- req
		}
		wg.Done()
	}()

	ctx := context.WithValue(context.TODO(), corecommon.ContextHeaderKey{}, map[string]string{})
	resp, err := core.NewRestInvoker().ContextDo(ctx, <-returnRes)

	if err != nil {
		lager.Logger.Error("do request failed.")
		return
	}
	reponseCall := common.Response{}
	json.Unmarshal(httputil.ReadBody(resp), &reponseCall)

	fmt.Printf("reponseCall.Data %+V", reponseCall.Data)
	test, _ := reponseCall.Data.(string)

	fmt.Printf("%+s", test)
	lager.Logger.Info("REST Server sayhello[GET]: " + string(httputil.ReadBody(resp)))
	reponse := common.Response{}
	if reponseCall.Status == 0 {
		reponse.Data = reponseCall.Data
		reponse.Status = 0
		reponse.Message = "调用成功"
	} else {
		reponse.Data = nil
		reponse.Status = reponseCall.Status
		reponse.Message = "调用失败"
	}

	b.WriteJSON(reponse, "application/json", "")
	defer resp.Body.Close()
	wg.Wait()

	return
}

// URLPatterns helps to respond for corresponding API calls
func (r *RestFulAdmin) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/", ResourceFunc: r.Root,
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodGet, Path: "/getInfo", ResourceFunc: r.GetUserInfo,
			Metadata: map[string]interface{}{
				"tags": []string{"users", "test"},
			},
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodPost, Path: "/updateUser", ResourceFunc: r.UpdateUserInfo,
			Metadata: map[string]interface{}{
				"tags": []string{"users", "test"},
			},
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodGet, Path: "/callphp", ResourceFunc: r.CallPhp,
			Metadata: map[string]interface{}{
				"tags": []string{"users", "test"},
			},
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodPost, Path: "/call", ResourceFunc: r.Call,
			Metadata: map[string]interface{}{
				"tags": []string{"users", "test"},
			},
			Returns: []*rf.Returns{{Code: 200}}},
	}
}
