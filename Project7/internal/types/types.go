package types

type Student struct {
	Id    int64 
	Name  string `validate:"required"`
	Email string `validate:"required"`
	Age   int `validate:"required"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}
