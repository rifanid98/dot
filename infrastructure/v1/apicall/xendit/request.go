package xendit

import (
	"dot/pkg/helper"
	"time"
)

type QrCodesCreate struct {
	ReferenceId string    `json:"reference_id"`
	Type        string    `json:"type"`
	Currency    string    `json:"currency"`
	Amount      int64     `json:"amount"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func (q *QrCodesCreate) Bind(data map[string]any) *QrCodesCreate {
	return &QrCodesCreate{
		ReferenceId: helper.DataToString(data["account_id"]),
		Type:        "DYNAMIC",
		Currency:    "IDR",
		Amount:      helper.DataToInt(data["amount"]),
		ExpiresAt:   time.Now().Add(time.Hour * 24),
	}
}
