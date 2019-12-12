package balance

import (
	"math/big"
	"sync"

	"github.com/alethio/web3-go/ethrpc"
	"github.com/alethio/web3-go/strhelper"
)

// New returns a new Bookkeeper struct
func New(eth ethrpc.ETHInterface, opts ...Option) *Bookkeeper {
	config := &Config{
		Retry: false,
		Attempts: 0,
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

	go b.fetchRequests(requests, results, done)

	for {
		select {
		case result := <-results:
			responses = append(responses, result)
		case err := <-done:
			return responses, err
		}
	}

}

func (b *Bookkeeper) fetchRequests(requests []*Request, results chan *RawResponse, done chan error) {
	var tries uint = 0
	wg := sync.WaitGroup{}

	for {
		failed := make(chan *RequestError, len(requests))
		errors := make(chan error, len(requests))
		for _, request := range requests {
			wg.Add(1)
			go func(req *Request, results chan *RawResponse, failed chan *RequestError) {
				defer wg.Done()
				var balance string
				var err error

				address := string(req.Address)
				if req.Currency == ETH {
					balance, err = b.eth.GetRawBalanceAtBlock(address, req.DefaultBlockParam)
				} else {
					token := string(req.Currency)
					balance, err = b.eth.GetRawTokenBalanceAtBlock(address, token, req.DefaultBlockParam)
				}

				if err != nil {
					failed <- &RequestError{req, err}
				} else {
					results <- &RawResponse{
						Request: req,
						Balance: balance,
					}
					errors <- err
				}
			}(request, results, failed)
		}

		wg.Wait()
		close(failed)

		requests = make([]*Request, 0, len(requests))
		reqErrors := make([]*RequestError, 0, len(requests))

		for reqError := range failed {
			reqErrors = append(reqErrors, reqError)
			requests = append(requests, reqError.Request)
		}

		if len(requests) == 0 {
			done <- nil
			return
		}

		tries++
		if b.config.Retry == false || tries > b.config.Attempts {
			done <- CollectError{reqErrors}
			return
		}
	}
}
