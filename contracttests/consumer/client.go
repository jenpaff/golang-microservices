package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/log"
	"github.com/jenpaff/golang-microservices/common"
	"github.com/jenpaff/golang-microservices/errors"
	"io/ioutil"
	"net/http"
)

type UserClient struct {
	BaseURL string
}

func (c *UserClient) GetUser(userName string) (*common.User, error) {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL+"/users/"+userName, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request due to error : %w", err)
	}
	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("retrieving user with username %s failed due to err %s : %w", userName, err.Error(), errors.UserClientError)
		log.Errorf(err.Error())
		return nil, err
	}
	defer response.Body.Close()
	statusCode := response.StatusCode
	if statusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(response.Body)
		err = fmt.Errorf("retrieving user with username %s failed with statusCode %d, status %s, err %s : %w", userName, statusCode, response.Status, string(respBody), errors.UserClientError)
		log.Errorf(err.Error())
		return nil, err

	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("could not read response: %s : %w", err.Error(), errors.UserClientError)
		log.Errorf(err.Error())
		return nil, err
	}

	var user common.User
	err = json.Unmarshal(respBody, &user)
	if err != nil {
		err = fmt.Errorf("could not parse response body: %s : %w", err.Error(), errors.UserClientError)
		log.Errorf(err.Error())
		return nil, err
	}

	return &user, nil
}
