package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func BuildAuthorizationRequest(config Config) string {
	return fmt.Sprintf(
		"http://localhost:8080/auth/realms/%v/protocol/openid-connect/auth?client_id=%v&redirect_uri=%v&response_mode=query&response_type=code&scope=openid",
		config.KeycloakConfig.Realm,
		config.KeycloakConfig.ClientID,
		config.EmbeddedServerConfig.GetCallbackURL(),
	)
}

func BuildTokenExchangeRequest(config Config, code string) (*http.Request, error) {
	tokenURL := fmt.Sprintf("%v/realms/%v/protocol/openid-connect/token", config.KeycloakConfig.KeycloakURL, config.KeycloakConfig.Realm)

	body := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"client_id":    {config.KeycloakConfig.ClientID},
		"redirect_uri": {config.EmbeddedServerConfig.GetCallbackURL()},
	}

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, err
}
