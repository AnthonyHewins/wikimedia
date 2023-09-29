package pages

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/AnthonyHewins/wikimedia/adapter"
	"github.com/stretchr/testify/assert"
)

var (
	mockErr = adapter.MediawikiErr{
		HTTPCode:   404,
		HTTPReason: "not found",
		MessageTranslations: map[string]string{
			"en": "failed",
		},
	}
)

func TestGetMedia(mainTest *testing.T) {
	nine := uint(932)
	testCases := []struct {
		name string
		arg  GetMedia

		expected    []File
		expectedErr error

		expectedRequestPath string

		mockStatus int
		mockResp   any
	}{
		{
			name:                "base case",
			expectedRequestPath: "core/v1/page/links/media",
			expected:            []File{{}},
			mockStatus:          200,
			mockResp:            filesWrapper{[]File{{}}},
		},
		{
			name:                "failed request",
			expectedErr:         &mockErr,
			expectedRequestPath: "core/v1/page/links/media",
			mockStatus:          404,
			mockResp:            &mockErr,
		},
		{
			name: "with options added",
			arg: GetMedia{
				Project:  "project",
				Language: "lang",
				Title:    "title",
			},
			expected:            []File{{}},
			expectedRequestPath: "core/v1/project/lang/page/title/links/media",
			mockStatus:          200,
			mockResp:            filesWrapper{[]File{{}}},
		},
		{
			name: "unmarshals correctly",
			arg: GetMedia{
				Project:  "project",
				Language: "lang",
				Title:    "title",
			},
			expected: []File{{
				Title:              "Commons-logo.svg",
				FileDescriptionURL: "//en.wikipedia.org/wiki/File:Commons-logo.svg",
				Latest: Latest{
					Timestamp: time.Time{},
					User: User{
						ID:   1,
						Name: "u",
					},
				},
				Preferred: Image{
					Mediatype: "DRAWING",
					SizeBytes: nil,
					Width:     446,
					Height:    599,
					Duration:  nil,
					URL:       "//upload.wikimedia.org/wikipedia/en/thumb/4/4a/Commons-logo.svg/446px-Commons-logo.svg.png",
				},
				Original: Image{
					Mediatype: "DRAWING",
					SizeBytes: &nine,
					Width:     1024,
					Height:    1376,
					Duration:  nil,
					URL:       "//upload.wikimedia.org/wikipedia/en/4/4a/Commons-logo.svg",
				},
			}},
			expectedRequestPath: "core/v1/project/lang/page/title/links/media",
			mockStatus:          200,
			mockResp: map[string][]any{"files": {
				map[string]any{
					"title":                "Commons-logo.svg",
					"file_description_url": "//en.wikipedia.org/wiki/File:Commons-logo.svg",
					"latest": map[string]any{
						"timestamp": nil,
						"user": map[string]any{
							"id":   1,
							"name": "u",
						},
					},
					"preferred": map[string]any{
						"mediatype": "DRAWING",
						"size":      nil,
						"width":     446,
						"height":    599,
						"duration":  nil,
						"url":       "//upload.wikimedia.org/wikipedia/en/thumb/4/4a/Commons-logo.svg/446px-Commons-logo.svg.png",
					},
					"original": map[string]any{
						"mediatype": "DRAWING",
						"size":      932,
						"width":     1024,
						"height":    1376,
						"duration":  nil,
						"url":       "//upload.wikimedia.org/wikipedia/en/4/4a/Commons-logo.svg",
					},
				},
			}},
		},
	}

	t := assert.New(mainTest)
	for _, tc := range testCases {
		tc.arg.a = &adapter.Core{
			HTTPClient: &http.Client{Timeout: time.Second},
		}

		ctx, ccl := context.WithTimeout(context.Background(), time.Millisecond*500)
		defer ccl()

		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.mockStatus)
			buf, err := json.Marshal(tc.mockResp)
			if err != nil {
				panic(err)
			}
			_, err = w.Write(buf)
			if err != nil {
				panic(err)
			}
		}))
		defer s.Close()

		tc.arg.WithOpts(func(r *http.Request) {
			t.Equal(tc.expectedRequestPath, r.URL.Path, "test case %s has the expected request URL", tc.name)
			path, err := url.Parse(s.URL)
			if err != nil {
				panic(err)
			}

			r.URL = path
		})

		actual, actualErr := tc.arg.Resolve(ctx)
		if t.Equal(tc.expectedErr, actualErr, "test case %s errors/succeeds as expected", tc.name) {
			t.Equal(tc.expected, actual, tc.name)
		}
	}
}
