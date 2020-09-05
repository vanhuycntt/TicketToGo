package main

import (
	"context"
	"errors"
	"net/http"
)

func Auth(user User) Middleware {
	md := func(after Handler) Handler {
		return func(ctx context.Context, rs http.ResponseWriter, rq *http.Request) error {

			if user.loginID == "unknown" {
				return errors.New("unknown user")
			}

			if user.session == "timeout" {
				return errors.New("login timeout")
			}
			return after(ctx, rs, rq)
		}
	}
	return md
}
