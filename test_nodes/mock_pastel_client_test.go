package testnodes

import (
	"context"
	"testing"

	"github.com/pastelnetwork/gonode/pastel"
	"github.com/stretchr/testify/require"
)

func TestNewMockPastelClient(t *testing.T) {
	tests := []struct {
		name    string
		client  pastel.Client
		want    []byte
		wantErr bool
	}{
		{
			"success",
			NewMockPastelClient(nil),
			[]byte("mock signature"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.Sign(context.Background(), []byte("mock data"), "mock pastel id", "mock passphrase", "mock algorithm")
			if tt.wantErr {
				require.Error(t, err, "tt.client.Sign")
			} else {
				require.EqualValues(t, tt.want, got)
			}
			got1, err := tt.client.Verify(context.Background(), []byte("mock data"), "mock signature", "mock pastel id", "mock algorithm")
			if tt.wantErr {
				require.Error(t, err, "tt.client.Sign")
			} else {
				require.True(t, got1)
			}
		})
	}
}
