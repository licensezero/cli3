package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Client responds to broker API requests.
type Client interface {
	Offer(api string, offerID string) (*Offer, error)
	Seller(api string, sellerID string) (*Seller, error)
}

type httpClient struct {
	Client *http.Client
}

// NewClient reutrns a client using the given Transport.
func NewClient(t http.RoundTripper) Client {
	return &httpClient{Client: &http.Client{Transport: t}}
}

func (c *httpClient) Offer(api string, offerID string) (offer *Offer, err error) {
	response, err := c.Client.Get(api + "/offers/" + offerID)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &offer)
	if err != nil {
		return nil, errors.New("invalid JSON")
	}
	err = offer.Validate()
	if err != nil {
		return nil, errors.New("invalid offer")
	}
	return
}

func (c *httpClient) Seller(api string, sellerID string) (seller *Seller, err error) {
	response, err := c.Client.Get(api + "/sellers/" + sellerID)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &seller)
	if err != nil {
		return nil, errors.New("invalid JSON")
	}
	err = seller.Validate()
	if err != nil {
		return nil, errors.New("invalid seller")
	}
	return
}
