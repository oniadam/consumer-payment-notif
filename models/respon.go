package models

type Respons struct {
	ResponseCode      string      `json:"responseCode"`
	ResponseMessage   string      `json:"responseMessage"`
	ResponseTimestamp string      `json:"responseTimestamp"`
	Errors            string      `json:"errors"`
	Data              interface{} `json:"data"`
}
