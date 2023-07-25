// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package validator // import "miniflux.app/validator"

import (
	"miniflux.app/model"
	"miniflux.app/storage"
)

// ValidateTagCreation validates tag creation.
func ValidateTagCreation(store *storage.Storage, userID int64, request *model.TagRequest) *ValidationError {
	if request.Title == "" {
		return NewValidationError("error.title_required")
	}

	if store.TagTitleExists(userID, request.Title) {
		return NewValidationError("error.tag_already_exists")
	}

	return nil
}

// ValidateTagModification validates tag modification.
func ValidateTagModification(store *storage.Storage, userID, tagID int64, request *model.TagRequest) *ValidationError {
	if request.Title == "" {
		return NewValidationError("error.title_required")
	}

	if store.AnotherTagExists(userID, tagID, request.Title) {
		return NewValidationError("error.tag_already_exists")
	}

	return nil
}
