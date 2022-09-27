package paginator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartIndexOption(t *testing.T) {
	type args struct {
		oss []OptionSet
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		{
			name: "",
			args: args{
				oss: []OptionSet{
					SetOrderOption(true),
					PerPageOption(15),
					DataKeyOption("data_key_test"),
					PerPageKeyOption("per_page_key_test"),
					PageKeyOption("page_key_test"),
					StartIndexOption(5),
				},
			},
			want: Option{
				startIndex:     5,
				startPage:      0,
				perPage:        15,
				perPageKey:     "per_page_key_test",
				pageKey:        "page_key_test",
				dataKey:        "data_key_test",
				firstPageKey:   "",
				lastPageKey:    "",
				nextPageKey:    "",
				prevPageKey:    "",
				currentPageKey: "",
				totalKey:       "",
				pathKey:        "",
				order:          true,
				lastKey:        "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var option Option
			for _, oss := range tt.args.oss {
				oss(&option)
			}
			assert.Equal(t, option.startPage, tt.want.startPage)
			assert.Equal(t, option.order, tt.want.order)
			assert.Equal(t, option.perPage, tt.want.perPage)
			assert.Equal(t, option.dataKey, tt.want.dataKey)
			assert.Equal(t, option.perPageKey, tt.want.perPageKey)
			assert.Equal(t, option.pageKey, tt.want.pageKey)
			assert.Equal(t, option.startIndex, tt.want.startIndex)
		})
	}
}
