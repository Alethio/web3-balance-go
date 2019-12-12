### Web3 Balance

An extension of [web3-go](https://github.com/Alethio/web3-go) that is used for reading token or ether balances from the chain.

#### Example

```go
eth, err := ethrpc.NewWithDefaults(ethClientURL)
b := balance.New(eth, 5)
block := fmt.Sprintf("0x%x", 7000000)

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
// Outputs:
// --------------- Raw Balances -----------------
// (map[string]map[string]map[balance.Currency]string) (len=1) {
//  (string) (len=8) "0x6acfc0": (map[string]map[balance.Currency]string) (len=1) {
//   (string) (len=42) "0xa838e871a02c6d883bf004352fc7dac8f781fed6": (map[balance.Currency]string) (len=5) {
//    (balance.Currency) (len=3) "ETH": (string) (len=17) "0x21264e1ec881d8a",
//    (balance.Currency) (len=42) "0x8aa33a7899fcc8ea5fbe6a608a109c3893a1b8b2": (string) (len=66) "0x0000000000000000000000000000000000000000000000c328093e61ee400000",
//    (balance.Currency) (len=42) "0xBEB9eF514a379B997e0798FDcC901Ee474B6D9A1": (string) (len=66) "0x000000000000000000000000000000000000000000000001d4ccdee9a074c000",
//    (balance.Currency) (len=42) "0x0f5d2fb29fb7d3cfee444a200298f468908cc942": (string) (len=66) "0x00000000000000000000000000000000000000000000017dbe4e10a6870d635c",
//    (balance.Currency) (len=42) "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07": (string) (len=66) "0x00000000000000000000000000000000000000000000000005afc055a2f44c2d"
//   }
//  }
// }
intBalances, err := b.GetIntSheet(requests)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println("--------------- big.Int Balances -----------------")
spew.Dump(intBalances)
// Outputs:
// --------------- big.Int Balances -----------------
// (map[string]map[string]map[balance.Currency]*big.Int) (len=1) {
//  (string) (len=8) "0x6acfc0": (map[string]map[balance.Currency]*big.Int) (len=1) {
//   (string) (len=42) "0xa838e871a02c6d883bf004352fc7dac8f781fed6": (map[balance.Currency]*big.Int) (len=5) {
//    (balance.Currency) (len=42) "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07": (*big.Int)(0xc0002b0140)(409757565152676909),
//    (balance.Currency) (len=42) "0x8aa33a7899fcc8ea5fbe6a608a109c3893a1b8b2": (*big.Int)(0xc0002b0180)(3600000000000000000000),
//    (balance.Currency) (len=42) "0x0f5d2fb29fb7d3cfee444a200298f468908cc942": (*big.Int)(0xc0002b01c0)(7041922408306145321820),
//    (balance.Currency) (len=3) "ETH": (*big.Int)(0xc0002b0260)(149292659155410314),
//    (balance.Currency) (len=42) "0xBEB9eF514a379B997e0798FDcC901Ee474B6D9A1": (*big.Int)(0xc0002b02e0)(33780620000000000000)
//   }
//  }
// }
// 


```
