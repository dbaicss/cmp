package routes
import (
	"net/http"

	"cmp-server/auth"
	"cmp-server/logic"
	"github.com/gorilla/mux"
	"cmp-server/api"
)

type Route struct {
	Method     string
	Pattern    string
	Handler    http.HandlerFunc
	Middleware mux.MiddlewareFunc
}

var routes []Route

func init() {
	register("GET", "/api/pod/list", api.HandlePodList, auth.TokenMiddleware)
	register("GET", "/api/crontab/list", api.HandleCrontabList, auth.TokenMiddleware)
	register("GET", "/api/asset/list", api.HandleAssetList, auth.TokenMiddleware)
	register("GET", "/api/server/list", api.HandleServerList, auth.TokenMiddleware)
	register("GET", "/api/service/list", api.HandleAssetServiceInfoList, auth.TokenMiddleware)
	register("GET", "/api/audit/list", api.HandleAuditList, auth.TokenMiddleware)
	register("GET", "/api/namespace/list", api.HandleNameSpaceList, auth.TokenMiddleware)
	register("GET", "/api/log/list", api.HandleLogList, auth.TokenMiddleware)
	register("GET","/ssh",api.WsHandler,nil)
	register("POST", "/api/user/register", logic.Register, nil)
	register("POST", "/api/user/login", logic.Login, nil)
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	for _, route := range routes {
		r := router.Methods(route.Method).
			Path(route.Pattern)
		if route.Middleware != nil {
			r.Handler(route.Middleware(route.Handler))
		} else {
			r.Handler(route.Handler)
		}
	}
	return router
}

func register(method, pattern string, handler http.HandlerFunc, middleware mux.MiddlewareFunc) {
	routes = append(routes, Route{method, pattern, handler, middleware})
}

