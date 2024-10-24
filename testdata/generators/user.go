package generators

import "auth/internal/domain"

func RandomUser() (result domain.User) {
	return domain.User{
		Uuid:     RandomUuidString(),
		Username: RandomString(255),
		Password: RandomString(255),
	}
}
