// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui  // import "miniflux.app/ui"

import (
	"net/http"

	"miniflux.app/http/context"
	"miniflux.app/http/response/html"
	"miniflux.app/locale"
	"miniflux.app/model"
	"miniflux.app/ui/form"
	"miniflux.app/ui/session"
	"miniflux.app/ui/view"
)

// ShowSettings shows the settings page.
func (c *Controller) ShowSettings(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)
	sess := session.New(c.store, ctx)
	view := view.New(c.tpl, ctx, sess)

	user, err := c.store.UserByID(ctx.UserID())
	if err != nil {
		html.ServerError(w, err)
		return
	}

	settingsForm := form.SettingsForm{
		Username:       user.Username,
		Theme:          user.Theme,
		Language:       user.Language,
		Timezone:       user.Timezone,
		EntryDirection: user.EntryDirection,
	}

	timezones, err := c.store.Timezones()
	if err != nil {
		html.ServerError(w, err)
		return
	}

	view.Set("form", settingsForm)
	view.Set("themes", model.Themes())
	view.Set("languages", locale.AvailableLanguages())
	view.Set("timezones", timezones)
	view.Set("menu", "settings")
	view.Set("user", user)
	view.Set("countUnread", c.store.CountUnreadEntries(user.ID))

	html.OK(w, r, view.Render("settings"))
}
