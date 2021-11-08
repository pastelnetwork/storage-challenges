// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	pastel "github.com/pastelnetwork/gonode/pastel"
	mock "github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// ActTickets provides a mock function with given fields: ctx, actType, minHeight
func (_m *Client) ActTickets(ctx context.Context, actType pastel.ActTicketType, minHeight int) (pastel.ActTickets, error) {
	ret := _m.Called(ctx, actType, minHeight)

	var r0 pastel.ActTickets
	if rf, ok := ret.Get(0).(func(context.Context, pastel.ActTicketType, int) pastel.ActTickets); ok {
		r0 = rf(ctx, actType, minHeight)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pastel.ActTickets)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pastel.ActTicketType, int) error); ok {
		r1 = rf(ctx, actType, minHeight)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindTicketByID provides a mock function with given fields: ctx, pastelid
func (_m *Client) FindTicketByID(ctx context.Context, pastelid string) (*pastel.IDTicket, error) {
	ret := _m.Called(ctx, pastelid)

	var r0 *pastel.IDTicket
	if rf, ok := ret.Get(0).(func(context.Context, string) *pastel.IDTicket); ok {
		r0 = rf(ctx, pastelid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pastel.IDTicket)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, pastelid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBalance provides a mock function with given fields: ctx, address
func (_m *Client) GetBalance(ctx context.Context, address string) (float64, error) {
	ret := _m.Called(ctx, address)

	var r0 float64
	if rf, ok := ret.Get(0).(func(context.Context, string) float64); ok {
		r0 = rf(ctx, address)
	} else {
		r0 = ret.Get(0).(float64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockCount provides a mock function with given fields: ctx
func (_m *Client) GetBlockCount(ctx context.Context) (int32, error) {
	ret := _m.Called(ctx)

	var r0 int32
	if rf, ok := ret.Get(0).(func(context.Context) int32); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int32)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockHash provides a mock function with given fields: ctx, blkIndex
func (_m *Client) GetBlockHash(ctx context.Context, blkIndex int32) (string, error) {
	ret := _m.Called(ctx, blkIndex)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, int32) string); ok {
		r0 = rf(ctx, blkIndex)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int32) error); ok {
		r1 = rf(ctx, blkIndex)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBlockVerbose1 provides a mock function with given fields: ctx, blkHeight
func (_m *Client) GetBlockVerbose1(ctx context.Context, blkHeight int32) (*pastel.GetBlockVerbose1Result, error) {
	ret := _m.Called(ctx, blkHeight)

	var r0 *pastel.GetBlockVerbose1Result
	if rf, ok := ret.Get(0).(func(context.Context, int32) *pastel.GetBlockVerbose1Result); ok {
		r0 = rf(ctx, blkHeight)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pastel.GetBlockVerbose1Result)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int32) error); ok {
		r1 = rf(ctx, blkHeight)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetInfo provides a mock function with given fields: ctx
func (_m *Client) GetInfo(ctx context.Context) (*pastel.GetInfoResult, error) {
	ret := _m.Called(ctx)

	var r0 *pastel.GetInfoResult
	if rf, ok := ret.Get(0).(func(context.Context) *pastel.GetInfoResult); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pastel.GetInfoResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNFTTicketFeePerKB provides a mock function with given fields: ctx
func (_m *Client) GetNFTTicketFeePerKB(ctx context.Context) (int64, error) {
	ret := _m.Called(ctx)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context) int64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNetworkFeePerMB provides a mock function with given fields: ctx
func (_m *Client) GetNetworkFeePerMB(ctx context.Context) (int64, error) {
	ret := _m.Called(ctx)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context) int64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRawTransactionVerbose1 provides a mock function with given fields: ctx, txID
func (_m *Client) GetRawTransactionVerbose1(ctx context.Context, txID string) (*pastel.GetRawTransactionVerbose1Result, error) {
	ret := _m.Called(ctx, txID)

	var r0 *pastel.GetRawTransactionVerbose1Result
	if rf, ok := ret.Get(0).(func(context.Context, string) *pastel.GetRawTransactionVerbose1Result); ok {
		r0 = rf(ctx, txID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pastel.GetRawTransactionVerbose1Result)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, txID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRegisterExDDFee provides a mock function with given fields: ctx, request
func (_m *Client) GetRegisterExDDFee(ctx context.Context, request pastel.GetRegisterExDDFeeRequest) (int64, error) {
	ret := _m.Called(ctx, request)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, pastel.GetRegisterExDDFeeRequest) int64); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pastel.GetRegisterExDDFeeRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRegisterExternalStorageFee provides a mock function with given fields: ctx, request
func (_m *Client) GetRegisterExternalStorageFee(ctx context.Context, request pastel.GetRegisterExternalStorageFeeRequest) (int64, error) {
	ret := _m.Called(ctx, request)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, pastel.GetRegisterExternalStorageFeeRequest) int64); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pastel.GetRegisterExternalStorageFeeRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRegisterNFTFee provides a mock function with given fields: ctx, request
func (_m *Client) GetRegisterNFTFee(ctx context.Context, request pastel.GetRegisterNFTFeeRequest) (int64, error) {
	ret := _m.Called(ctx, request)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, pastel.GetRegisterNFTFeeRequest) int64); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pastel.GetRegisterNFTFeeRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransaction provides a mock function with given fields: ctx, txID
func (_m *Client) GetTransaction(ctx context.Context, txID string) (*pastel.GetTransactionResult, error) {
	ret := _m.Called(ctx, txID)

	var r0 *pastel.GetTransactionResult
	if rf, ok := ret.Get(0).(func(context.Context, string) *pastel.GetTransactionResult); ok {
		r0 = rf(ctx, txID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pastel.GetTransactionResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, txID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IDTickets provides a mock function with given fields: ctx, idType
func (_m *Client) IDTickets(ctx context.Context, idType pastel.IDTicketType) (pastel.IDTickets, error) {
	ret := _m.Called(ctx, idType)

	var r0 pastel.IDTickets
	if rf, ok := ret.Get(0).(func(context.Context, pastel.IDTicketType) pastel.IDTickets); ok {
		r0 = rf(ctx, idType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pastel.IDTickets)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pastel.IDTicketType) error); ok {
		r1 = rf(ctx, idType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAvailableTradeTickets provides a mock function with given fields: ctx
func (_m *Client) ListAvailableTradeTickets(ctx context.Context) ([]pastel.TradeTicket, error) {
	ret := _m.Called(ctx)

	var r0 []pastel.TradeTicket
	if rf, ok := ret.Get(0).(func(context.Context) []pastel.TradeTicket); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]pastel.TradeTicket)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterNodeConfig provides a mock function with given fields: ctx
func (_m *Client) MasterNodeConfig(ctx context.Context) (*pastel.MasterNodeConfig, error) {
	ret := _m.Called(ctx)

	var r0 *pastel.MasterNodeConfig
	if rf, ok := ret.Get(0).(func(context.Context) *pastel.MasterNodeConfig); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pastel.MasterNodeConfig)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterNodeStatus provides a mock function with given fields: ctx
func (_m *Client) MasterNodeStatus(ctx context.Context) (*pastel.MasterNodeStatus, error) {
	ret := _m.Called(ctx)

	var r0 *pastel.MasterNodeStatus
	if rf, ok := ret.Get(0).(func(context.Context) *pastel.MasterNodeStatus); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pastel.MasterNodeStatus)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterNodesExtra provides a mock function with given fields: ctx
func (_m *Client) MasterNodesExtra(ctx context.Context) (pastel.MasterNodes, error) {
	ret := _m.Called(ctx)

	var r0 pastel.MasterNodes
	if rf, ok := ret.Get(0).(func(context.Context) pastel.MasterNodes); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pastel.MasterNodes)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterNodesList provides a mock function with given fields: ctx
func (_m *Client) MasterNodesList(ctx context.Context) (pastel.MasterNodes, error) {
	ret := _m.Called(ctx)

	var r0 pastel.MasterNodes
	if rf, ok := ret.Get(0).(func(context.Context) pastel.MasterNodes); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pastel.MasterNodes)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterNodesTop provides a mock function with given fields: ctx
func (_m *Client) MasterNodesTop(ctx context.Context) (pastel.MasterNodes, error) {
	ret := _m.Called(ctx)

	var r0 pastel.MasterNodes
	if rf, ok := ret.Get(0).(func(context.Context) pastel.MasterNodes); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pastel.MasterNodes)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegTicket provides a mock function with given fields: ctx, regTxid
func (_m *Client) RegTicket(ctx context.Context, regTxid string) (pastel.RegTicket, error) {
	ret := _m.Called(ctx, regTxid)

	var r0 pastel.RegTicket
	if rf, ok := ret.Get(0).(func(context.Context, string) pastel.RegTicket); ok {
		r0 = rf(ctx, regTxid)
	} else {
		r0 = ret.Get(0).(pastel.RegTicket)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, regTxid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegTickets provides a mock function with given fields: ctx
func (_m *Client) RegTickets(ctx context.Context) (pastel.RegTickets, error) {
	ret := _m.Called(ctx)

	var r0 pastel.RegTickets
	if rf, ok := ret.Get(0).(func(context.Context) pastel.RegTickets); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pastel.RegTickets)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterActTicket provides a mock function with given fields: ctx, regTicketTxid, artistHeight, fee, pastelID, passphrase
func (_m *Client) RegisterActTicket(ctx context.Context, regTicketTxid string, artistHeight int, fee int64, pastelID string, passphrase string) (string, error) {
	ret := _m.Called(ctx, regTicketTxid, artistHeight, fee, pastelID, passphrase)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int64, string, string) string); ok {
		r0 = rf(ctx, regTicketTxid, artistHeight, fee, pastelID, passphrase)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, int, int64, string, string) error); ok {
		r1 = rf(ctx, regTicketTxid, artistHeight, fee, pastelID, passphrase)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterExDDTicket provides a mock function with given fields: ctx, request
func (_m *Client) RegisterExDDTicket(ctx context.Context, request pastel.RegisterExDDRequest) (string, error) {
	ret := _m.Called(ctx, request)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, pastel.RegisterExDDRequest) string); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pastel.RegisterExDDRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterExternalStorageTicket provides a mock function with given fields: ctx, request
func (_m *Client) RegisterExternalStorageTicket(ctx context.Context, request pastel.RegisterExternalStorageRequest) (string, error) {
	ret := _m.Called(ctx, request)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, pastel.RegisterExternalStorageRequest) string); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pastel.RegisterExternalStorageRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterNFTTicket provides a mock function with given fields: ctx, request
func (_m *Client) RegisterNFTTicket(ctx context.Context, request pastel.RegisterNFTRequest) (string, error) {
	ret := _m.Called(ctx, request)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, pastel.RegisterNFTRequest) string); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, pastel.RegisterNFTRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendFromAddress provides a mock function with given fields: ctx, fromID, toID, amount
func (_m *Client) SendFromAddress(ctx context.Context, fromID string, toID string, amount float64) (string, error) {
	ret := _m.Called(ctx, fromID, toID, amount)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string, float64) string); ok {
		r0 = rf(ctx, fromID, toID, amount)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, float64) error); ok {
		r1 = rf(ctx, fromID, toID, amount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendToAddress provides a mock function with given fields: ctx, pastelID, amount
func (_m *Client) SendToAddress(ctx context.Context, pastelID string, amount int64) (string, error) {
	ret := _m.Called(ctx, pastelID, amount)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, int64) string); ok {
		r0 = rf(ctx, pastelID, amount)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, int64) error); ok {
		r1 = rf(ctx, pastelID, amount)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Sign provides a mock function with given fields: ctx, data, pastelID, passphrase, algorithm
func (_m *Client) Sign(ctx context.Context, data []byte, pastelID string, passphrase string, algorithm string) ([]byte, error) {
	ret := _m.Called(ctx, data, pastelID, passphrase, algorithm)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, []byte, string, string, string) []byte); ok {
		r0 = rf(ctx, data, pastelID, passphrase, algorithm)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []byte, string, string, string) error); ok {
		r1 = rf(ctx, data, pastelID, passphrase, algorithm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StorageNetworkFee provides a mock function with given fields: ctx
func (_m *Client) StorageNetworkFee(ctx context.Context) (float64, error) {
	ret := _m.Called(ctx)

	var r0 float64
	if rf, ok := ret.Get(0).(func(context.Context) float64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(float64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TicketOwnership provides a mock function with given fields: ctx, txID, pastelID, passphrase
func (_m *Client) TicketOwnership(ctx context.Context, txID string, pastelID string, passphrase string) (string, error) {
	ret := _m.Called(ctx, txID, pastelID, passphrase)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) string); ok {
		r0 = rf(ctx, txID, pastelID, passphrase)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, txID, pastelID, passphrase)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Verify provides a mock function with given fields: ctx, data, signature, pastelID, algorithm
func (_m *Client) Verify(ctx context.Context, data []byte, signature string, pastelID string, algorithm string) (bool, error) {
	ret := _m.Called(ctx, data, signature, pastelID, algorithm)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, []byte, string, string, string) bool); ok {
		r0 = rf(ctx, data, signature, pastelID, algorithm)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []byte, string, string, string) error); ok {
		r1 = rf(ctx, data, signature, pastelID, algorithm)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
