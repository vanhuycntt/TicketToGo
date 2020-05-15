package delegate

import "TicketToGo/cmd/tests/interfaces/broker"

type restAPIHandler struct {
	postHandler RestPostHandler
	getHandler  RestGetHandler
}

func New() broker.Handler {
	rapi := restAPIHandler{}
}

type RestPostHandler interface {
	Post(body string, headers map[string]interface{})
	Close() error
}

type RestGetHandler interface {
	Get(params map[string]interface{}) error
}
