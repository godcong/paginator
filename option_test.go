package paginator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartIndexOption(t *testing.T) {
	type args struct {
		oss []*Option
	}
	tests := []struct {
		name string
		args args
		want Config
	}{
		{
			name: "",
			args: args{
				oss: []*Option{
					SettingOption().
						SetPerPageKey("per_page_key_test").
						SetPageKey("page_key_test"),
					SettingPerPage(15).
						SetDataKey("data_key_test"),
					SettingStartIndex(5),
					SettingStartPage(20),
				},
			},
			want: Config{
				Keys: map[Key]string{
					KeyPerPage: "per_page_key_test",
					KeyPage:    "page_key_test",
					KeyData:    "data_key_test",
				},
				StartIndex: 5,
				StartPage:  20,
				PerPage:    15,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfig(tt.args.oss...)
			fmt.Println("config", config)
			assert.Equal(t, tt.want.GetKey(KeyPage), config.GetKey(KeyPage))
			assert.Equal(t, tt.want.GetKey(KeyPerPage), config.GetKey(KeyPerPage))
			assert.Equal(t, tt.want.GetKey(KeyData), config.GetKey(KeyData))
			assert.Equal(t, tt.want.PerPage, config.PerPage)
			assert.Equal(t, tt.want.StartIndex, config.StartIndex)
			assert.Equal(t, tt.want.StartPage, config.StartPage)
		})
	}
}
