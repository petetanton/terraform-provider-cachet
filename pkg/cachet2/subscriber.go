package cachet2

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Subscriber struct {
	ID     int    `json:"id,omitempty"`
	Email  string `json:"email"`
	Verify bool   `json:"verify"`
}

type SubscriberGetResponse struct {
	Subscribers []*Subscriber `json:"data"`
}

type SubscriberPostResponse struct {
	Subscriber *Subscriber `json:"data"`
}

func (c *Client) SubscribersGetAll() ([]*Subscriber, error) {
	path := "api/v1/subscribers?per_page=10000"

	response, err := c.get(path)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got a %d status code when trying to get subscribers: %s", response.StatusCode, string(responseBytes))
	}

	var subscriberResponse SubscriberGetResponse
	if err := json.Unmarshal(responseBytes, &subscriberResponse); err != nil {
		return nil, err
	}

	return subscriberResponse.Subscribers, nil
}

func (c *Client) SubscriberGet(id int) (*Subscriber, error) {
	subscribers, err := c.SubscribersGetAll()
	if err != nil {
		return nil, err
	}

	for _, subscriber := range subscribers {
		if subscriber.ID == id {
			return subscriber, nil
		}
	}

	return nil, nil
}

func (c *Client) SubscriberDelete(id int) error {
	path := fmt.Sprintf("api/v1/subscribers/%d", id)

	response, err := c.delete(path)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("got a %d status code when trying to delete subscriber %d", response.StatusCode, id)
	}

	return nil
}

func (c *Client) SubscriberCreate(s *Subscriber) (*Subscriber, error) {
	path := "api/v1/subscribers"

	reqBytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	response, err := c.post(path, reqBytes)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got a %d status code when trying to create subscriber: %s", response.StatusCode, string(responseBytes))
	}

	var subscriberResponse SubscriberPostResponse
	if err := json.Unmarshal(responseBytes, &subscriberResponse); err != nil {
		return nil, err
	}

	return subscriberResponse.Subscriber, nil
}
