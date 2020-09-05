package main

import (
	"context"
	"net/http"
)

func RequestJWTToken(ctx context.Context, res http.ResponseWriter, req *http.Request) error {
	res.Write([]byte(`{msg:"about page"}`))
	res.WriteHeader(http.StatusOK)
	return nil
}

func SignIn(ctx context.Context, res http.ResponseWriter, req *http.Request) error {

}
