package domain

type PaymentRequest struct {
	Email       string                 `json:"email"`
	Amount      int64                  `json:"amount"`
	Currency    string                 `json:"currency,omitempty"`
	Reference   string                 `json:"reference,omitempty"`
	CallbackURL string                 `json:"callback_url,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type PaymentResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}
