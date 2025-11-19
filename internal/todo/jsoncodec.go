package todo

import "encoding/json"

type jsonCodec struct {
	name string
}

func (c jsonCodec) Name() string { return c.name }

func (jsonCodec) Marshal(v any) ([]byte, error) { return json.Marshal(v) }

func (jsonCodec) Unmarshal(data []byte, v any) error { return json.Unmarshal(data, v) }

func codecsForJSON() []jsonCodec {
	return []jsonCodec{
		{name: "json"},
		{name: "json; charset=utf-8"},
	}
}
