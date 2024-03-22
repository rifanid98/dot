package deps

import (
	"dot/infrastructure/v1/apicall/xendit"
)

func (d *dependency) initApiCall() {
	d.initXenditApiCall()
}

func (d *dependency) initXenditApiCall() {
	d.apicall.XenditApiCall = xendit.New(d.base.Httpc, d.base.Cfg.ApiCall.Xendit)
}
