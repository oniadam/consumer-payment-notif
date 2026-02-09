package models

type NotifPaymentWa struct {
	CustomerName string `json:"customerName"`
	AggrNo       string `json:"aggrNo"`
	// Amount           string `json:"amount"`
	WaNo             string  `json:"waNo"`
	Senddtm          string  `json:"sendDtm"`
	Sendby           string  `json:"sendby"`
	Templatecode     string  `json:"templatecode"`
	TotalPaid        float64 `json:"totalPaid"`
	TransactionSrc   string  `json:"transactionSrc"`
	Paymentmetodcode string  `json:"paymentmetodcode"`
	Refno            string  `json:"refno"`
	RefNoWa          string  `json:"refNoWa"`
	Filepath         string  `json:"filepath"`
	Flagreversal     string  `json:"flagreversal"`
	Createdby        string  `json:"createdby"`
	Createddtm       string  `json:"createddtm"`
}

type InstReqToMeta struct {
	MessagingProduct string   `json:"messaging_product"`
	To               string   `json:"to"`
	Type             string   `json:"type"`
	Text             TextBody `json:"text"`
}

type TextBody struct {
	Body string `json:"body"`
}

type ResFrMeta struct {
	MessagingProduct string `json:"messaging_product"`
	Contacts         []struct {
		Input string `json:"input"`
		WaID  string `json:"wa_id"`
	} `json:"contacts"`
	Messages []struct {
		ID string `json:"id"`
	} `json:"messages"`
}
