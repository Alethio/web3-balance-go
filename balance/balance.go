package balance

import (
	"math/big"

	"github.com/alethio/web3-go/ethrpc"
	"github.com/alethio/web3-go/strhelper"
)

// New returns a new Bookkeeper struct
func New(eth ethrpc.ETHInterface, opts ...Option) *Bookkeeper {
	config := &Config{
		Retry: false,
		Attempts: 1,
		Loader: rpcLoader{},
	}

	for _, opt := range opts {
		opt(config)
	}

	return &Bookkeeper{
		eth:     eth,
		config: config,
	}
}

// GetIntBalanceSheet takes a list of balance requests and returns a tree like
// structure containing all int balances
func (b *Bookkeeper) GetIntSheet(requests []*Request) (IntSheet, error) {
	balances := make(IntSheet)
	intResponses, err := b.GetIntResults(requests)
	for _, result := range intResponses {
		block := result.Request.DefaultBlockParam
		address := result.Request.Address
		source := result.Request.Currency

		if balances[block] == nil {
			balances[block] = make(map[Address]map[Currency]*big.Int)
		}

		if balances[block][address] == nil {
			balances[block][address] = make(map[Currency]*big.Int)

		}

		balances[block][address][source] = result.Balance
	}
	return balances, err
}

// GetRawBalanceSheet takes a list of balance requests and returns a tree like
// structure containing all hex string balances
func (b *Bookkeeper) GetRawSheet(requests []*Request) (RawSheet, error) {
	balances := make(RawSheet)
	rawResponses, err := b.GetRawResults(requests)
	for _, result := range rawResponses {
		block := result.Request.DefaultBlockParam
		address := result.Request.Address
		source := result.Request.Currency

		if balances[block] == nil {
			balances[block] = make(map[Address]map[Currency]string)
		}

		if balances[block][address] == nil {
			balances[block][address] = make(map[Currency]string)

		}

		balances[block][address][source] = result.Balance
	}
	return balances, err
}

// GetIntBalanceResults returns an array of *big.Int balance results for the provided requests
func (b *Bookkeeper) GetIntResults(requests []*Request) ([]*IntResponse, error) {
	intResponses := make([]*IntResponse, 0, len(requests))
	failedRequests := make([]*RequestError, 0, len(requests))

	rawResponses, err := b.GetRawResults(requests)
	if err != nil {
		if collectError, ok := err.(CollectError); ok {
			failedRequests = append(failedRequests, collectError.Errors...)
		} else {
			return intResponses, err
		}
	}

	for _, rawResponse := range rawResponses {
		intBalance, err := strhelper.HexStrToBigInt(rawResponse.Balance)
		if err != nil {
			failedRequests = append(failedRequests, &RequestError{rawResponse.Request, err})
		} else {
			intResponses = append(intResponses, &IntResponse{
				Request: rawResponse.Request,
				Balance: intBalance,
			})
		}
	}

	if len(failedRequests) > 0 {
		return intResponses, CollectError{failedRequests}
	}
	return intResponses, nil
}

// GetRawBalanceResults returns an array of hex string balance results for the provided requests
func (b *Bookkeeper) GetRawResults(requests []*Request) ([]*RawResponse, error) {
	results := make(chan *RawResponse)
	responses := make([]*RawResponse, 0, len(requests))

	done := make(chan error, 1)

	go b.config.Loader.fetchRequests(b, requests, results, done)

	for {
		select {
		case result := <-results:
			responses = append(responses, result)
		case err := <-done:
			return responses, err
		}
	}

}

