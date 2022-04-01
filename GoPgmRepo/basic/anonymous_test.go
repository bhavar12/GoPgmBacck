package main

import "testing"

func Test_another(t *testing.T) {
	type args struct {
		f func(string) string
	}
	anon := func(str string) string {
		return str
	}
	tests := []struct {
		args args
		want string
	}{
		{anon, "hello"},
		{anon, "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if str := another(tt.args.f); str != tt.want {
				t.Error("Test Failed: {} inputted, {} expected, recieved: {}", tt.args, tt.want, str)
			}
		})
	}
}
