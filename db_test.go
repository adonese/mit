package main

import (
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"
)

func Test_getUser(t *testing.T) {
	type args struct {
		db       *gorm.DB
		username string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 User
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getUser(tt.args.db, tt.args.username)
			if got != tt.want {
				t.Errorf("getUser() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getUser() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
