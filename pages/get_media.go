package pages

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/AnthonyHewins/wikimedia/adapter"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Latest struct {
	Timestamp time.Time `json:"timestamp"`
	User      User      `json:"user"`
}

type Image struct {
	Mediatype string `json:"mediatype"`
	SizeBytes *uint  `json:"size"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Duration  any    `json:"duration"`
	URL       string `json:"url"`
}

type File struct {
	Title              string `json:"title"`
	FileDescriptionURL string `json:"file_description_url"`
	Latest             Latest `json:"latest"`
	Preferred          Image  `json:"preferred"`
	Original           Image  `json:"original"`
}

type filesWrapper struct {
	Files []File `json:"files"`
}

// GetPage is a request object for fetching a page. Create it with NewGetPageReq,
// and pass the required parameters to create it. Resolve it when you're ready
// with a context
//
//	client.NewPageReq("", "", "")
type GetMedia struct {
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

func (c *Client) NewGetMediaReq(project, language, title string) *GetPage {
	return &GetPage{
		a:        c.a,
		Project:  project,
		Language: language,
		Title:    title,
	}
}

func (g *GetMedia) WithOpts(opts ...func(*http.Request)) *GetMedia {
	g.opts = opts
	return g
}

func (g *GetMedia) Resolve(ctx context.Context) ([]File, error) {
	path, err := url.JoinPath(g.a.BaseURL, "core/v1", g.Project, g.Language, "page", g.Title, "links/media")
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

	var getResp filesWrapper
	if err = json.Unmarshal(buf, &getResp); err != nil {
		return nil, err
	}

	return getResp.Files, nil
}
