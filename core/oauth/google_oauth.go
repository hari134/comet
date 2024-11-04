package oauth

import (
	"context"
	"net/http"
	"os"

	oauthTypes "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
)

var googleOAuthConfig *oauthTypes.Config

func init() {
	googleOAuthConfig = &oauthTypes.Config{
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"),
		Scopes:       []string{oauth2.UserinfoEmailScope, oauth2.UserinfoProfileScope},
		Endpoint:     google.Endpoint,
	}
}

func GetGoogleOAuthURL(state string) string {
	return googleOAuthConfig.AuthCodeURL(state, oauthTypes.AccessTypeOffline)
}

func HandleGoogleOAuthCallback(r *http.Request) (*oauthTypes.Token, *oauth2.Userinfo, error) {
	ctx := context.Background()

	code := r.URL.Query().Get("code")
	if code == "" {
		return nil, nil, ErrMissingCode
	}

	token, err := googleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, nil, err
	}

	userInfo, err := fetchGoogleUserInfo(ctx, token)
	if err != nil {
		return nil, nil, err
	}

	return token, userInfo, nil
}

func fetchGoogleUserInfo(ctx context.Context, token *oauthTypes.Token) (*oauth2.Userinfo, error) {
	client := googleOAuthConfig.Client(ctx, token)

	oauth2Service, err := oauth2.New(client)
	if err != nil {
		return nil, err
	}

	userInfoService := oauth2Service.Userinfo
	userInfo, err := userInfoService.Get().Do()
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

var ErrMissingCode = http.ErrNoLocation
