package core

import (
	"context"
	"encoding/json"
	"reflect"
)

type InternalContext struct {
	trackerId string
	data      map[string]any
}

func NewInternalContext(trackerID string) *InternalContext {
	return &InternalContext{
		trackerId: trackerID,
	}
}

func (ic *InternalContext) ToContext() context.Context {
	ctx := context.WithValue(context.Background(), "tracker_id", ic.trackerId)
	return ctx
}

func (ic *InternalContext) InjectData(data map[string]any) *InternalContext {
	ic.data = data
	return ic
}

func (ic *InternalContext) AppendData(data map[string]any) *InternalContext {
	if ic.data != nil {
		for k, v := range data {
			ic.data[k] = v
		}
		return ic
	}

	return ic.InjectData(data)
}

func (ic *InternalContext) GetData() map[string]any {
	return ic.data
}

func (ic *InternalContext) GetTrackerId() string {
	return ic.trackerId
}

func (ic *InternalContext) Clone() *InternalContext {
	return NewInternalContext(ic.trackerId)
}

func (ic *InternalContext) GetClientId() string {
	ctxData := ic.GetData()
	cid := ctxData["client_id"]
	if cid == nil {
		return ""
	}

	return dataToString(cid)
}

func (ic *InternalContext) GetAppId() string {
	ctxData := ic.GetData()
	aid := ctxData["app_id"]
	if aid == nil {
		return ""
	}

	return dataToString(aid)
}

func (ic *InternalContext) GetAccountId() string {
	ctxData := ic.GetData()
	aid := ctxData["id"]
	if aid == nil {
		return ""
	}

	return dataToString(aid)
}

func dataToString(data any) string {
	if reflect.TypeOf(data).String() == "string" {
		return data.(string)
	}
	bts, _ := json.Marshal(data)
	return string(bts)
}
