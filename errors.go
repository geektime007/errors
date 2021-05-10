package errors

//  错误列表

// * 只有内部用到的错误，错误走 ErrorCodeInternal， 成功 ErrorCodeNil
// * 需要提供给前端的错误，走具体的 pb 里面的错误码

var (
	// 通用错误
	Nil             *Error = nil
	Unknown                = NewError(ErrorCodeUnknown, 500, "未知错误")
	BadFormat              = NewError(ErrorCodeBadFormat, 400, "格式错误")
	RunCommandError        = NewError(ErrorRunCommandError, 500, "执行命令异常")
	Timeout                = NewError(ErrorCodeTimeout, 408, "超时")
	Unimplemented          = NewError(ErrorCodeUnimplemented, 404, "未实现!")
)
