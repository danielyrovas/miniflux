// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package storage // import "miniflux.app/storage"

import (
	"database/sql"
	"errors"
	"fmt"

	"miniflux.app/model"
)

// AnotherCategoryExists checks if another category exists with the same title.
func (s *Storage) AnotherCategoryExists(userID, categoryID int64, title string) bool {
	var result bool
	query := `SELECT true FROM categories WHERE user_id=$1 AND id != $2 AND title=$3`
	s.db.QueryRow(query, userID, categoryID, title).Scan(&result)
	return result
}

// CategoryExists checks if the given category exists into the database.
func (s *Storage) CategoryExists(userID, categoryID int64) bool {
	var result bool
	query := `SELECT true FROM categories WHERE user_id=$1 AND id=$2`
	s.db.QueryRow(query, userID, categoryID).Scan(&result)
	return result
}

// Category returns a category from the database.
func (s *Storage) Category(userID, categoryID int64) (*model.Category, error) {
	var category model.Category

	query := `SELECT id, user_id, title, view FROM categories WHERE user_id=$1 AND id=$2`
	err := s.db.QueryRow(query, userID, categoryID).Scan(&category.ID, &category.UserID, &category.Title, &category.View)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf(`store: unable to fetch category: %v`, err)
	default:
		return &category, nil
	}
}

// FirstCategory returns the first category for the given user.
func (s *Storage) FirstCategory(userID int64) (*model.Category, error) {
	query := `SELECT id, user_id, title FROM categories WHERE user_id=$1 ORDER BY title ASC LIMIT 1`

	var category model.Category
	err := s.db.QueryRow(query, userID).Scan(&category.ID, &category.UserID, &category.Title)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf(`store: unable to fetch category: %v`, err)
	default:
		return &category, nil
	}
}

// CategoryByTitle finds a category by the title.
func (s *Storage) CategoryByTitle(userID int64, title string) (*model.Category, error) {
	var category model.Category

	query := `SELECT id, user_id, title FROM categories WHERE user_id=$1 AND title=$2`
	err := s.db.QueryRow(query, userID, title).Scan(&category.ID, &category.UserID, &category.Title)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf(`store: unable to fetch category: %v`, err)
	default:
		return &category, nil
	}
}

// Categories returns all categories that belongs to the given user.
func (s *Storage) Categories(userID int64) (model.Categories, error) {
	query := `SELECT id, user_id, title FROM categories WHERE user_id=$1 ORDER BY title ASC`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf(`store: unable to fetch categories: %v`, err)
	}
	defer rows.Close()

	categories := make(model.Categories, 0)
	for rows.Next() {
		var category model.Category
		if err := rows.Scan(&category.ID, &category.UserID, &category.Title); err != nil {
			return nil, fmt.Errorf(`store: unable to fetch category row: %v`, err)
		}

		categories = append(categories, &category)
	}

	return categories, nil
}

// CategoriesWithFeedCount returns all categories with the number of feeds.
func (s *Storage) CategoriesWithFeedCount(userID int64) (model.Categories, error) {
	query := `
		SELECT
			c.id,
			c.user_id,
			c.title,
			(SELECT count(*) FROM feeds WHERE feeds.category_id=c.id) AS count
		FROM categories c
		WHERE
			user_id=$1
		ORDER BY c.title ASC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf(`store: unable to fetch categories: %v`, err)
	}
	defer rows.Close()

	categories := make(model.Categories, 0)
	for rows.Next() {
		var category model.Category
		if err := rows.Scan(&category.ID, &category.UserID, &category.Title, &category.FeedCount); err != nil {
			return nil, fmt.Errorf(`store: unable to fetch category row: %v`, err)
		}

		categories = append(categories, &category)
	}

	return categories, nil
}

// CreateCategory creates a new category.
func (s *Storage) CreateCategory(category *model.Category) error {
	query := `
		INSERT INTO categories
			(user_id, title, view)
		VALUES
			($1, $2, $3)
		RETURNING 
			id
	`
	err := s.db.QueryRow(
		query,
		category.UserID,
		category.Title,
		category.View,
	).Scan(&category.ID)

	if err != nil {
		return fmt.Errorf(`store: unable to create category: %v`, err)
	}

	return nil
}

// UpdateCategory updates an existing category.
func (s *Storage) UpdateCategory(category *model.Category) error {
	query := `UPDATE categories SET title=$1, view=$2 WHERE id=$3 AND user_id=$4`
	_, err := s.db.Exec(
		query,
		category.Title,
		category.View,
		category.ID,
		category.UserID,
	)

	if err != nil {
		return fmt.Errorf(`store: unable to update category: %v`, err)
	}

	return nil
}

// RemoveCategory deletes a category.
func (s *Storage) RemoveCategory(userID, categoryID int64) error {
	query := `DELETE FROM categories WHERE id = $1 AND user_id = $2`
	result, err := s.db.Exec(query, categoryID, userID)
	if err != nil {
		return fmt.Errorf(`store: unable to remove this category: %v`, err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(`store: unable to remove this category: %v`, err)
	}

	if count == 0 {
		return errors.New(`store: no category has been removed`)
	}

	return nil
}

// UpdateCategoryView updates a category view setting.
func (s *Storage) UpdateCategoryView(userID int64, categoryID int64, view string) error {
	if _, ok := model.Views()[view]; !ok {
		return fmt.Errorf("invalid view value: %v", view)
	}
	query := `UPDATE categories SET view = $3 WHERE user_id=$1 AND id=$2`
	result, err := s.db.Exec(query, userID, categoryID, view)
	if err != nil {
		return fmt.Errorf("unable to set view for category #%d: %v", categoryID, err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to set view for category #%d: %v", categoryID, err)
	}

	if count == 0 {
		return errors.New("nothing has been updated")
	}

	return nil
}
