package model

import (
	"errors"
	"fmt"
)

var (
	ErrorOrderExistAnotherUser = errors.New("order already exist for another user")
	ErrorOrderExist            = errors.New("order already exist for that user")
	ErrorConflictSaveUser      = errors.New("user already exist")
	ErrorItemNotFound          = errors.New("item not found")
	ErrorNotEnoughFounds       = errors.New("not enough founds ")
	ErrorWithdrawExist         = errors.New("withdraw already exist for that user")
)

var _ error = (*ErrorAccrualLimitAchieved)(nil)

type ErrorAccrualLimitAchieved struct {
	Err         error
	WaitSeconds int
	PerMinute   int
}

func (e *ErrorAccrualLimitAchieved) Error() string {
	return fmt.Sprintf("To many requests, pause %v seconds and make %v requests per minute", e.WaitSeconds, e.PerMinute)
}

func (e *ErrorAccrualLimitAchieved) Is(tgt error) bool {
	_, ok := tgt.(*ErrorAccrualLimitAchieved)
	return ok
}
