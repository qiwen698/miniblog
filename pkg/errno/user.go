package errno

var (
	ErrUserAlreadyExist = &Errno{HTTP: 400, Code: "FailedOperation.UserAlreadyExist", Message: "User already exist."}
)
