package presenter

type Error struct {
	ReqId  string `json:"requestID"`
	ErrMsg string `json:"message"`
}

func ErrResp(reqID string, err error) *Error {
	return &Error{
		ReqId:  reqID,
		ErrMsg: err.Error(),
	}
}
