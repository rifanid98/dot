package adapter

import (
	"dot/core"
)

//go:generate mockery --name XenditApiCall --filename xendit_apicall.go --output ./mocks
type XenditApiCall interface {
	QRCreate(ic *core.InternalContext, data map[string]any) (map[string]any, *core.CustomError)
	QRCheck(ic *core.InternalContext, data map[string]any) (map[string]any, *core.CustomError)
}
