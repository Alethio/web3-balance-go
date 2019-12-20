package balance

import (
	"github.com/alethio/web3-go/ethrpc"
	"sync"
)

type rpcLoader struct {
	eth ethrpc.ETHInterface
}

func (loader rpcLoader) fetchRequests(b *Bookkeeper, requests []*Request, results chan *RawResponse, done chan error) {
	var tries uint = 0
	wg := sync.WaitGroup{}

	for {
		failed := make(chan *RequestError, len(requests))
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

		if len(reqErrors) == 0 {
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
