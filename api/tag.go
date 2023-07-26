// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package api // import "miniflux.app/api"

import (
	json_parser "encoding/json"
	"fmt"
	"net/http"
	"time"

	"miniflux.app/http/request"
	"miniflux.app/http/response/json"
	"miniflux.app/model"
	"miniflux.app/validator"
)

func (h *handler) createTag(w http.ResponseWriter, r *http.Request) {
	userID := request.UserID(r)

	var tagRequest model.TagRequest
	if err := json_parser.NewDecoder(r.Body).Decode(&tagRequest); err != nil {
		json.BadRequest(w, r, err)
		return
	}

	if validationErr := validator.ValidateTagCreation(h.store, userID, &tagRequest); validationErr != nil {
		json.BadRequest(w, r, validationErr.Error())
		return
	}

	tag, err := h.store.CreateTag(userID, &tagRequest)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	json.Created(w, r, tag)
}

func (h *handler) getTags(w http.ResponseWriter, r *http.Request) {
	var tags model.Tags
	var err error
	// includeCounts := request.QueryStringParam(r, "counts", "false")

	// if includeCounts == "true" {
	// 	tags, err = h.store.TagsWithFeedCount(request.UserID(r))
	// } else {
	tags, err = h.store.Tags(request.UserID(r))
	// }

	if err != nil {
		json.ServerError(w, r, err)
		return
	}
	json.OK(w, r, tags)
}

func (h *handler) updateTag(w http.ResponseWriter, r *http.Request) {
	fmt.Println(`doggy`)
	userID := request.UserID(r)
	tagID := request.RouteInt64Param(r, "tagID")

	tag, err := h.store.Tag(userID, tagID)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	if tag == nil {
		json.NotFound(w, r)
		return
	}

	var tagRequest model.TagRequest
	if err := json_parser.NewDecoder(r.Body).Decode(&tagRequest); err != nil {
		json.BadRequest(w, r, err)
		return
	}

	if validationErr := validator.ValidateTagModification(h.store, userID, tag.ID, &tagRequest); validationErr != nil {
		json.BadRequest(w, r, validationErr.Error())
		return
	}

	tagRequest.Patch(tag)
	err = h.store.UpdateTag(tag)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	json.Created(w, r, tag)
}

func (h *handler) removeTag(w http.ResponseWriter, r *http.Request) {
	userID := request.UserID(r)
	tagID := request.RouteInt64Param(r, "tagID")

	if !h.store.TagIDExists(userID, tagID) {
		json.NotFound(w, r)
		return
	}

	if err := h.store.RemoveTag(userID, tagID); err != nil {
		json.ServerError(w, r, err)
		return
	}

	json.NoContent(w, r)
}

// get all feeds that are associated with a tag
func (h *handler) getTagFeeds(w http.ResponseWriter, r *http.Request) {
	userID := request.UserID(r)
	tagID := request.RouteInt64Param(r, "tagID")

	tag, err := h.store.Tag(userID, tagID)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	if tag == nil {
		json.NotFound(w, r)
		return
	}

	feeds, err := h.store.FeedsByTagWithCounters(userID, tagID)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	json.OK(w, r, feeds)
}

func (h *handler) markTagAsRead(w http.ResponseWriter, r *http.Request) {
	userID := request.UserID(r)
	tagID := request.RouteInt64Param(r, "tagID")

	tag, err := h.store.Tag(userID, tagID)
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	if tag == nil {
		json.NotFound(w, r)
		return
	}

	if err = h.store.MarkTagAsRead(userID, tagID, time.Now()); err != nil {
		json.ServerError(w, r, err)
		return
	}

	json.NoContent(w, r)
}

/*
func (h *handler) getTags(w http.ResponseWriter, r *http.Request) {
	var tags model.Tags
	var err error
	includeCounts := request.QueryStringParam(r, "counts", "false")

	if includeCounts == "true" {
		tags, err = h.store.TagsWithFeedCount(request.UserID(r))
	} else {
		tags, err = h.store.Tags(request.UserID(r))
	}

	if err != nil {
		json.ServerError(w, r, err)
		return
	}
	json.OK(w, r, tags)
}

func (h *handler) refreshTag(w http.ResponseWriter, r *http.Request) {
	userID := request.UserID(r)
	tagID := request.RouteInt64Param(r, "tagID")

	jobs, err := h.store.NewTagBatch(userID, tagID, h.store.CountFeeds(userID))
	if err != nil {
		json.ServerError(w, r, err)
		return
	}

	go func() {
		h.pool.Push(jobs)
	}()

	json.NoContent(w, r)
}
*/
