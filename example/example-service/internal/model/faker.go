package model

type ResponseFaker struct {
	BaseResponseFaker
	Data []DataFaker `json:"data"`
}

type DataFaker struct {
	Id        int    `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Birthday  string `json:"birthday"`
	Gender    string `json:"gender"`
	Website   string `json:"website"`
	Image     string `json:"image"`
}
