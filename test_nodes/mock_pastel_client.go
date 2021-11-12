package testnodes

import (
	"context"
	"log"

	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/storage-challenges/test_nodes/mocks"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func NewMockPastelClient(db *gorm.DB) pastel.Client {
	mc := &mocks.Client{}
	mc.On("Sign", mock.Anything, mock.MatchedBy(func(data []byte) bool {
		log.Println("MOCK PASTEL CLIENT -- Sign data param:", string(data))
		return true
	}), mock.MatchedBy(func(pastelID string) bool {
		log.Println("MOCK PASTEL CLIENT -- Sign pastel id param:", pastelID)
		return true
	}), mock.MatchedBy(func(passphrase string) bool {
		log.Println("MOCK PASTEL CLIENT -- Sign passphrase param:", passphrase)
		return true
	}), mock.MatchedBy(func(algorithm string) bool {
		log.Println("MOCK PASTEL CLIENT -- Sign algorithm param:", algorithm)
		return true
	})).Return([]byte("mock signature"), nil)

	mc.On("Verify", mock.Anything, mock.MatchedBy(func(data []byte) bool {
		log.Println("MOCK PASTEL CLIENT -- Verify data param:", string(data))
		return true
	}), mock.MatchedBy(func(signature string) bool {
		log.Println("MOCK PASTEL CLIENT -- Verify signature param:", signature)
		return true
	}), mock.MatchedBy(func(pastelID string) bool {
		log.Println("MOCK PASTEL CLIENT -- Verify pastel id param:", pastelID)
		return true
	}), mock.MatchedBy(func(algorithm string) bool {
		log.Println("MOCK PASTEL CLIENT -- Verify algorithm param:", algorithm)
		return true
	})).Return(true, nil)

	mc.On("GetBlockCount", mock.Anything).Run(func(args mock.Arguments) {
		log.Println("MOCK PASTEL CLIENT -- GetBlockHash")
	}).Return(func(ctx context.Context) int32 { return GetLastBlockNumer(db) }, func(ctx context.Context) error { return nil })

	mc.On("GetBlockHash", mock.Anything, mock.MatchedBy(func(blockHeight int32) bool {
		log.Println("MOCK PASTEL CLIENT -- GetBlockHash blockHeight param:", blockHeight)
		return true
	})).Return(func(ctx context.Context, blockHeight int32) string { return GetPastelBlockHash(blockHeight) }, func(ctx context.Context, blockHeight int32) error { return nil })

	mc.On("MasterNodesExtra", mock.Anything).Run(func(args mock.Arguments) {
		log.Println("MOCK PASTEL CLIENT -- MasterNodesExtra")
	}).Return(func(ctx context.Context) pastel.MasterNodes {
		var ret = pastel.MasterNodes{}
		for _, mn := range GetMasternodes(db) {
			ret = append(ret, pastel.MasterNode{
				ExtAddress: mn.MasternodeIPAddress,
				ExtKey:     mn.NodeID,
			})
		}
		return ret
	}, func(ctx context.Context) error { return nil })

	return mc
}
