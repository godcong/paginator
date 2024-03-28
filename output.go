package paginator

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type KeyValue struct {
	Key   string
	Value any
}

type CustomOutput struct {
	*Config
	KeyValues []KeyValue
}

func (o *CustomOutput) Add(key Key, value any) {
	if o.IsIgnored(key) {
		return
	}
	o.KeyValues = append(o.KeyValues, KeyValue{o.GetKey(key), value})
}

// MarshalJSON for map slice.
func (o *CustomOutput) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.Write([]byte{'{'})
	for i, mi := range o.KeyValues {
		b, err := json.Marshal(&mi.Value)
		if err != nil {
			return nil, err
		}
		buf.WriteString(fmt.Sprintf("%q:", fmt.Sprint(mi.Key)))
		buf.Write(b)
		if i < len(o.KeyValues)-1 {
			buf.Write([]byte{','})
		}
	}
	buf.Write([]byte{'}'})
	return buf.Bytes(), nil
}

func customOutput(cfg *Config, result *PageResult) any {
	output := CustomOutput{
		Config: cfg,
	}

	output.Add(KeySuccess, result.Success)
	output.Add(KeyTotal, result.Total)
	output.Add(KeyPerPage, result.PerPage)
	output.Add(KeyCurrentPage, result.CurrentPage)
	output.Add(KeyLastPage, result.LastPage)
	output.Add(KeyPreviousPage, result.PreviousPage)
	output.Add(KeyNextPage, result.NextPage)
	output.Add(KeyPath, result.Path)
	output.Add(KeyData, result.Data)

	if cfg.SkipPath {
		return output
	}
	output.Add(KeyCurrentPath, result.CurrentPath)
	output.Add(KeyNextPageUrl, result.NextPageURL)
	output.Add(KeyPrevPageUrl, result.PrevPageURL)
	output.Add(KeyFirstPageUrl, result.FirstPageURL)
	output.Add(KeyLastPageUrl, result.LastPageURL)
	return &output
}
