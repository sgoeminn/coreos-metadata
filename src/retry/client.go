package retry

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	MaxAttempts    int
	Header         http.Header
}

func (c Client) Get(url string) ([]byte, error) {
	delay := c.InitialBackoff
	for attempt := 1; attempt <= c.MaxAttempts; attempt++ {
		fmt.Printf("fetching %q: attempt #%d\n", url, attempt)

		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		request.Header = c.Header

		if response, err := (&http.Client{}).Do(request); err != nil {
			fmt.Printf("failed to fetch: %v\n", err)
		} else if response.StatusCode != http.StatusOK {
			fmt.Printf("failed to fetch: %s\n", http.StatusText(response.StatusCode))
		} else {
			defer response.Body.Close()
			return ioutil.ReadAll(response.Body)
		}

		time.Sleep(delay)
		delay *= 2
		if delay > c.MaxBackoff {
			delay = c.MaxBackoff
		}
	}

	return nil, fmt.Errorf("timed out while fetching %q", url)
}

func (c Client) Getf(format string, a ...interface{}) ([]byte, error) {
	return c.Get(fmt.Sprintf(format, a...))
}
