package http

import (
	"net"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/request"
	"github.com/codecrafters-io/http-server-starter-go/app/response"
)

type Handler func(net.Conn, request.Request)

type route struct {
	method  string
	path    string
	handler Handler
}

type App struct {
	routes []route
}

func New() *App {
	return &App{}
}

func (app *App) Get(path string, handler Handler) {
	app.add("GET", path, handler)
}

func (app *App) Post(path string, handler Handler) {
	app.add("POST", path, handler)
}

func (app *App) Serve(conn net.Conn, req request.Request) {
	for _, route := range app.routes {
		if route.method == req.Method && match(route.path, req.Path) {
			route.handler(conn, req)
			return
		}
	}

	response.Write(conn, response.Response{
		StatusCode: 404,
	})
}

func (app *App) add(method string, path string, handler Handler) {
	app.routes = append(app.routes, route{
		method:  method,
		path:    path,
		handler: handler,
	})
}

func match(routePath string, requestPath string) bool {
	if strings.HasSuffix(routePath, "/*") {
		prefix := strings.TrimSuffix(routePath, "/*")
		return strings.HasPrefix(requestPath, prefix+"/")
	}

	return routePath == requestPath
}
