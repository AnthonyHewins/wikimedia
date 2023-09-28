package wikimedia

import (
	"github.com/AnthonyHewins/wikimedia/adapter"
	"github.com/AnthonyHewins/wikimedia/pages"
)

type Client struct {
	Pages *pages.Client
}

func NewClient(core *adapter.Core) *Client {
	return &Client{
		Pages: pages.NewClient(core),
	}
}
