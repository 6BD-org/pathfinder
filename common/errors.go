package common

import (
	"fmt"

	"github.com/6BD-org/pathfinder/consts"
)

type PathFinderError struct {
	error
	ErrCode consts.ErrCode
	Msg     string
}

func (pfe PathFinderError) Error() string {
	return fmt.Sprintf("[%v] PathFinderError: %v", pfe.ErrCode, pfe.Msg)
}

// NewErr PathFinder error
func NewErr(errCode consts.ErrCode, msg string, args ...interface{}) PathFinderError {
	return PathFinderError{
		ErrCode: errCode,
		Msg:     fmt.Sprintf(msg, args...),
	}
}
