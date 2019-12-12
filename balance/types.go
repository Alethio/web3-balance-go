package balance

import (
	"math/big"

	"github.com/alethio/web3-go/ethrpc"
)

// Bookkeeper wraps the operations
type Bookkeeper struct {
	eth     ethrpc.ETHInterface
	config  *Config
}

// BlockNumber : wrapper type for a block number
type BlockNumber = uint64

// Address : wrapper type for an ETH address
type Address = string

// RawBalanceSheet : the tree-like structure representing all un-parsed balances
type RawSheet = map[string]map[Address]map[Currency]string

// IntBalanceSheet : the tree-like structure with all balances converted to big.Int
type IntSheet = map[string]map[Address]map[Currency]*big.Int

// Source : either "ETH" or a token address
type Currency string

const (
	// ETH : Ethereum
	ETH Currency = "ETH"
)

// BalanceRequest : a unit of work
type Request struct {
	DefaultBlockParam string
	Address Address
	Currency Currency
}

// RawBalanceResponse : the raw response associated with a balance request
type RawResponse struct {
	Request *Request
	Balance string
}

// IntBalanceResponse : the converted response associated with a balance request
type IntResponse struct {
	Request *Request
	Balance *big.Int
}
