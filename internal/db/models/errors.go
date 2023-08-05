package models

import "errors"

var (
	// Common errors
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")

	// User signup errors
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")

	// File upload errors
	ErrUploadFailed  = errors.New("file upload has failed")
	ErrGetFileFailed = errors.New("file delete has failed")

	// Follow errors
	ErrAlreadyFollowed      = errors.New("already followed")
	ErrNotFollowed          = errors.New("not followed")
	ErrCannotFollowYourself = errors.New("cannot follow yourself")

	// Post errors
	ErrAlreadyLiked = errors.New("already liked")
	ErrNotLiked     = errors.New("not liked")
)
