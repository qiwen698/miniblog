package v1

// CreateUserRequest 指定了 `POST /v1/users` 接口的请求参数.
type CreateUserRequest struct {
	Username string `json:"username" valid:"alphanum,required,stringlength(1|255)"`
	Password string `json:"password" valid:"required,stringlength(6|18)"`
	Nickname string `json:"nickname" valid:"required,stringlength(1|255)"`
	Email    string `json:"email" valid:"required,email"`
	Phone    string `json:"phone" valid:"required,stringlength(11|11)"`
}

// GetUserResponse 指定了 `GET /v1/users/{name}` 接口的返回参数.
type GetUserResponse UserInfo

// UserInfo 指定了用户的详细信息.
type UserInfo struct {
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	PostCount int64  `json:"postCount"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// LoginRequest 指定了 `POST /login` 接口的请求参数

type LoginRequest struct {
	Username string `json:"username" valid:"alphanum,required,stringlength(1|255)"`
	Password string `json:"password" valid:"required,stringlength(6|18)"`
}

// LoginResponse 指定了 `POST /login` 接口的返回参数.

type LoginResponse struct {
	Token string `json:"token"`
}

type ChangePasswordRequest struct {
	// 旧密码.
	OldPassword string `json:"oldPassword" valid:"required,stringlength(6|18)"`
	// 新密码.
	NewPassword string `json:"newPassword" valid:"required,stringlength(6|18)"`
}

// UpdateUserRequest 指定了 `PUT /v1/users/{name}` 接口的请求参数.
type UpdateUserRequest struct {
	Nickname *string `json:"nickname" valid:"stringlength(1|255)"`
	Email    *string `json:"email" valid:"email"`
	Phone    *string `json:"phone" valid:"stringlength(11|11)"`
}
