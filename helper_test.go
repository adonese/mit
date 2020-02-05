package main

import (
	"reflect"
	"testing"

	"github.com/jinzhu/gorm"
)

func Test_getGrinderFromAgent(t *testing.T) {
	type args struct {
		db      *gorm.DB
		agentID int
	}
	tests := []struct {
		name string
		args args
		want Grinder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, got := getGrinderFromAgent(tt.args.db, tt.args.agentID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getGrinderFromAgent() = %v, want %v", got, tt.want)
			}
		})
	}
}
