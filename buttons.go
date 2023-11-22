// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kid

import (
	"context"
	"embed"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/yalue/merged_fs"
	"goki.dev/colors"
	"goki.dev/gi/v2/gi"
	"goki.dev/girl/styles"
	"goki.dev/glop/dirs"
	"goki.dev/goosi/events"
	"goki.dev/icons"
	"golang.org/x/oauth2"
)

//go:embed svg/*.svg
var providerIcons embed.FS

func init() {
	icons.Icons = merged_fs.NewMergedFS(icons.Icons, providerIcons)
}

// ButtonsConfig is the configuration information passed to [Buttons].
type ButtonsConfig struct {
	// SuccessFunc, if non-nil, is the function called after the user successfully
	// authenticates. It is passed the user's authentication token and info.
	SuccessFunc func(token *oauth2.Token, userInfo *oidc.UserInfo)

	// TokenFile, if non-nil, is the function used to determine what token file is
	// passed to [Auth]. It is passed the provider being used (eg: "google").
	TokenFile func(provider string) string

	// Scopes, if non-nil, is a map of scopes to pass to [Auth], keyed by the
	// provider being used (eg: "google").
	Scopes map[string][]string
}

// Buttons adds a new vertical layout to the given parent with authentication
// buttons for major platforms, using the given configuration options. See
// [ButtonsConfig] for more information on the configuration options. The
// configuration options can be nil, in which case default values will be used.
func Buttons(par gi.Widget, c *ButtonsConfig) *gi.Layout {
	ly := gi.NewLayout(par, "auth-buttons")
	ly.Style(func(s *styles.Style) {
		s.Direction = styles.Column
	})

	if c == nil {
		c = &ButtonsConfig{}
	}
	if c.SuccessFunc == nil {
		c.SuccessFunc = func(token *oauth2.Token, userInfo *oidc.UserInfo) {}
	}
	if c.TokenFile == nil {
		c.TokenFile = func(provider string) string { return "" }
	}
	if c.Scopes == nil {
		c.Scopes = map[string][]string{}
	}

	GoogleButton(ly, c.SuccessFunc, c.TokenFile("google"), c.Scopes["google"]...)
	return ly
}

// GoogleButton adds a new button for signing in with Google.
// It calls the given function when the token and user info are obtained.
// See [Auth] for more information about token files and scopes.
func GoogleButton(par gi.Widget, fun func(token *oauth2.Token, userInfo *oidc.UserInfo), tokenFile string, scopes ...string) *gi.Button {
	bt := gi.NewButton(par, "sign-in-with-google").SetType(gi.ButtonOutlined).
		SetText("Sign in with Google").SetIcon("sign-in-with-google")
	bt.Style(func(s *styles.Style) {
		s.Color = colors.Scheme.OnSurface
	})

	auth := func() {
		token, userInfo, err := Google(context.TODO(), tokenFile, scopes...)
		if err != nil {
			gi.ErrorDialog(bt, err, "Error signing in with Google").Run()
			return
		}
		fun(token, userInfo)
	}
	bt.OnClick(func(e events.Event) {
		auth()
	})

	// if we have a valid token file, we auth immediately without the user clicking on the button
	if tokenFile != "" {
		exists, err := dirs.FileExists(tokenFile)
		if err != nil {
			gi.ErrorDialog(bt, err, "Error searching for saved Google auth token file").Run()
			return bt
		}
		if exists {
			// have to wait until the scene is shown in case any dialogs are created
			bt.OnShow(func(e events.Event) {
				auth()
			})
		}
	}
	return bt
}
