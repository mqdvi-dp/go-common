package model

type BaseResponseFaker struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Total  int    `json:"total"`
}
