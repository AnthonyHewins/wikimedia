package pages

import (
	"github.com/AnthonyHewins/wikimedia/adapter"
)

type Client struct {
	a *adapter.Core
}

func NewClient(coreAdapter *adapter.Core) *Client {
	return &Client{a: coreAdapter}
}
