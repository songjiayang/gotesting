package httpmocktesting

import (
	"encoding/json"
	"errors"
	post "gotesting/post"
	"io/ioutil"
	"net/http"
)

type PostClient struct {
	*http.Client
}

func (c *PostClient) Fetch(url string, page int) ([]*post.PostModel, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return nil, errors.New(string(data))
	}

	var items []*post.PostModel
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	return items, nil
}
