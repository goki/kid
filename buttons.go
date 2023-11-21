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
	"goki.dev/goosi/events"
	"goki.dev/icons"
	"golang.org/x/oauth2"
)

//go:embed svg/*.svg
var providerIcons embed.FS

func init() {
	icons.Icons = merged_fs.NewMergedFS(icons.Icons, providerIcons)
}

// Buttons adds a new vertical layout to the given parent with authentication
// buttons for major platforms. It calls the given function with the resulting
// oauth token and user info when the user successfully authenticates. See [Auth]
// for more information about scopes.
func Buttons(par gi.Widget, fun func(token *oauth2.Token, userInfo *oidc.UserInfo), scopes ...string) *gi.Layout {
	ly := gi.NewLayout(par, "auth-buttons")
	ly.Style(func(s *styles.Style) {
		s.Direction = styles.Column
	})
	GoogleButton(ly, fun, scopes...)
	return ly
}

// GoogleButton adds a new button for signing in with Google.
// It calls the given function when the token and user info are obtained.
// See [Auth] for more information about scopes.
func GoogleButton(par gi.Widget, fun func(token *oauth2.Token, userInfo *oidc.UserInfo), scopes ...string) *gi.Button {
	bt := gi.NewButton(par, "sign-in-with-google").SetType(gi.ButtonOutlined).
		SetText("Sign in with Google").SetIcon("sign-in-with-google")
	bt.Style(func(s *styles.Style) {
		s.Color = colors.Scheme.OnSurface
	})
	bt.OnClick(func(e events.Event) {
		token, userInfo, err := Google(context.TODO(), scopes...)
		if err != nil {
			gi.ErrorDialog(par, err, "Error signing in with Google")
		}
		fun(token, userInfo)
	})
	return bt
}
