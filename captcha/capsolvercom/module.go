package capsolvercom

import "github.com/go-resty/resty/v2"

type Provider struct {
	s      *resty.Client
	apiKey string
}

func New(apiKey string) *Provider {
	return &Provider{
		s:      resty.New(),
		apiKey: apiKey,
	}
}
