package models

type Merchant struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	AccountNo    string `json:"account_no"`
	IsRegistered bool   `json:"is_registered"`
}
