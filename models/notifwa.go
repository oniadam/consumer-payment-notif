package models

type NotifPaymentWa struct {
	CustomerName      string  `json:"customerName"`
	AggrNo            string  `json:"aggrNo"`
	IdNo              string  `json:"idNo"`
	WaNo              string  `json:"waNo"`
	CustomerServiceNo string  `json:"customerServiceNo"`
	Senddtm           string  `json:"sendDtm"`
	Sendby            string  `json:"sendby"`
	Templatecode      string  `json:"templatecode"`
	LanguageCode      string  `json:"languageCode"`
	TotalPaid         float64 `json:"totalPaid"`
	TransactionSrc    string  `json:"transactionSrc"`
	Paymentmetodcode  string  `json:"paymentmetodcode"`
	Refno             string  `json:"refno"`
	RefNoWa           string  `json:"refNoWa"`
	Filepath          string  `json:"filepath"`
	Flagreversal      string  `json:"flagreversal"`
	Createdby         string  `json:"createdby"`
	Createddtm        string  `json:"createddtm"`
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

type NotifPaymentTemplateToMetaFbReq struct {
	MessagingProduct string                `json:"messaging_product"`
	To               string                `json:"to"`
	Type             string                `json:"type"`
	Template         NotifPaymentTemplates `json:"template"`
}

type NotifPaymentTemplates struct {
	Name        string       `json:"name"`
	LanguageS   LanguageS    `json:"language"`
	ComponentsS []ComponentS `json:"components"`
}

type LanguageS struct {
	Code string `json:"code"`
}

type ComponentS struct {
	Type       string       `json:"type"`
	Parameters []ParameterS `json:"parameters"`
}

type ParameterS struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type NotifMyestaReq struct {
	IdNo             string `json:"idNo"`
	Token            string `json:"token"`
	Title            string `json:"title" `
	ShortDescription string `json:"shortDescription"`
	FullDescription  string `json:"fullDescription" db:"full_description"`
	ResponseFcm      string `json:"responseFcm" db:"response_fcm"`
	Timestamp        int64  `json:"timestamp"`
	DeviceId         string `json:"deviceId"`
	CreatedBy        string `json:"createdBy"`
}
