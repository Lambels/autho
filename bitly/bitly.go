package bitly

import (
	"encoding/json"
	"net/http"

	"github.com/Lambels/autho"
	"golang.org/x/oauth2"
)

const profileEndpoint string = "https://api-ssl.bitly.com/v4/user"

var Endpoint *oauth2.Endpoint = &oauth2.Endpoint{
	AuthURL:  "https://bitly.com/oauth/authorize",
	TokenURL: "https://api-ssl.bitly.com/oauth/access_token",
}

// User represents fields accessible on a bitly account.
//
// https://dev.bitly.com/api-reference/#getUser
type User struct {
	Login            string   `json:"login"`
	Name             string   `json:"name"`
	IsActive         bool     `json:"is_active"`
	Created          string   `json:"created"`
	Modified         string   `json:"modified"`
	IsSSOUser        bool     `json:"is_sso_user"`
	Emails           []*Email `json:"emails"`
	Is2FAEnabled     bool     `json:"is_2fa_enabled"`
	DefaultGroupGuid string   `json:"default_group_guid"`
}

// Email represents an email on the users account.
type Email struct {
	Email      string `json:"email"`
	IsPrimary  string `json:"is_primary"`
	IsVerified string `json:"is_verified"`
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
	if resp.StatusCode != http.StatusOK {
		return nil, autho.ErrNoUser
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, autho.ErrNoUser
	}

	return &user, nil
}
