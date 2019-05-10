package models

type MSG struct {
	Message string `json:"Message"`
	RequestId string `json:"RequestId"`
	BizId string `json:"BizId"`
	Code string `json:"Code"`
}