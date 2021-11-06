package testnodes

import (
	"encoding/hex"
	"log"

	"github.com/pastelnetwork/gonode/pastel"
	"github.com/pastelnetwork/storage-challenges/test_nodes/mocks"
	"github.com/stretchr/testify/mock"
)

func NewMockPastelClient() pastel.Client {
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
	}).Return(int32(1), nil)

	mc.On("GetBlockHash", mock.Anything, mock.MatchedBy(func(blockHeight int32) bool {
		log.Println("MOCK PASTEL CLIENT -- GetBlockHash blockHeight param:", blockHeight)
		return true
	})).Return(hex.EncodeToString([]byte("mock block hash")), nil)

	mc.On("MasterNodesExtra", mock.Anything).Run(func(args mock.Arguments) {
		log.Println("MOCK PASTEL CLIENT -- MasterNodesExtra")
	}).Return(pastel.MasterNodes{pastel.MasterNode{
		ExtAddress: "localhost:9000",
		ExtKey:     "jXlzy0y3L1gYG04DBEZSKI9KV5BReiRzrW5bDBls3M2gtS6R0Ed8MHrEW9hzzgi4aW1taxNzChPSHEgJY4aTbw",
	}, {
		ExtAddress: "localhost:9001",
		ExtKey:     "jXEZVtIEVmSkYw0v8qGjsBrrELBOPuedNYMctelLWSlw6tiVNljFMpZFir30SN9r645tEAKwEAYfKR3o4Ek5YM",
	}, {
		ExtAddress: "localhost:9002",
		ExtKey:     "jXqBzHsk8P1cuRFrsRkQR5IhPzwFyCxE369KYqFLSITr8l5koLWcabZZDUVltIJ8666bE53G5fbtCz4veU2FCP",
	}, {
		ExtAddress: "localhost:9003",
		ExtKey:     "jXTwS1eCNDopMUIZAQnvpGlVe9lEnbauoh8TNDRoZcRTJVxCmZu1oSySBM1UwwyHDh7npbn01tZG0q2xyGmVJr",
	}}, nil)

	return mc
}
