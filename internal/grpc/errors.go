package grpc

import "errors"

var (
	ErrorURLNotFounded  = errors.New("url not founded")
	ErrorURLIsDeleted   = errors.New("url is deleted")
	ErrorURLIsExist     = errors.New("url is exist")
	ErrorWrongUserID    = errors.New("wrong user id")
	ErrorURLListIsEmpty = errors.New("url list is empty")
)
