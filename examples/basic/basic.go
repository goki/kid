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
	sc := gi.NewScene("kid-basic").SetTitle("Kid Basic Example")
	kid.Buttons(sc, func(token *oauth2.Token, userInfo *oidc.UserInfo) {
		d := gi.NewDialog(sc).Title("User info:").FullWindow(true)
		giv.NewStructView(d).SetStruct(userInfo)
		claims := map[string]any{}
		grr.Log0(userInfo.Claims(&claims))
		giv.NewMapView(d).SetMap(&claims)
		d.Ok().Run()
	})
	gi.NewWindow(sc).Run().Wait()
}
