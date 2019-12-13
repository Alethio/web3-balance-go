package balance

import "github.com/alethio/web3-multicall-go/multicall"

type Option func(*Config)

type Config struct {
	Retry bool
	Attempts uint
	Multicall *multicall.Multicall
	Loader Loader

}

type Loader interface {
	fetchRequests(b *Bookkeeper, requests []*Request, results chan *RawResponse, done chan error)
}


func RetryOnError(attempts uint) Option {
	return func(c *Config) {
		c.Retry = true;
		c.Attempts = attempts
	}
}

func UseMulticall(mc *multicall.Multicall) Option {
	return func(c *Config) {
		c.Loader = multicallLoader{mc}
	}
}

