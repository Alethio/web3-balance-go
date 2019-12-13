package balance

import (
	"encoding/hex"
	"fmt"
	"github.com/alethio/web3-multicall-go/multicall"
	"github.com/avast/retry-go"
	"golang.org/x/sync/errgroup"
)

type multicallLoader struct {
	mc *multicall.Multicall
}

func (loader multicallLoader) fetchRequests(b *Bookkeeper, requests []*Request, results chan *RawResponse, done chan error) {
	viewCallsByBlock := make(map[string]multicall.ViewCalls)
	reqMap := make(map[string]*Request)

	for index, req := range requests {
		key := fmt.Sprintf("%s-%d", req.DefaultBlockParam, index)
		reqMap[key] = req
		viewCalls, ok := viewCallsByBlock[req.DefaultBlockParam]
		if !ok {
			viewCalls = make(multicall.ViewCalls, 0, 0)
		}

		var viewCall multicall.ViewCall
		if req.Currency == ETH {
			viewCall = multicall.ViewCall{
				Key: key,
				Target: loader.mc.Contract(),
				Method: "getEthBalance(address)(uint256)",
				Arguments: []interface{}{req.Address},
			}
		} else {
			viewCall = multicall.ViewCall{
				Key: key,
				Target: string(req.Currency),
				Method: "balanceOf(address)(uint256)",
				Arguments: []interface{}{req.Address},
			}
		}

		viewCalls = append(viewCalls, viewCall)
		viewCallsByBlock[req.DefaultBlockParam] = viewCalls
	}


	group := errgroup.Group{}
	failed := make(chan *RequestError, len(requests))

	for defaultBlockParam, viewCalls := range viewCallsByBlock {
		defaultBlockParam := defaultBlockParam
		viewCalls := viewCalls

		group.Go(func() error {
			err := retry.Do(
				func() error {
					res, err := loader.mc.CallRaw(viewCalls, defaultBlockParam)
					if err != nil {
						fmt.Println(err)
						return err
					}
					for key, result := range res.Calls {
						if result.Success {
							balance := result.ReturnValues[0].([]byte)
							hexBalance := hex.EncodeToString(balance)
							results <- &RawResponse{
								Request: reqMap[key],
								Balance: hexBalance,
							}
						} else {
							failed <- &RequestError{reqMap[key], fmt.Errorf("VM Error")}
						}
					}
					return nil
				},
				retry.Attempts(b.config.Attempts),
			)
			return err
		})
	}

	err := group.Wait()
	close(failed)

	if err != nil {
		reqErrors := make([]*RequestError, 0, len(requests))
		for reqError := range failed {
			reqErrors = append(reqErrors, reqError)
		}

		if len(reqErrors) == 0 {
			done <- err
		} else {
			done <- CollectError{reqErrors}
		}
	} else {
		done <- nil
	}
}
