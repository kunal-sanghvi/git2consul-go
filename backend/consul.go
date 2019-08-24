package backend

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type consul struct {
	host string
}

func (c *consul) Populate(key string, value string) error {
	req, err := http.NewRequest("PUT", c.host + "/v1/kv/" + key, strings.NewReader(value))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("Non 200 response code %d", resp.StatusCode))
	}
	return nil
}

func NewConsulBackend(host string) ConfigBackend {
	return &consul{
		host: host,
	}
}