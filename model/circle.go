package model

type (
	CircleId struct {
		V string
	}

	CircleName struct {
		V string
	}
	// Aggregate Root
	Circle struct {
		id      CircleId
		name    CircleName
		owner   UserId
		members []UserId
	}

	ICircleRepository interface {
		Save(circle Circle) error
		FindById(id CircleId) (Circle, error)
		FindByName(name CircleName) (Circle, error)
	}

	ICircleFactory interface {
		Create(name CircleName, owner User) (Circle, error)
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

func NewCircle(id CircleId, name CircleName, owner UserId, users []UserId) (Circle, bool) {
	if id.V == "" {
		return Circle{}, false
	}
	if name.V == "" {
		return Circle{}, false
	}
	if owner.V == "" {
		return Circle{}, false
	}

	if len(users) == 0 {
		return Circle{}, false
	}

	return Circle{
		id:      id,
		name:    name,
		owner:   owner,
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

	_, err := cas.userRepository.Find(memberId)

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
	circle.members = append(circle.members, memberId)

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

	c.members = append(c.members, member.Id)
	return true
}

func (c *Circle) IsFull() bool {
	return c.CountMembers() >= 30
}

func (c *Circle) CountMembers() int {
	return len(c.members) + 1
}

func NewCircleJoinCommand(userId string, circleId string) CircleJoinCommand {
	return CircleJoinCommand{userId: userId, circleId: circleId}
}
