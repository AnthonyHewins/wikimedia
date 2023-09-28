package pages

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/AnthonyHewins/wikimedia/adapter"
)

type GetPageResp struct {
	ID           int        `json:"id"`
	Key          string     `json:"key"`
	Title        string     `json:"title"`
	Latest       LatestPage `json:"latest"`
	ContentModel string     `json:"content_model"`
	License      License    `json:"license"`
	HTMLURL      string     `json:"html_url"`
}

type LatestPage struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

type License struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

// GetPage is a request object for fetching a page. Create it with NewGetPageReq,
// and pass the required parameters to create it. Resolve it when you're ready
// with a context
//
//	client.NewPageReq("", "", "")
type GetPage struct {
	a *adapter.Core

	opts []func(*http.Request)

	// Required: Project name. For example: wikipedia (encyclopedia articles), commons (images, audio, and video), wiktionary (dictionary entries)
	Project string

	// Required? Here's what the docs say:
	// Language code. For example: ar (Arabic), en (English), es (Spanish).
	// Note: The language parameter is prohibited for commons and other multilingual projects.
	Language string

	// name of the article
	Title string
}

func (c *Client) NewGetPageReq(project, language, title string) *GetPage {
	return &GetPage{
		a:        c.a,
		Project:  project,
		Language: language,
		Title:    title,
	}
}

func (g *GetPage) WithOpts(opts ...func(*http.Request)) *GetPage {
	g.opts = opts
	return g
}

func (g *GetPage) Resolve(ctx context.Context) (*GetPageResp, error) {
	path, err := url.JoinPath(g.a.BaseURL, "core/v1", g.Project, g.Language, "page", g.Title, "bare")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	for _, opt := range g.opts {
		opt(req)
	}

	buf, err := g.a.Do(req)
	if err != nil {
		return nil, err
	}

	var getResp GetPageResp
	if err = json.Unmarshal(buf, &getResp); err != nil {
		return nil, err
	}

	return &getResp, nil
}
