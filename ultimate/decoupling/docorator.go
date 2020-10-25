package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

type Handler func(ctx context.Context, rs http.ResponseWriter, rq *http.Request) error

type Middleware func(handler Handler) Handler

type User struct {
	session string
	loginID string
}

type RestHandler struct {
	pathsToHandler map[string]http.Handler
}

func (rh RestHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	pathHandler, ok := rh.pathsToHandler[req.URL.Path]
	if !ok {
		res.Write([]byte(`{msg: "not found", code :404}`))
		res.WriteHeader(http.StatusNotFound)

		return
	}
	pathHandler.ServeHTTP(res, req)
}

func (rh *RestHandler) Handle(path string, handler Handler, mids ...Middleware) {
	handler = rh.filterHandler(handler, mids...)
	h := func(res http.ResponseWriter, req *http.Request) {
		err := handler(context.Background(), res, req)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
		}
	}
	rh.pathsToHandler[path] = http.HandlerFunc(h)
}

func (rh *RestHandler) filterHandler(handler Handler, mids ...Middleware) Handler {
	for idx := len(mids) - 1; idx >= 0; idx-- {
		handler = mids[idx](handler)

	}
	return handler
}

func main() {
	rh := RestHandler{
		pathsToHandler: map[string]http.Handler{},
	}

	rh.Handle("/api/v1/about", RequestJWTToken, Auth(User{}))
	rh.Handle("/api/v1/sign_in", SignIn)
	rh.Handle("/api/v1/excel/invoices/{code}", DownloadExcel(func(id string) (bytes []byte, e error) {
		return []byte(""), nil
	}))

	rh.Handle("/api/v1/excel/invoices/{code}", DownloadExcel(func(id string) (bytes []byte, e error) {
		return []byte(""), nil
	}))

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: rh,
	}
	server.RegisterOnShutdown(func() {

	})

	var syncWait sync.WaitGroup
	syncWait.Add(1)
	go func() {
		defer syncWait.Done()
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("server error", err)
		}
	}()
	if err := server.Shutdown(context.Background()); err != nil {

	}
	syncWait.Wait()
}
