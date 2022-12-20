package app

type Credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type JWTTokenResp struct {
	Token string `json:"token"`
}

type Message struct {
	Message string `json:"message"`
}

type IDResponse struct {
	ID int64 `json:"id"`
}
