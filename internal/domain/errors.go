package domain

import "errors"

var (
	ErrFriendshipRequestExists    = errors.New("friendship request exists")
	ErrFriendshipRequestNotExists = errors.New("friendship request doesn't exist")
	ErrWrongCredentials           = errors.New("wrong credentials")
	ErrBadRequest                 = errors.New("bad request")
)
