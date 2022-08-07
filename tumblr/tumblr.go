package tumblr

import (
	"encoding/json"
	"net/http"

	"github.com/Lambels/autho"
)

const profileEndpoint string = "https://api.tumblr.com/v2/user/info"

type meta struct {
	Status int    `json:"stauts"`
	Msg    string `json:"msg"`
}

// User represents fields accessible on a tumblr account.
//
// https://www.tumblr.com/docs/en/api/v2#userinfo--get-a-users-information
type User struct {
	Name string `json:"name"`
	// The ammout of people the user follows.
	Following int `json:"following"`
	Likes     int `json:"likes"`
}

type response struct {
	Metadata meta  `json:"meta"`
	UserInfo *User `json:"response"`
}

func me(client *http.Client) (*User, error) {
	req, err := http.NewRequest(http.MethodGet, profileEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return validateResp(resp)
}

func validateResp(resp *http.Response) (*User, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, autho.ErrNoUser
	}

	var data response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, autho.ErrNoUser
	}

	return data.UserInfo, nil
}
