// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kid

import (
	"context"
	"embed"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"goki.dev/icons"
	"golang.org/x/oauth2"

	"github.com/yalue/merged_fs"
)

var (
	clientID     = os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET")
)

//go:embed svg/google.svg
var googleIcon embed.FS

func init() {
	icons.Icons = merged_fs.NewMergedFS(icons.Icons, googleIcon)
}

// Google authenticates the user with Google and returns the
// resulting oauth token and user info.
func Google(ctx context.Context) (*oauth2.Token, *oidc.UserInfo, error) {
	return Auth(ctx, "google", "https://accounts.google.com")
}
