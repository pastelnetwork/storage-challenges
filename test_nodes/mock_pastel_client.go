package testnodes

import (
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
	return mc
}
