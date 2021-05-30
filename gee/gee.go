package gee

import (
	"log"
	"net/http"
	"strings"
)

// 	HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	engine      *Engine       // all groups share a Engine instance
}

// Engine implement the interface of ServerHTTP
type Engine struct {
	*RouterGroup // 将Engine作为最顶层的分组，也就是说Engine拥有RouterGroup所有的能力。
	router       *router
	groups       []*RouterGroup // store all groups
}

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)

	return newGroup
}

func (group *RouterGroup) addRoute(method, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (group *RouterGroup) GET(comp string, handler HandlerFunc) {
	group.addRoute("GET", comp, handler)
}

// POST defines the method to add POST request
func (group *RouterGroup) POST(comp string, handler HandlerFunc) {
	group.addRoute("POST", comp, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// Use is define to add middleware to the group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := NewContext(w, r)
	c.middlewares = middlewares
	engine.router.handle(c)
}
