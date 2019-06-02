package wssv

type RegisterReq struct {
	Username string
	Password string
}

type RegisterResp struct {
	Error string
}

type LoginReq struct {
	Username string
	Password string
}

type LoginResp struct {
	Token string
	Error string
}

type MessageReq struct {
	Authorization string
	Message       string
}

type MessageResp struct {
	Author  string
	Message string
}
