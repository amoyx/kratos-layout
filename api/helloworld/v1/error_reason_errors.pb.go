// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package v1

import (
	fmt "fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

func IsCreateUserFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_CREATE_USER_FAILED.String() && e.Code == 201
}

func ErrorCreateUserFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(201, ErrorReason_CREATE_USER_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsUpdateUserFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_UPDATE_USER_FAILED.String() && e.Code == 202
}

func ErrorUpdateUserFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(202, ErrorReason_UPDATE_USER_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsDeleteUserFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_DELETE_USER_FAILED.String() && e.Code == 203
}

func ErrorDeleteUserFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(203, ErrorReason_DELETE_USER_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsUserNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == ErrorReason_USER_NOT_FOUND.String() && e.Code == 204
}

func ErrorUserNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(204, ErrorReason_USER_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}