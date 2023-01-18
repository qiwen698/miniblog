package errno

//定义一些预定义的错误类型，供程序直接引用

var (
	// OK 代表请求成功
	OK = &Errno{HTTP: 200, Code: "", Message: ""}
	// InternalServerError 表示所有未知的服务端错误.
	InternalServerError = &Errno{HTTP: 500, Code: "InternalError", Message: "Internal server error."}
	// ErrPageNotFound 表示路由不匹配错误.
	ErrPageNotFound = &Errno{HTTP: 404, Code: "ResourceNotFound.PageNotFound", Message: "Page not found"}
)
