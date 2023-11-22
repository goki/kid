// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"path/filepath"

	"github.com/coreos/go-oidc/v3/oidc"
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/gimain"
	"goki.dev/gi/v2/giv"
	"goki.dev/goosi"
	"goki.dev/goosi/events"
	"goki.dev/grr"
	"goki.dev/kid"
	"golang.org/x/oauth2"
)

func main() { gimain.Run(app) }

func app() {
	gi.SetAppName("kid-scopes")
	b := gi.NewBody().SetTitle("Kid Scopes and Token File Example")
	fun := func(token *oauth2.Token, userInfo *oidc.UserInfo) {
		b.Sc.OnShow(func(e events.Event) {
			d := gi.NewBody().AddTitle("User info")
			gi.NewLabel(d).SetType(gi.LabelHeadlineMedium).SetText("Basic info")
			giv.NewStructView(d).SetStruct(userInfo)
			gi.NewLabel(d).SetType(gi.LabelHeadlineMedium).SetText("Detailed info")
			claims := map[string]any{}
			grr.Log0(userInfo.Claims(&claims))
			giv.NewMapView(d).SetMap(&claims)
			d.AddBottomBar(func(pw gi.Widget) {
				d.AddOk(pw)
			})
			ds := d.NewFullDialog(b)
			ds.Run()
		})
	}
	kid.Buttons(b, &kid.ButtonsConfig{
		SuccessFunc: fun,
		TokenFile: func(provider string) string {
			return filepath.Join(goosi.TheApp.AppPrefsDir(), provider+"-token.json")
		},
		Scopes: map[string][]string{
			"google": {"https://mail.google.com/"},
		},
	})
	b.NewWindow().Run().Wait()
}
