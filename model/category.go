// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package model // import "miniflux.app/model"

import "fmt"

// Tag represents a feed Tag.
type Tag struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	UserID       int64  `json:"user_id"`
	HideGlobally bool   `json:"hide_globally"`
	FeedCount    int    `json:"-"`
	TotalUnread  int    `json:"-"` // TODO how does this apply to tags
}

func (c *Tag) String() string {
	return fmt.Sprintf("ID=%d, UserID=%d, Title=%s", c.ID, c.UserID, c.Title)
}

// TagRequest represents the request to create or update a Tag.
type TagRequest struct {
	Title        string `json:"title"`
	HideGlobally string `json:"hide_globally"`
}

// Patch updates Tag fields.
func (cr *TagRequest) Patch(tag *Tag) {
	tag.Title = cr.Title
	tag.HideGlobally = cr.HideGlobally != ""
}

// Categories represents a list of categories.
type Tags []*Tag
