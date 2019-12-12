package balance

import (
	"fmt"
	"testing"

	"github.com/alethio/web3-go/ethrpc"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-test/deep"
)

type MockETH struct {
	ethrpc.ETH
	balances   RawSheet
	throwError bool
}

func blockNumber(bn uint64) string {

}

func (m *MockETH) GetRawBalanceAtBlock(address, block string) (string, error) {
	if m.throwError == true {
		return "", fmt.Errorf("fatal error")
	}
	return m.balances[block][address][ETH], nil
}

func (m *MockETH) GetRawTokenBalanceAtBlock(address, token, block string) (string, error) {
	if m.throwError == true {
		return "", fmt.Errorf("fatal error")
	}
	return m.balances[block][address][Currency(token)], nil
}

func ExampleBookkeeper_GetIntBalanceResults() {
	r, err := ethrpc.NewWithDefaults("wss://mainnet.infura.io/ws")
	if err != nil {
		fmt.Println(err)
		return
	}

	b := New(r, 10)
	results, err := b.GetIntResults(balanceRequests())

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, res := range results {
		fmt.Printf("%s[%s]: %s\n", res.Request.Address, res.Request.Currency, res.Balance)
	}
	// Example output:
	// 0xa838e871a02c6d883bf004352fc7dac8f781fed6[0xd26114cd6EE289AccF82350c8d8487fedB8A0C07]: 409757565152676909
	// 0xa838e871a02c6d883bf004352fc7dac8f781fed6[0x8aa33a7899fcc8ea5fbe6a608a109c3893a1b8b2]: 3600000000000000000000
	// 0xa838e871a02c6d883bf004352fc7dac8f781fed6[0x0f5d2fb29fb7d3cfee444a200298f468908cc942]: 7041922408306145321820
	// 0xa838e871a02c6d883bf004352fc7dac8f781fed6[eth]: 1000670436501076869
	// 0xa838e871a02c6d883bf004352fc7dac8f781fed6[0xBEB9eF514a379B997e0798FDcC901Ee474B6D9A1]: 33780620000000000000
}

func TestGetBalancesWithOneAddressAndNoTokens(t *testing.T) {
	mockEth := &MockETH{}
	block := fmt.Sprintf("0x%d", 7500000)
	mockEth.balances = RawSheet{
		block: map[Address]map[Currency]string{
			"0x9fc201b6bc40cccbd5b588532ce98b845f95af51": map[Currency]string{
				ETH: fmt.Sprintf("0x%x", 100),
			},
		},
	}

	bookkeeper := New(mockEth, 10)
	requests := []*Request{
		{
			Address:    "0x9fc201b6bc40cccbd5b588532ce98b845f95af51",
			Currency:   ETH,
			DefaultBlockParam: block,
		},
	}

	balances, err := bookkeeper.GetRawSheet(requests)
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	spew.Dump(mockEth.balances)
	spew.Dump(balances)
	if diff := deep.Equal(balances, mockEth.balances); diff != nil {
		t.Error(diff)
	}
}

func TestGetBalancesWithOneAddressAndTokens(t *testing.T) {
	mockEth := &MockETH{}
	block := blockNumber(7500000)
	mockEth.balances = make(RawSheet)
	mockEth.balances = RawSheet{
		block: map[Address]map[Currency]string{
			"0x9fc201b6bc40cccbd5b588532ce98b845f95af51": map[Currency]string{
				ETH:     fmt.Sprintf("0x%x", 100),
				"0xabc": fmt.Sprintf("0x%x", 102),
			},
		},
	}

	bookkeeper := New(mockEth, 10)
	requests := []*Request{
		&Request{
			Address: "0x9fc201b6bc40cccbd5b588532ce98b845f95af51",
			Currency:  ETH,
			DefaultBlockParam: block,
		},
		&Request{
			Address: "0x9fc201b6bc40cccbd5b588532ce98b845f95af51",
			Currency:  "0xabc",
			DefaultBlockParam:   block,
		},
	}

	balances, err := bookkeeper.GetRawSheet(requests)
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	if diff := deep.Equal(balances, mockEth.balances); diff != nil {
		t.Error(diff)
	}
}

func TestGetBalancesWithMultipleAddressesAndTokens(t *testing.T) {
	mockEth := &MockETH{}
	block := blockNumber(7500000)
	mockEth.balances = make(RawSheet)
	mockEth.balances = RawSheet{
		block: map[Address]map[Currency]string{
			"0x9fc201b6bc40cccbd5b588532ce98b845f95af51": map[Currency]string{
				ETH:     fmt.Sprintf("0x%x", 100),
				"0xabc": fmt.Sprintf("0x%x", 102),
				"0xabd": fmt.Sprintf("0x%x", 105),
			},
			"0x9fc201b6bc40cccbd5b588532ce98b845f95af52": map[Currency]string{
				ETH:     fmt.Sprintf("0x%x", 101),
				"0xabc": fmt.Sprintf("0x%x", 103),
				"0xabd": fmt.Sprintf("0x%x", 104),
			},
		},
	}

	bookkeeper := New(mockEth, 10)
	requests := make([]*Request, 0, 0)
	for _, address := range []Address{"0x9fc201b6bc40cccbd5b588532ce98b845f95af51", "0x9fc201b6bc40cccbd5b588532ce98b845f95af52"} {
		for _, Currency := range []Currency{ETH, "0xabc", "0xabd"} {
			requests = append(requests, &Request{
				Address: address,
				Currency:  Currency,
				DefaultBlockParam:  block,
			})
		}
	}

	balances, err := bookkeeper.GetRawSheet(requests)
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	if diff := deep.Equal(balances, mockEth.balances); diff != nil {
		t.Error(diff)
	}
}

func TestGetBalancesWithError(t *testing.T) {
	mockEth := &MockETH{throwError: true}
	bookkeeper := New(mockEth, 10)
	block := blockNumber(7500000)
	requests := []*Request{
		&Request{
			Address: "0x9fc201b6bc40cccbd5b588532ce98b845f95af51",
			Currency:  ETH,
			DefaultBlockParam:   block,
		},
		&Request{
			Address: "0x9fc201b6bc40cccbd5b588532ce98b845f95af51",
			Currency:  "0xabc",
			DefaultBlockParam:   block,
		},
	}

	_, err := bookkeeper.GetRawResults(requests)
	if err == nil {
		t.Fatal("Expecting error")
	}
}

func balanceRequests() []*Request {
	address := Address("0xa838e871a02c6d883bf004352fc7dac8f781fed6")
	block := blockNumber(7500000)
	return []*Request{
		&Request{
			Address: address,
			DefaultBlockParam:   block,
			Currency:  ETH,
		},
		&Request{
			Address: address,
			DefaultBlockParam:   block,
			Currency:  Currency("0xBEB9eF514a379B997e0798FDcC901Ee474B6D9A1"),
		},
		&Request{
			Address: address,
			DefaultBlockParam:   block,
			Currency:  Currency("0x0f5d2fb29fb7d3cfee444a200298f468908cc942"),
		},
		&Request{
			Address: address,
			DefaultBlockParam:   block,
			Currency:  Currency("0xd26114cd6EE289AccF82350c8d8487fedB8A0C07"),
		},
		&Request{
			Address: address,
			DefaultBlockParam:   block,
			Currency:  Currency("0x8aa33a7899fcc8ea5fbe6a608a109c3893a1b8b2"),
		},
	}
}
