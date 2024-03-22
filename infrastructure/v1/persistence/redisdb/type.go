package redisdb

import "encoding/json"

type Map map[string]any

func (m Map) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}
