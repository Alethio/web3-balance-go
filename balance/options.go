package balance

type Option func(*Config)

type Config struct {
	Retry bool
	Attempts uint
}

func RetryOnError(attempts uint) Option {
	return func(c *Config) {
		c.Retry = true;
		c.Attempts = attempts
	}
}
