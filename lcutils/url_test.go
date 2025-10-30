package lcutils

import (
	"fmt"
	"testing"
)

func TestSanitizeURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				url: "0.2+0.3.png",
			},
			want: "0.2_0.3.png",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeFileName4URL(tt.args.url); got != tt.want {
				t.Errorf("SanitizeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessUnderscoresOptimized(t *testing.T) {
	// 测试用例
	testCases := []string{
		"_hello_world_",
		"__test___string__",
		"___multiple___underscores___",
		"no_underscores",
		"___",
		"",
		"_single_",
		"start_",
		"_end",
		"normal_string",
		"normal_string ",
		"normal_string  ",
		" normal_string",
		"  normal_string",
		"  normal_  str ing",
		"  nor   mal_  str ing",
		"  nor   mal_  \tstr ing",
		"  nor   mal_  \t\nstr ing",
		"  nor   mal_  \t\n\rstr ing",
		"a__b___c____d",
		"___开头_中间___结尾___",
	}

	fmt.Println("处理结果:")
	fmt.Println("=  * 50")

	for _, test := range testCases {
		result := SanitizeFileName4URL(test)
		fmt.Printf("输入: %-30s -> 输出: %-30s\n",
			fmt.Sprintf("\"%s\"", test),
			fmt.Sprintf("\"%s\"", result))
	}
}
