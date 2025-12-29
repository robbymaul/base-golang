package web

type ListTransactionRequest struct {
	Currency  string `json:"currency"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}
