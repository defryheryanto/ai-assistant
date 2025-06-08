package user

type CreateUserParams struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Role  Role   `json:"role"`
	Email string `json:"email"`
}
