// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package storage // import "miniflux.app/storage"

import (
	"database/sql"
	"errors"
	"fmt"

	// "github.com/lib/pq"
	"miniflux.app/model"
)

// CreateTag creates a new tag.
func (s *Storage) CreateTag(userID int64, request *model.TagRequest) (*model.Tag, error) {
	var tag model.Tag

	query := `
		INSERT INTO tags
			(user_id, title)
		VALUES
			($1, $2)
		RETURNING
			id,
			user_id,
			title
	`
	err := s.db.QueryRow(
		query,
		userID,
		request.Title,
	).Scan(
		&tag.ID,
		&tag.UserID,
		&tag.Title,
	)
	if err != nil {
		return nil, fmt.Errorf(`store: unable to create tag %q: %v`, request.Title, err)
	}

	return &tag, nil
}

// AnotherTagExists checks if another tag exists with the same title.
func (s *Storage) AnotherTagExists(userID, tagID int64, title string) bool {
	var result bool
	query := `SELECT true FROM tags WHERE user_id=$1 AND id != $2 AND lower(title)=lower($3) LIMIT 1`
	s.db.QueryRow(query, userID, tagID, title).Scan(&result)
	return result
}

// TagTitleExists checks if the given tag exists in the database.
func (s *Storage) TagTitleExists(userID int64, title string) bool {
	var result bool
	query := `SELECT true FROM tags WHERE user_id=$1 AND lower(title)=lower($2) LIMIT 1`
	s.db.QueryRow(query, userID, title).Scan(&result)
	return result
}

// Tags returns all tags that belongs to the given user.
func (s *Storage) Tags(userID int64) (model.Tags, error) {
	query := `SELECT id, user_id, title FROM tags WHERE user_id=$1 ORDER BY title ASC`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf(`store: unable to fetch tags: %v`, err)
	}
	defer rows.Close()

	tags := make(model.Tags, 0)
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.ID, &tag.UserID, &tag.Title); err != nil {
			return nil, fmt.Errorf(`store: unable to fetch tag row: %v`, err)
		}

		tags = append(tags, &tag)
	}

	return tags, nil
}

// Tag returns a tag from the database.
func (s *Storage) Tag(userID, tagID int64) (*model.Tag, error) {
	var tag model.Tag

	query := `SELECT id, user_id, title FROM tags WHERE user_id=$1 AND id=$2`
	err := s.db.QueryRow(query, userID, tagID).Scan(&tag.ID, &tag.UserID, &tag.Title)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf(`store: unable to fetch tag: %v`, err)
	default:
		return &tag, nil
	}
}

// UpdateTag updates an existing tag.
func (s *Storage) UpdateTag(tag *model.Tag) error {
	query := `UPDATE tags SET title=$1 WHERE id=$2 AND user_id=$3`
	_, err := s.db.Exec(
		query,
		tag.Title,
		tag.ID,
		tag.UserID,
	)
	if err != nil {
		return fmt.Errorf(`store: unable to update tag: %v`, err)
	}

	return nil
}

// RemoveTag deletes a tag.
func (s *Storage) RemoveTag(userID, tagID int64) error {
	query := `DELETE FROM tags WHERE id = $1 AND user_id = $2`
	result, err := s.db.Exec(query, tagID, userID)
	if err != nil {
		return fmt.Errorf(`store: unable to remove this tag: %v`, err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(`store: unable to remove this tag: %v`, err)
	}

	if count == 0 {
		return errors.New(`store: no tag has been removed`)
	}

	return nil
}

// TagIDExists checks if the given tag exists into the database.
func (s *Storage) TagIDExists(userID, tagID int64) bool {
	var result bool
	query := `SELECT true FROM tags WHERE user_id=$1 AND id=$2`
	s.db.QueryRow(query, userID, tagID).Scan(&result)
	return result
}

// TagByTitle finds a tag by the title.
func (s *Storage) TagByTitle(userID int64, title string) (*model.Tag, error) {
	var tag model.Tag

	query := `SELECT id, user_id, title FROM tags WHERE user_id=$1 AND title=$2`
	err := s.db.QueryRow(query, userID, title).Scan(&tag.ID, &tag.UserID, &tag.Title)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf(`store: unable to fetch tag: %v`, err)
	default:
		return &tag, nil
	}
}

/*
// FirstTag returns the first tag for the given user.
func (s *Storage) FirstTag(userID int64) (*model.Tag, error) {
	query := `SELECT id, user_id, title, hide_globally FROM tags WHERE user_id=$1 ORDER BY title ASC LIMIT 1`

	var tag model.Tag
	err := s.db.QueryRow(query, userID).Scan(&tag.ID, &tag.UserID, &tag.Title, &tag.HideGlobally)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf(`store: unable to fetch tag: %v`, err)
	default:
		return &tag, nil
	}
}

// TagsWithFeedCount returns all tags with the number of feeds.
func (s *Storage) TagsWithFeedCount(userID int64) (model.Tags, error) {
	user, err := s.UserByID(userID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			c.id,
			c.user_id,
			c.title,
			c.hide_globally,
			(SELECT count(*) FROM feeds WHERE feeds.tag_id=c.id) AS count,
			(SELECT count(*)
			   FROM feeds
			     JOIN entries ON (feeds.id = entries.feed_id)
			   WHERE feeds.tag_id = c.id AND entries.status = 'unread') AS count_unread
		FROM tags c
		WHERE
			user_id=$1
	`

	if user.TagsSortingOrder == "alphabetical" {
		query = query + `
			ORDER BY
				c.title ASC
		`
	} else {
		query = query + `
			ORDER BY
				count_unread DESC,
				c.title ASC
		`
	}

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf(`store: unable to fetch tags: %v`, err)
	}
	defer rows.Close()

	tags := make(model.Tags, 0)
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.ID, &tag.UserID, &tag.Title, &tag.HideGlobally, &tag.FeedCount, &tag.TotalUnread); err != nil {
			return nil, fmt.Errorf(`store: unable to fetch tag row: %v`, err)
		}

		tags = append(tags, &tag)
	}

	return tags, nil
}

// delete the given tags, replacing those tags with the user's first
// tag on affected feeds
func (s *Storage) RemoveAndReplaceTagsByName(userid int64, titles []string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return errors.New("unable to begin transaction")
	}

	titleParam := pq.Array(titles)
	var count int
	query := "SELECT count(*) FROM tags WHERE user_id = $1 and title != ANY($2)"
	err = tx.QueryRow(query, userid, titleParam).Scan(&count)
	if err != nil {
		tx.Rollback()
		return errors.New("unable to retrieve tag count")
	}
	if count < 1 {
		tx.Rollback()
		return errors.New("at least 1 tag must remain after deletion")
	}

	query = `
		WITH d_cats AS (SELECT id FROM tags WHERE user_id = $1 AND title = ANY($2))
		UPDATE feeds
		 SET tag_id =
		  (SELECT id
			FROM tags
			WHERE user_id = $1 AND id NOT IN (SELECT id FROM d_cats)
			ORDER BY title ASC
			LIMIT 1)
		WHERE user_id = $1 AND tag_id IN (SELECT id FROM d_cats)
	`
	_, err = tx.Exec(query, userid, titleParam)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("unable to replace tags: %v", err)
	}

	query = "DELETE FROM tags WHERE user_id = $1 AND title = ANY($2)"
	_, err = tx.Exec(query, userid, titleParam)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("unable to delete tags: %v", err)
	}
	tx.Commit()
	return nil
}
*/
