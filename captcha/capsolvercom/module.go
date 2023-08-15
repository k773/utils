package capsolvercom

import "github.com/go-resty/resty/v2"

type Provider struct {
	S      *resty.Client
	ApiKey string
}

func New(apiKey string) *Provider {
	return &Provider{
		S:      resty.New(),
		ApiKey: apiKey,
	}
}
