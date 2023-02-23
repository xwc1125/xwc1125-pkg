package safecurl

import "testing"

func TestCheckIp(t *testing.T) {
	type args struct {
		ipStr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "ipv4Test", args: args{ipStr: "127.0.0.1"}, want: true},
		{name: "ipv6Test", args: args{ipStr: "2001:0db8:3c4d:0015:0000:0000:1a2f:1a2b"}, want: true},
		{name: "invalidIPTest", args: args{ipStr: "123456"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckIp(tt.args.ipStr); got != tt.want {
				t.Errorf("CheckIp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSafeURL(t *testing.T) {
	type args struct {
		inputUrl string
	}
	tests := []struct {
		name    string
		args    args
		wantR   string
		wantErr bool
	}{
		{name: "innerDomain", args: args{inputUrl: "http://tst.xwc1125.com"}, wantR: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := GetSafeURL(tt.args.inputUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSafeURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotR != tt.wantR {
				t.Errorf("GetSafeURL() gotR = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestIsInnerIp(t *testing.T) {
	type args struct {
		ipStr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "innerIP", args: args{ipStr: "192.168.1.109"}, want: true},
		{name: "outterIP", args: args{ipStr: "101.91.80.190"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInnerIp(tt.args.ipStr); got != tt.want {
				t.Errorf("IsInnerIp() = %v, want %v", got, tt.want)
			}
		})
	}
}
