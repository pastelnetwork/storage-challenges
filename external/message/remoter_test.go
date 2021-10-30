package message

import (
	"testing"
)

func TestNewActorPIDSet(t *testing.T) {
	type args struct {
		properties []ActorProperties
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			"success",
			args{
				properties: []ActorProperties{
					{Address: "127.0.0.1:8080", Kind: "id"},
					{Address: "127.0.0.1:8001", Kind: "id"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(NewActorPIDSet(tt.args.properties))
		})
	}
}
