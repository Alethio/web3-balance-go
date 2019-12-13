package main

import (
	"flag"
	"fmt"
	"github.com/alethio/web3-balance-go/balance"
	"github.com/alethio/web3-go/ethrpc"
	"github.com/alethio/web3-multicall-go/multicall"
	"github.com/davecgh/go-spew/spew"
	"log"
	"os"
)

var ethURL string

type worker struct {
	eth ethrpc.ETHInterface
}

func main() {
	flag.StringVar(&ethURL, "eth-client-url", "http://localhost:8546", "URL of an Ethereum Client (parity needed)")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		log.Println("Please issue a command:")
		log.Println("  exampleBalances")
		os.Exit(0)
	}

	eth, err := ethrpc.NewWithDefaults(ethURL)
	if err != nil {
		log.Fatal(err)
	}
	b := balance.New(eth)

	cmd := args[0]
	switch cmd {
	case "exampleMulticallBalances":
		// block := fmt.Sprintf("0x%x", 7000000)
		block := "latest"
		mc, err := multicall.New(eth, multicall.ContractAddress(multicall.MainnetAddress))
		if err != nil {
			panic(err)
		}
		b := balance.New(eth, balance.UseMulticall(mc))
		requests := []*balance.Request{
			&balance.Request{
				Address: "0xa838e871a02c6d883bf004352fc7dac8f781fed6",
				Currency:  balance.ETH,
				DefaultBlockParam:   block,
			},
			&balance.Request{
				Address: "0xa838e871a02c6d883bf004352fc7dac8f781fed6",
				Currency:  "0xBEB9eF514a379B997e0798FDcC901Ee474B6D9A1",
				DefaultBlockParam:   block,
			},
			&balance.Request{
				Address: "0xa838e871a02c6d883bf004352fc7dac8f781fed6",
				Currency:  "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07",
				DefaultBlockParam:   block,
			},
			&balance.Request{
				Address: "0xa838e871a02c6d883bf004352fc7dac8f781fed6",
				Currency:  "0x8aa33a7899fcc8ea5fbe6a608a109c3893a1b8b2",
				DefaultBlockParam:   block,
			},
			&balance.Request{
				Address: "0xa838e871a02c6d883bf004352fc7dac8f781fed6",
				Currency: "0x0f5d2fb29fb7d3cfee444a200298f468908cc942",
				DefaultBlockParam:  block,
			},
		}
		rawBalances, err := b.GetRawSheet(requests)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("--------------- Raw Balances -----------------")
		spew.Dump(rawBalances)
		intBalances, err := b.GetIntSheet(requests)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("--------------- big.Int Balances -----------------")
		spew.Dump(intBalances)
	case "exampleBalances":
		// block := fmt.Sprintf("0x%x", 7000000)
		block := "latest"
		requests := []*balance.Request{
			&balance.Request{
				Address: "0xa838e871a02c6d883bf004352fc7dac8f781fed6",
				Currency:  balance.ETH,
				DefaultBlockParam:   block,
			},
			&balance.Request{
				Address: "0xa838e871a02c6d883bf004352fc7dac8f781fed6",
				Currency:  "0xBEB9eF514a379B997e0798FDcC901Ee474B6D9A1",
				DefaultBlockParam:   block,
			},
			&balance.Request{
				Address: "0xa838e871a02c6d883bf004352fc7dac8f781fed6",
				Currency:  "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07",
				DefaultBlockParam:   block,
			},
			&balance.Request{
				Address: "0xa838e871a02c6d883bf004352fc7dac8f781fed6",
				Currency:  "0x8aa33a7899fcc8ea5fbe6a608a109c3893a1b8b2",
				DefaultBlockParam:   block,
			},
			&balance.Request{
				Address: "0xa838e871a02c6d883bf004352fc7dac8f781fed6",
				Currency: "0x0f5d2fb29fb7d3cfee444a200298f468908cc942",
				DefaultBlockParam:  block,
			},
		}
		rawBalances, err := b.GetRawSheet(requests)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("--------------- Raw Balances -----------------")
		spew.Dump(rawBalances)
		intBalances, err := b.GetIntSheet(requests)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("--------------- big.Int Balances -----------------")
		spew.Dump(intBalances)
	}
}
