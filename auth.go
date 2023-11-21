// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package kid (Ki ID) provides a system for identifying and authenticating
// users through third party cloud systems in GoKi apps.
package kid

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"goki.dev/goosi"
	"goki.dev/grows/jsons"
	"golang.org/x/oauth2"
)

// Auth authenticates the user with the given provider and returns the
// resulting oauth token and user info. The provider is specified in terms
// of its name (eg: "google") and its URL (eg: "https://accounts.google.com").
// Also, Auth uses the given Client ID and Client Secret for the app that needs
// the user information, which are typically obtained through a developer oauth
// portal (eg: the Credentials section of https://console.developers.google.com/).
// If the given token file is not "", Auth also saves the token to the file as JSON and
// skips the user-facing authentication step if it finds a valid token at the file
// (ie: remember me). By default, Auth requests the "openid", "profile", and "email"
// scopes, but more scopes can be specified on top of those via the scopes parameter.
func Auth(ctx context.Context, providerName, providerURL, clientID, clientSecret string, tokenFile string, scopes ...string) (*oauth2.Token, *oidc.UserInfo, error) {
	if clientID == "" || clientSecret == "" {
		slog.Warn("got empty client id or client secret; do you need to set env variables?")
	}

	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		return nil, nil, err
	}

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://127.0.0.1:5556/auth/" + providerName + "/callback",
		Endpoint:     provider.Endpoint(),
		Scopes:       append([]string{oidc.ScopeOpenID, "profile", "email"}, scopes...),
	}

	var token *oauth2.Token

	if tokenFile != "" {
		err := jsons.Open(&token, tokenFile)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, nil, err
		}
	}

	// if we didn't get it through remember me, we have to get it manually
	if token == nil {
		b := make([]byte, 16)
		rand.Read(b)
		state := base64.RawURLEncoding.EncodeToString(b)

		code := make(chan string)

		sm := http.NewServeMux()
		sm.HandleFunc("/auth/"+providerName+"/callback", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("state") != state {
				http.Error(w, "state did not match", http.StatusBadRequest)
				return
			}
			code <- r.URL.Query().Get("code")
			w.Write([]byte("<h1>Signed in</h1><p>You can return to the app</p>"))
		})
		// TODO(kai/auth): more graceful closing / error handling
		go http.ListenAndServe("127.0.0.1:5556", sm)

		goosi.TheApp.OpenURL(config.AuthCodeURL(state))

		cs := <-code

		token, err = config.Exchange(ctx, cs)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to exchange token: %w", err)
		}
		if tokenFile != "" {
			// TODO(kai/kid): more secure saving of token file
			err := jsons.Save(token, tokenFile)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	tokenSource := config.TokenSource(ctx, token)
	// the access token could have changed
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, nil, err
	}

	// if the access token changed, we have to save it again
	if newToken.AccessToken != token.AccessToken && tokenFile != "" {
		// TODO(kai/kid): more secure saving of token file
		err := jsons.Save(newToken, tokenFile)
		if err != nil {
			return nil, nil, err
		}
	}

	userInfo, err := provider.UserInfo(ctx, tokenSource)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user info: %w", err)
	}
	return newToken, userInfo, nil
}
