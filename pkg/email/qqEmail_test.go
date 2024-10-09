package qqEmail

import (
	"testing"

	"github.com/GoSimplicity/LinkMe/utils"
)

func TestQQEmail(t *testing.T) {
	type args struct {
		to      string
		subject string
		body    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test-1",
			args: args{
				to:      "yansaitao@qq.com",
				subject: "test-1",
				body:    "验证码为" + utils.GenRandomCode(6),
			},
			wantErr: false,
		},
		{
			name: "test-2",
			args: args{
				to:   "yansaitao@gmail.com",
				body: "验证码为" + utils.GenRandomCode(6),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendEmail(tt.args.to, tt.args.subject, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("sendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
