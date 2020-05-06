package example

import (
	"log"
	"net/http"

	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
	rf "github.com/go-chassis/go-chassis/server/restful"

	"github.com/go-mesh/openlogging"
)

//TracingHello is a struct
type TracingHello struct {
}

//Trace is a method
func (r *TracingHello) Trace(b *rf.Context) {
	log.Println("tracing===", b.Ctx)
	req, err := rest.NewRequest("GET", "http://RESTServerB/sayhello/world", nil)
	if err != nil {
		openlogging.GetLogger().Errorf("err %s", err)

		b.WriteError(500, err)
		return
	}

	resp, err := core.NewRestInvoker().ContextDo(b.Ctx, req)
	if err != nil {
		b.WriteError(500, err)
		return
	}
	defer resp.Body.Close()

	openlogging.GetLogger().Errorf("-----------------------req %s", httputil.ReadBody(resp))
	b.Write(httputil.ReadBody(resp))
}

//URLPatterns helps to respond for corresponding API calls
func (r *TracingHello) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/trace", ResourceFunc: r.Trace},
	}
}
