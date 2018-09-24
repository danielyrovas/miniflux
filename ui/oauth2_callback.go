// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui  // import "miniflux.app/ui"

import (
	"net/http"

	"miniflux.app/http/cookie"
	"miniflux.app/http/request"
	"miniflux.app/http/response"
	"miniflux.app/http/response/html"
	"miniflux.app/http/route"
	"miniflux.app/locale"
	"miniflux.app/logger"
	"miniflux.app/model"
	"miniflux.app/ui/session"
)

// OAuth2Callback receives the authorization code and create a new session.
func (c *Controller) OAuth2Callback(w http.ResponseWriter, r *http.Request) {
	printer := locale.NewPrinter(request.UserLanguage(r))
	sess := session.New(c.store, request.SessionID(r))

	provider := request.RouteStringParam(r, "provider")
	if provider == "" {
		logger.Error("[OAuth2] Invalid or missing provider")
		response.Redirect(w, r, route.Path(c.router, "login"))
		return
	}

	code := request.QueryStringParam(r, "code", "")
	if code == "" {
		logger.Error("[OAuth2] No code received on callback")
		response.Redirect(w, r, route.Path(c.router, "login"))
		return
	}

	state := request.QueryStringParam(r, "state", "")
	if state == "" || state != request.OAuth2State(r) {
		logger.Error(`[OAuth2] Invalid state value: got "%s" instead of "%s"`, state, request.OAuth2State(r))
		response.Redirect(w, r, route.Path(c.router, "login"))
		return
	}

	authProvider, err := getOAuth2Manager(c.cfg).Provider(provider)
	if err != nil {
		logger.Error("[OAuth2] %v", err)
		response.Redirect(w, r, route.Path(c.router, "login"))
		return
	}

	profile, err := authProvider.GetProfile(code)
	if err != nil {
		logger.Error("[OAuth2] %v", err)
		response.Redirect(w, r, route.Path(c.router, "login"))
		return
	}

	if request.IsAuthenticated(r) {
		user, err := c.store.UserByExtraField(profile.Key, profile.ID)
		if err != nil {
			html.ServerError(w, err)
			return
		}

		if user != nil {
			logger.Error("[OAuth2] User #%d cannot be associated because %s is already associated", request.UserID(r), user.Username)
			sess.NewFlashErrorMessage(printer.Printf("error.duplicate_linked_account"))
			response.Redirect(w, r, route.Path(c.router, "settings"))
			return
		}

		if err := c.store.UpdateExtraField(request.UserID(r), profile.Key, profile.ID); err != nil {
			html.ServerError(w, err)
			return
		}

		sess.NewFlashMessage(printer.Printf("alert.account_linked"))
		response.Redirect(w, r, route.Path(c.router, "settings"))
		return
	}

	user, err := c.store.UserByExtraField(profile.Key, profile.ID)
	if err != nil {
		html.ServerError(w, err)
		return
	}

	if user == nil {
		if !c.cfg.IsOAuth2UserCreationAllowed() {
			html.Forbidden(w)
			return
		}

		user = model.NewUser()
		user.Username = profile.Username
		user.IsAdmin = false
		user.Extra[profile.Key] = profile.ID

		if err := c.store.CreateUser(user); err != nil {
			html.ServerError(w, err)
			return
		}
	}

	sessionToken, _, err := c.store.CreateUserSession(user.Username, r.UserAgent(), request.ClientIP(r))
	if err != nil {
		html.ServerError(w, err)
		return
	}

	logger.Info("[Controller:OAuth2Callback] username=%s just logged in", user.Username)
	c.store.SetLastLogin(user.ID)
	sess.SetLanguage(user.Language)
	sess.SetTheme(user.Theme)

	http.SetCookie(w, cookie.New(
		cookie.CookieUserSessionID,
		sessionToken,
		c.cfg.IsHTTPS,
		c.cfg.BasePath(),
	))

	response.Redirect(w, r, route.Path(c.router, "unread"))
}
