package paginator

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testServer *httptest.Server
var testPaginator Paginator

type query struct {
}

func (q query) Count(ctx context.Context) (int64, error) {
	return 999, nil
}

func (q query) Get(ctx context.Context, page PageQuery) (any, error) {
	pages := page.PerPage
	if has := int(page.Total) - (page.CurrentPage-1)*page.PerPage; has < pages {
		pages = has
	}

	var values []any
	for i := 0; i < pages; i++ {
		values = append(values, fmt.Sprintf("page.current.%d.%d", page.CurrentPage, i+1))
	}
	return values, nil
}

func (q query) Finder(p Parser) Finder {
	return q
}

func handler(res http.ResponseWriter, req *http.Request) {
	p := NewHTTPParser(req)
	parse, err := testPaginator.Parse(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(parse)
	marshal, err := json.Marshal(parse)
	res.WriteHeader(http.StatusOK)
	res.Write(marshal)
}

func init() {
	testPaginator = New().SetDefaultQuery(&query{})

	http.HandleFunc("/test", handler)
	//http.ListenAndServe(":3030", nil)
	//testServer..ListenAndServe(":3000", testServer)
}

func TestNew(t *testing.T) {
	type args struct {
		ops []string
	}
	tests := []struct {
		name string
		args args
		want Paginator
	}{
		{
			name: "",
			args: args{
				ops: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://127.0.0.1:18080/test?page=2&per_page=1", nil)
			//resp, err := http.Get("")

			w := httptest.NewRecorder()
			handler(w, req)
			resp := w.Result()
			//resp, err := http.DefaultClient.Do(req)
			//checkError(t, err)
			all, err := io.ReadAll(resp.Body)
			checkError(t, err)
			t.Logf("response: %v", string(all))
		})
	}
}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
