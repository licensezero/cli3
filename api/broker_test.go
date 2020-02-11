package api

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/licensezero/helptest"
	"net/http"
	"testing"
)

func TestInvalidResponses(t *testing.T) {
	brokerAPI := "https://broker.licensezero.com"
	offerID := uuid.New().String()
	sellerID := uuid.New().String()
	transport := helptest.RoundTripFunc(func(req *http.Request) *http.Response {
		url := req.URL.String()
		if url == brokerAPI+"/offers/"+offerID ||
			url == brokerAPI+"/sellers/"+sellerID {
			return &http.Response{
				StatusCode: 200,
				Body:       helptest.NoopCloser{bytes.NewBufferString(`{"invalid":"response"}`)},
				Header:     make(http.Header),
			}
		}
		return &http.Response{
			StatusCode: 404,
			Body:       helptest.NoopCloser{bytes.NewBufferString("")},
			Header:     make(http.Header),
		}
	})

	client := http.Client{Transport: transport}
	brokerServer := BrokerServer{Client: &client, Base: brokerAPI}

	_, err := brokerServer.Offer(offerID)
	if err.Error() != "invalid offer" {
		t.Error("failed to return invalid offer error")
	}

	_, err = brokerServer.Seller(sellerID)
	if err.Error() != "invalid seller" {
		t.Error("failed to return invalid seller error")
	}
}

func TestInvalidJSON(t *testing.T) {
	brokerAPI := "https://licensezero.com"
	offerID := "186d34a9-c8f7-414c-91bc-a34b4553b91d"
	sellerID := "2d7aa023-8eca-4d15-ad43-6b15768a5293"
	transport := helptest.RoundTripFunc(func(req *http.Request) *http.Response {
		url := req.URL.String()
		if url == brokerAPI+"/offers/"+offerID ||
			url == brokerAPI+"/sellers/"+sellerID {
			return &http.Response{
				StatusCode: 200,
				Body:       helptest.NoopCloser{Reader: bytes.NewBufferString(`notvalidjson`)},
				Header:     make(http.Header),
			}
		}
		return &http.Response{
			StatusCode: 404,
			Body:       helptest.NoopCloser{Reader: bytes.NewBufferString("")},
			Header:     make(http.Header),
		}
	})

	client := http.Client{Transport: transport}
	brokerServer := BrokerServer{Client: &client, Base: brokerAPI}

	_, err := brokerServer.Offer(offerID)
	if err.Error() != "invalid JSON" {
		t.Error("failed to return invalid offer error")
	}

	_, err = brokerServer.Seller(sellerID)
	if err.Error() != "invalid JSON" {
		t.Error("failed to return invalid seller error")
	}
}