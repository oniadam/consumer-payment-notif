package models

import "gopkg.in/guregu/null.v4"

type GetWaNoRes struct {
	CustomerCode string      `json:"customerCode"`
	Fullname     string      `json:"fullname"`
	WaNo         null.String `json:"waNo"`
}
