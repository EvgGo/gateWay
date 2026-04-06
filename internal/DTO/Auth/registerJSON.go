package DTO

// RegisterJSON - DTO для HTTP JSON входа.
// Мы отделяем HTTP-формат от protobuf, чтобы:
//
// - контролировать json-теги (snake_case),
// - использовать *string для optional полей там, где нужно,
// - проще вводить legacy/aliased поля в будущем.
type RegisterJSON struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}
