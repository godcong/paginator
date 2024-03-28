package paginator

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	fmt "github.com/k0kubun/pp/v3"
)

var testServer *httptest.Server
var testPaginator Paginator
var testPageReady PageReady

func init() {
	testPaginator = New(SettingOption().AddIgnore(KeyTotal).UsePathPaginate())
	testPageReady = testPaginator.Ready()
	testPageReady.CountHook(func(ctx context.Context, query Query) (int64, error) {
		return 1024, nil
	})
	testPageReady.ResultHook(func(ctx context.Context, query Query) (any, error) {
		var ret []string
		for i := 0; i < query.Limit(); i++ {
			ret = append(ret, "test"+strconv.Itoa(i+query.Offset()))
		}
		return ret, nil
	})
}

func handler(res http.ResponseWriter, req *http.Request) {
	q := FromHTTP(req)

	parse, err := testPageReady.DoParse(req.Context(), q, "page", "per_page")
	if err != nil {
		return
	}

	marshal, err := testPaginator.ToJSON(parse)
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(marshal)
}

func init() {
	// testPaginator = New()

	// http.HandleFunc("/test", handler)
	// http.ListenAndServe(":3030", nil)
	// testServer..ListenAndServe(":3000", testServer)
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
			req := httptest.NewRequest("GET", "http://127.0.0.1:18080/test?page=2&per_page=4", nil)
			// resp, err := http.Get("")

			w := httptest.NewRecorder()
			handler(w, req)
			resp := w.Result()
			// resp, err := http.DefaultClient.Do(req)
			// checkError(t, err)
			all, err := io.ReadAll(resp.Body)
			checkError(t, err)
			fmt.Println(string(all))
			var ret PageResult
			json.Unmarshal(all, &ret)
			fmt.Println(ret)
		})
	}
}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
