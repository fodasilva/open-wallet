package users

type CreateUserInput struct {
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	AvatarURL *string `json:"avatar_url"`
	Username  string  `json:"username"`
}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	CreatedAt string `json:"created_at"`
	Username  string `json:"username"`
}

type UserFilter struct {
	Email    string
	Username string
}
