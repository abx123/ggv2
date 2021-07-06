package constant

import "fmt"

var ErrTableNotFound = fmt.Errorf("table not found")

var ErrGuestNotFound = fmt.Errorf("guest not found")

var ErrGuestAlreadyRSVP = fmt.Errorf("guest already RSVP")

var ErrGuestAlreadyArrived = fmt.Errorf("guest already arrived")

var ErrTableIsFull = fmt.Errorf("table is full")

var ErrFailedOptimisticLock = fmt.Errorf("unable to secure optimistic lock, please retry")

var ErrCapacityLessThanOne = fmt.Errorf("capacity cannot be less than 1")

var ErrInvalidRequest = fmt.Errorf("invalid request parameter")

var ErrDBErr = fmt.Errorf("database returns error")

var ErrAccompanyingGuestLessThanZero = fmt.Errorf("accompanying guest cannot be less than 0")

var ErrGuestNeverRSVP = fmt.Errorf("guest never rsvp")

var ErrGuestNotArrived = fmt.Errorf("guest not arrived")
