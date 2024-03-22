package repository

import (
	"context"
	"dot/pkg/util"

	"dot/core"
)

var log = util.NewLogger()

func ctx(ic *core.InternalContext) context.Context {
	ctxData := ic.GetData()
	if ctxData != nil {
		session := ctxData["session"]
		if session != nil {
			return session.(context.Context) // implements context.Context
		}
		return ic.ToContext()
	}
	return ic.ToContext()
}
