// Copyright (c) 2023, The Goki Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"path/filepath"

	"github.com/coreos/go-oidc/v3/oidc"
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/giv"
	"goki.dev/grr"
	"goki.dev/kid"
	"golang.org/x/oauth2"
)

func main() {
	b := gi.NewAppBody("kid-scopes").SetTitle("Kid Scopes and Token File Example")
	fun := func(token *oauth2.Token, userInfo *oidc.UserInfo) {
		d := gi.NewBody()
		gi.NewLabel(d).SetType(gi.LabelHeadlineMedium).SetText("Basic info")
		giv.NewStructView(d).SetStruct(userInfo)
		gi.NewLabel(d).SetType(gi.LabelHeadlineMedium).SetText("Detailed info")
		claims := map[string]any{}
		grr.Log(userInfo.Claims(&claims))
		giv.NewMapView(d).SetMap(&claims)
		d.AddBottomBar(func(pw gi.Widget) {
			d.AddOk(pw)
		})
		ds := d.NewFullDialog(b)
		ds.Run()
	}
	kid.Buttons(b, &kid.ButtonsConfig{
		SuccessFunc: fun,
		TokenFile: func(provider, email string) string {
			return filepath.Join(b.App().DataDir(), provider+"-token.json")
		},
		Scopes: map[string][]string{
			"google": {"https://mail.google.com/"},
		},
	})
	b.NewWindow().Run().Wait()
}
