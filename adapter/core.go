package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// MediawikiErr is the common error struct for errors in the API.
// Error details are placed into a JSON object called messageTranslations
// to handle multiple languages; therefore getting the details is a little difficult.
// The Error() method joins them all together
//
// Here's an example of what comes back on error, and what you can expect
// if you need to use it:
//
//	{
//	  "messageTranslations": {
//	    "en": "The specified title does not exist"
//	  },
//	  "httpCode": 404,
//	  "httpReason": "Not Found"
//	}
type MediawikiErr struct {
	HTTPCode            int               `json:"httpCode"`
	HTTPReason          string            `json:"httpReason"`
	MessageTranslations map[string]string `json:"messageTranslations"`
}

func (a *MediawikiErr) Is(err error) bool {
	_, ok := err.(*MediawikiErr)
	return ok
}

func (a *MediawikiErr) Error() string {
	var sb strings.Builder
	i := 0
	for k, v := range a.MessageTranslations {
		if i != 0 {
			sb.WriteRune(';')
		}

		sb.WriteString(k)
		sb.WriteRune(':')
		sb.WriteString(v)
	}

	return fmt.Sprintf("%d: %s", a.HTTPCode, sb.String())
}

type Core struct {
	BaseURL    string
	HTTPClient *http.Client
}

// Do will perform HTTP requests using the specified HTTP client and
// check for any wikipedia errors
func (c *Core) Do(req *http.Request) ([]byte, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	if s := resp.StatusCode; s >= 200 && s < 300 {
		return buf, nil
	}

	var x MediawikiErr
	if err := json.Unmarshal(buf, &x); err != nil {
		return nil, fmt.Errorf("server did not return valid JSON: %w", err)
	}

	return nil, &x
}
