package model

type (
	UserService struct {
		userRepository IUserRepository
	}
)

func NewUserService(userRepository IUserRepository) UserService {
	return UserService{userRepository: userRepository}
}

func (us *UserService) Exists(user *User) bool {
	duplicatedUser, _ := us.userRepository.FindByName(&user.Name)
	return duplicatedUser != nil
}
