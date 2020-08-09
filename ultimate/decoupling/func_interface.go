package main

import "net/http"

type Handler func(rs http.ResponseWriter, rq *http.Request) error

type Middleware func(handler Handler) Handler

type User struct {
}

func Authenticate(user User) Middleware {
	md := func(after Handler) Handler {

		return func(rs http.ResponseWriter, rq *http.Request) error {

			return after(rs, rq)
		}
	}
	return md
}
func main() {

}
