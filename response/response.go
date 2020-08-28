package response

type Auth struct {
	Owner       interface{} `json:"owner"`
	AccessToken string      `json:"access_token"`
}

type Error struct {
	Message string `json:"message"`
}
