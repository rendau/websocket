package rest

type ErrRepSt struct {
	ErrorCode string `json:"error_code"`
}

type AuthRepSt struct {
	Id        int64  `json:"id"`
	ErrorCode string `json:"error_code"`
	ErrorDesc string `json:"error_desc"`
}
