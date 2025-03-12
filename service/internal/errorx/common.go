package errorx

const serviceInternal errCode = 5000

const (
	// api relative
	paramKeyRequired errCode = 4000 + iota
	passwordIncorrect
	invalidEmailFormat
	invalidPasswordFormat
)

const (
	// database relative
	noRecord errCode = 4100 + iota
	recordExisted
)

var (
	ParamKeyRequired      = newErrorResponse(paramKeyRequired, "some parameters or body's field are missing")
	NoRecord              = newErrorResponse(noRecord, "there are no records that match the specified query parameters")
	RecordExisted         = newErrorResponse(recordExisted, "the record already existed")
	PasswordIncorrect     = newErrorResponse(passwordIncorrect, "the password is incorrect")
	InvalidEmailFormat    = newErrorResponse(invalidEmailFormat, "invalid email format")
	InvalidPasswordFormat = newErrorResponse(invalidPasswordFormat, "invalid password format")
)
