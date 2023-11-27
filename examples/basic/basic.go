// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/gimain"
	"goki.dev/gi/v2/giv"
	"goki.dev/grr"
	"goki.dev/kid"
	"golang.org/x/oauth2"
)

func main() { gimain.Run(app) }

func app() {
	gi.SetAppName("kid-basic")
	b := gi.NewBody().SetTitle("Kid Basic Example")
	fun := func(token *oauth2.Token, userInfo *oidc.UserInfo) {
		d := gi.NewBody().AddTitle("User info")
		gi.NewLabel(d).SetType(gi.LabelHeadlineMedium).SetText("Basic info")
		giv.NewStructView(d).SetStruct(userInfo)
		gi.NewLabel(d).SetType(gi.LabelHeadlineMedium).SetText("Detailed info")
		claims := map[string]any{}
		grr.Log(userInfo.Claims(&claims))
		giv.NewMapView(d).SetMap(&claims)
		d.AddBottomBar(func(pw gi.Widget) {
			d.AddOk(pw)
		})
		d.NewFullDialog(b).Run()
	}
	kid.Buttons(b, &kid.ButtonsConfig{SuccessFunc: fun})
	b.NewWindow().Run().Wait()
}
