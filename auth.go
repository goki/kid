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
	"fmt"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"goki.dev/goosi"
	"golang.org/x/oauth2"
)

// Auth authenticates the user with the given provider and returns the
// resulting oauth token and user info. The provider is specified in terms
// of its name (eg: "google") and its URL (eg: "https://accounts.google.com").
// Also, Auth uses the given Client ID and Client Secret for the app that needs
// the user information, which are typically obtained through a developer oauth
// portal (eg: the Credentials section of https://console.developers.google.com/).
func Auth(ctx context.Context, providerName, providerURL, clientID, clientSecret string) (*oauth2.Token, *oidc.UserInfo, error) {
	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		return nil, nil, err
	}

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://127.0.0.1:5556/auth/" + providerName + "/callback",
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	fmt.Println(config)

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

	oauth2Token, err := config.Exchange(ctx, cs)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	userInfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		return oauth2Token, nil, fmt.Errorf("failed to get user info: %w", err)
	}
	return oauth2Token, userInfo, nil
}