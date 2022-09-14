package app

type TokenAuth interface {
	TokenGenerator
	TokenDecoder
}

type TokenGenerator interface {
	GenerateTokenString(userID string) (string, error)
}

type TokenDecoder interface {
	GetSubject(tokenString string) (string, error)
}
