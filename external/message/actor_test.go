package message

import (
	"testing"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/pastelnetwork/storage-challenges/application/dto"
)

func TestNewActorPIDSet(t *testing.T) {
	type args struct {
		addresses map[string][]string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			"success",
			args{
				map[string][]string{"id": {"127.0.0.1:8080", "127.0.0.1:8001"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(NewActorPIDSet(tt.args.addresses))
		})
	}
}

func TestDo(t *testing.T) {
	type args struct {
		clients  *actor.PIDSet
		message  interface{}
		callback func(context actor.Context, message interface{})
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			"success",
			args{
				clients: actor.NewPIDSet(actor.NewPID("127.0.0.1:8000", "storage-challenge")),
				message: &dto.StorageChallengeRequest{},
				callback: func(context actor.Context, message interface{}) {
					t.Log("CALLBACK", message)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Do(tt.args.clients, tt.args.message, tt.args.callback)
			c()
		})
	}
}
