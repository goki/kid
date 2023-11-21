// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kid

import (
	"context"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// Google authenticates the user with Google using [Auth] and returns the
// resulting oauth token and user info.
func Google(ctx context.Context, scopes ...string) (*oauth2.Token, *oidc.UserInfo, error) {
	return Auth(ctx, "google", "https://accounts.google.com", os.Getenv("GOOGLE_OAUTH2_CLIENT_ID"), os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET"), scopes...)
}
