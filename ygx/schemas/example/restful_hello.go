package example

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"

	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
	rf "github.com/go-chassis/go-chassis/server/restful"
	"github.com/go-chassis/ygx/libs/common"

	corecommon "github.com/go-chassis/go-chassis/core/common"
	"github.com/go-mesh/openlogging"
)

var num = rand.Intn(100)

// RestFulHello is a struct used for implementation of restfull hello program
type RestFulHello struct {
}

//Sayhello is a method used to reply user with hello
func (r *RestFulHello) Root(b *rf.Context) {
	b.Write([]byte(fmt.Sprintf("x-forwarded-host %s", b.ReadRequest().Host)))
}

type Test struct {
	Test string `json:"test"`
}

// Sayhello is a method used to reply user with hello
func (r *RestFulHello) Sayhello(b *rf.Context) {
	// id := b.ReadPathParameter("userid")
	// log.Printf("get user id: " + id)
	// log.Printf("get user name: " + b.ReadRequest().Header.Get("user"))

	// b.Write([]byte(fmt.Sprintf("user %s from %d", id, num)))
	body, err := ioutil.ReadAll(b.ReadRequest().Body)
	//body, err := base64.StdEncoding.DecodeString(string(base64Body))
	if err != nil {
		lager.Logger.Error(fmt.Sprintf("body:%s \n", err.Error()))
	}
	lager.Logger.Error(fmt.Sprintf("Sayhello body:%s \n", string(body)))

	var s = Test{
		b.ReadRequest().Header.Get("user"),
	}

	reponse := common.Response{}
	reponse.Data = s
	reponse.Status = 0
	reponse.Message = "执行成功"
	fmt.Printf("reponse %+V", reponse)

	openlogging.GetLogger().Errorf("-----------------------req %s", reponse)
	b.WriteJSON(reponse, "application/json", "")
}

// Sayhi is a method used to reply user with hello world text
func (r *RestFulHello) Sayhi(b *rf.Context) {
	result := struct {
		Name string
	}{}
	err := b.ReadEntity(&result)
	if err != nil {
		b.Write([]byte(err.Error() + ":hello world"))
		return
	}
	b.Write([]byte(result.Name + ":hello world"))
	return
}

// SayJSON is a method used to reply user hello in json format
func (r *RestFulHello) SayJSON(b *rf.Context) {
	reslut := struct {
		Name string
	}{}

	err := b.ReadEntity(&reslut)
	if err != nil {
		b.WriteHeaderAndJSON(http.StatusInternalServerError, reslut, "application/json")
		return
	}

	reslut.Name = "hello:" + reslut.Name
	b.WriteJSON(reslut, "application/json", "")
	return
}

func (r *RestFulHello) CallPhp(b *rf.Context) {
	var returnRes = make(chan *http.Request)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		//req, err := rest.NewRequest("GET", "http://mesher-consumer:8088/client.php", nil)
		req, err := rest.NewRequest("GET", "http://hellomesher:80/api.php?name=pen", nil)
		if err != nil {
			lager.Logger.Error("new request failed.")
			returnRes <- nil
			panic(" call sayhello request failed. ")
		} else {
			lager.Logger.Error("new request ok.")

			req.Header.Set("Content-Type", "multipart/form-data")
			returnRes <- req
		}
		wg.Done()
	}()

	ctx := context.WithValue(context.TODO(), corecommon.ContextHeaderKey{}, map[string]string{
		"name": "pen",
	})

	resp, err := core.NewRestInvoker().ContextDo(ctx, <-returnRes)
	if err != nil {
		lager.Logger.Error("do request failed.")
		return
	}

	reponse := common.Response{}
	reponse.Data = string(httputil.ReadBody(resp))
	reponse.Status = 0
	reponse.Message = "执行成功"
	b.WriteJSON(reponse, "application/json", "")

	defer resp.Body.Close()

	wg.Wait()
}

// URLPatterns helps to respond for corresponding API calls
func (r *RestFulHello) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/", ResourceFunc: r.Root,
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodPost, Path: "/sayhello/{userid}", ResourceFunc: r.Sayhello,
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodGet, Path: "/sayhello/{userid}", ResourceFunc: r.Sayhello,
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodPost, Path: "/sayhi", ResourceFunc: r.Sayhi,
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodPost, Path: "/sayjson",
			ResourceFunc: r.SayJSON,
			Metadata: map[string]interface{}{
				"tags": []string{"users", "test"},
			},
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodGet, Path: "/callphp", ResourceFunc: r.CallPhp,
			Returns: []*rf.Returns{{Code: 200}}},
	}
}

// RestFulMessage is a struct used to implement restful message
type RestFulMessage struct {
}

// Saymessage is used to reply user with his name
func (r *RestFulMessage) Saymessage(b *rf.Context) {
	id := b.ReadPathParameter("name")

	b.Write([]byte("get name: " + id))
}

// Sayhi is a method used to reply request user with hello world text
func (r *RestFulMessage) Sayhi(b *rf.Context) {
	reslut := struct {
		Name string
	}{}
	err := b.ReadEntity(&reslut)
	if err != nil {
		b.Write([]byte(err.Error() + ":hello world"))
		return
	}
	b.Write([]byte(reslut.Name + ":hello world"))
	return
}

// Sayerror is a method used to reply request user with error
func (r *RestFulMessage) Sayerror(b *rf.Context) {
	b.WriteError(http.StatusInternalServerError, errors.New("test hystric"))
	return
}

// URLPatterns helps to respond for corresponding API calls
func (r *RestFulMessage) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/saymessage/{name}", ResourceFunc: r.Saymessage},
		{Method: http.MethodPost, Path: "/sayhimessage", ResourceFunc: r.Sayhi},
		{Method: http.MethodGet, Path: "/sayerror", ResourceFunc: r.Sayhi},
	}
}

// Hello is a struct used for implementation of restfull hello program
type Hello struct{}

// Hello
func (r *Hello) Hello(b *rf.Context) { b.Write([]byte("hi from hello")) }

// URLPatterns helps to respond for corresponding API calls
func (r *Hello) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/hello", ResourceFunc: r.Hello},
	}
}

// Legacy is a struct
type Legacy struct{}

// Do
func (r *Legacy) Do(b *rf.Context) { b.Write([]byte("hello from legacy")) }

// URLPatterns helps to respond for corresponding API calls
func (r *Legacy) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/legacy", ResourceFunc: r.Do},
	}
}

// Legacy is a struct
type Admin struct{}

// Do
func (r *Admin) Do(b *rf.Context) { b.Write([]byte("hello from admin")) }

//URLPatterns helps to respond for corresponding API calls
func (r *Admin) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/admin", ResourceFunc: r.Do},
	}
}
