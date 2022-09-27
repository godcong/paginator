package paginator

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type pageReady struct {
	parser Parser
	page   *PageQuery
	values func(page *PageQuery, op *Option) any
}

type KeyValue struct {
	Key   string
	Value any
}

type kvList []KeyValue

// MarshalJSON for map slice.
func (kv kvList) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.Write([]byte{'{'})
	for i, mi := range kv {
		b, err := json.Marshal(&mi.Value)
		if err != nil {
			return nil, err
		}
		buf.WriteString(fmt.Sprintf("%q:", fmt.Sprint(mi.Key)))
		buf.Write(b)
		if i < len(kv)-1 {
			buf.Write([]byte{','})
		}
	}
	buf.Write([]byte{'}'})
	return buf.Bytes(), nil
}

func unorderedValues(page *PageQuery, op *Option) any {
	values := make(map[string]any)
	values[op.TotalKey()] = page.Total
	values[op.PerPageKey()] = page.PerPage
	values[op.CurrentPageKey()] = page.CurrentPage
	values[op.LastKey()] = page.LastPage
	values[op.DataKey()] = page.Data
	values[op.PathKey()] = page.Path
	values[op.NextPageKey()] = page.NextPageURL
	values[op.PrevPageKey()] = page.PrevPageURL
	values[op.FirstPageKey()] = page.FirstPageURL
	values[op.LastPageKey()] = page.LastPageURL
	return values
}

func orderedValues(page *PageQuery, op *Option) any {
	list := kvList{
		KeyValue{op.TotalKey(), page.Total},
		KeyValue{op.PerPageKey(), page.PerPage},
		KeyValue{op.CurrentPageKey(), page.CurrentPage},
		KeyValue{op.LastKey(), page.LastPage},
		KeyValue{op.DataKey(), page.Data},
		KeyValue{op.PathKey(), page.Path},
		KeyValue{op.NextPageKey(), page.NextPageURL},
		KeyValue{op.PrevPageKey(), page.PrevPageURL},
		KeyValue{op.FirstPageKey(), page.FirstPageURL},
		KeyValue{op.LastPageKey(), page.LastPageURL},
	}
	return list
}
