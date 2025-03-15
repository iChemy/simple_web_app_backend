package handler

type Handlers struct {
	Ping *PingController
	User *UserController
}

func NewHandlers() Handlers {
	return Handlers{
		Ping: &PingController{},
	}
}
