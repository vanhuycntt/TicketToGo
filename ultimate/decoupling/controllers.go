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
	return nil
}

func DownloadExcel(downloader func(id string) ([]byte, error)) Handler {
	return func(ctx context.Context, res http.ResponseWriter, req *http.Request) error {
		id := req.URL.Query().Get("id")
		if id == "" {
			_, _ = res.Write([]byte(`{"msg":"invalid param"}`))
			res.WriteHeader(http.StatusBadRequest)
			return nil
		}

		bytes, err := downloader(id)
		if err != nil {
			_, _ = res.Write([]byte(`"msg":"not found""`))
			res.WriteHeader(http.StatusBadRequest)
			return err
		}
		if _, err := res.Write(bytes); err != nil {
			_, _ = res.Write([]byte(`"msg":"not found""`))
			res.WriteHeader(http.StatusInternalServerError)
			return err
		}

		res.Header().Set("content-type", "application/octet-stream")
		return nil

	}
}
