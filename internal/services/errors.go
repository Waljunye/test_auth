package services

type ErrSignIn struct{}

func (e ErrSignIn) Error() string {
	return "username or password is incorrect"
}
