// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package model // import "miniflux.app/model"

import "fmt"

// Tag represents a feed tag.
type Tag struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	UserID int64  `json:"user_id"`
}

func (c *Tag) String() string {
	return fmt.Sprintf("ID=%d, UserID=%d, Title=%s", c.ID, c.UserID, c.Title)
}

// TagRequest represents the request to create or update a tag.
type TagRequest struct {
	Title string `json:"title"`
}

// Patch updates tag fields.
func (cr *TagRequest) Patch(tag *Tag) {
	tag.Title = cr.Title
}

// Tags represents a list of tags.
type Tags []*Tag
