package server

type CreateUserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  string `json:"age"`
	Job  string `json:"job"`
}

type GetUserResponse struct {
	Name string `json:"name"`
	Age  string `json:"age"`
	Job  string `json:"job"`
}
