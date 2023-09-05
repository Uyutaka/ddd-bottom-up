package main

import "fmt"

type (
	UserName struct {
		V string
	}

	UserId struct {
		V string
	}

	CircleId struct {
		V string
	}

	CircleName struct {
		V string
	}

	// Aggregate Root
	User struct {
		Id   UserId
		Name UserName
	}

	// Aggregate Root
	Circle struct {
		id      CircleId
		name    CircleName
		user    User
		members []User
	}

	ICircleRepository interface {
		Save(circle Circle) error
		FindById(id CircleId) (Circle, error)
		FindByName(name CircleName) (Circle, error)
	}

	ICircleFactory interface {
		Create(name CircleName, owner User) (Circle, error)
	}
	IUserRepository interface {
		Save(user User) error
		Find(id UserId) (User, error)
		Exists(user User) (bool, error)
	}
	CircleService struct {
		repo ICircleRepository
	}

	CircleCreateCommand struct {
		userId string
		name   string
	}

	CircleApplicationService struct {
		circleFactory    ICircleFactory
		circleRepository ICircleRepository
		circleService    CircleService
		userRepository   IUserRepository
	}

	CircleJoinCommand struct {
		userId   string
		circleId string
	}
)

func NewUserName(v string) (UserName, bool) {
	if len(v) == 0 {
		return UserName{}, false
	}
	if len(v) < 3 {
		return UserName{}, false
	}

	return UserName{V: v}, true
}

func NewUserId(v string) (UserId, bool) {
	if len(v) == 0 {
		return UserId{}, false
	}

	return UserId{V: v}, true
}

func NewCircleId(v string) (CircleId, bool) {
	if len(v) == 0 {
		return CircleId{}, false
	}
	return CircleId{V: v}, true
}

func NewCircleName(v string) (CircleName, bool) {
	if len(v) == 0 {
		return CircleName{}, false
	}
	if len(v) < 3 {
		return CircleName{}, false
	}
	if len(v) > 20 {
		return CircleName{}, false
	}

	return CircleName{V: v}, true
}

func NewUser(id UserId, name UserName) (User, bool) {
	if len(id.V) == 0 {
		return User{}, false
	}

	return User{Id: id}, true
}

func (u *User) ChangeUserName(name string) bool {

	if len(name) == 0 {
		return false
	}
	if len(name) < 3 {
		return false
	}

	userName, _ := NewUserName(name)
	u.Name = userName
	return true
}

func (c *CircleName) Equals(other CircleName) bool {
	return c.V == other.V
}

func NewCircleService(repo ICircleRepository) CircleService {
	return CircleService{repo: repo}
}

func (s *CircleService) Exist(circle Circle) bool {
	duplicated, _ := s.repo.FindByName(circle.name)
	return duplicated.name.V != ""
}

func NewCircle(id CircleId, name CircleName, user User, users []User) (Circle, bool) {
	if id.V == "" {
		return Circle{}, false
	}
	if name.V == "" {
		return Circle{}, false
	}
	if user.Id.V == "" {
		return Circle{}, false
	}

	if len(users) == 0 {
		return Circle{}, false
	}

	return Circle{
		id:      id,
		name:    name,
		user:    user,
		members: users,
	}, true
}

func NewCircleCreateCommand(userId string, userName string) CircleCreateCommand {
	return CircleCreateCommand{userId: userId, name: userName}
}

func NewCircleApplicationService(circleFactory ICircleFactory, circleRepository ICircleRepository, circleService CircleService, userRepository IUserRepository) CircleApplicationService {
	return CircleApplicationService{
		circleFactory:    circleFactory,
		circleRepository: circleRepository,
		circleService:    circleService,
		userRepository:   userRepository,
	}
}

func (cas *CircleApplicationService) Create(command CircleCreateCommand) bool {

	// TX Starts

	// find owner's user id
	ownerId, _ := NewUserId(command.userId)
	owner, err := cas.userRepository.Find(ownerId)
	if err != nil {
		return false
	}

	name, _ := NewCircleName(command.name)
	circle, _ := cas.circleFactory.Create(name, owner)

	// check duplication
	if cas.circleService.Exist(circle) {
		return false
	}

	cas.circleRepository.Save(circle)
	return true
	// TX Ends
}

func (cas *CircleApplicationService) Join(command CircleJoinCommand) bool {
	// TX Starts

	memberId, _ := NewUserId(command.userId)

	member, err := cas.userRepository.Find(memberId)

	if err != nil {
		return false
	}

	circleId, _ := NewCircleId(command.circleId)
	circle, err := cas.circleRepository.FindById(circleId)
	if err != nil {
		return false
	}

	if len(circle.members) >= 29 {
		return false
	}

	// This violates Law of Demeter (See List 12.2 & Chap 12.1.2)
	circle.members = append(circle.members, member)

	cas.circleRepository.Save(circle)
	return true
	// TX Ends

}

func (c *Circle) Join(member User) bool {
	if len(member.Name.V) == 0 {
		return false
	}

	if c.IsFull() {
		return false
	}

	c.members = append(c.members, member)
	return true
}

func (c *Circle) IsFull() bool {
	return len(c.members) >= 29
}

func NewCircleJoinCommand(userId string, circleId string) CircleJoinCommand {
	return CircleJoinCommand{userId: userId, circleId: circleId}
}

func main() {
	var u2 User
	var c1 Circle
	fmt.Println(u2)
	fmt.Println(c1)

}
