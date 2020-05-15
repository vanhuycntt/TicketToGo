package broker

type Handler interface {
	PostHandler
	GetHandler
}

type PostHandler interface {
	Post(body string, headers map[string]interface{})
}

type GetHandler interface {
	Get(params map[string]interface{}) error
}
