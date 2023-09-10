package model

import (
	"time"
)

type (
	CircleId struct {
		V string
	}

	CircleName struct {
		V string
	}
	// Aggregate Root
	Circle struct {
		id      *CircleId
		name    *CircleName
		owner   *UserId
		members []UserId
		created time.Time
	}

	ICircleRepository interface {
		Save(circle *Circle) error
		FindById(id CircleId) (*Circle, error)
		FindByName(name *CircleName) (Circle, error)
		// ng because condition of searching circles is not in repository of domain model
		// FindRecommended(time time.Time) ([]Circle, error)
		FindAll() ([]Circle, error)
	}

	ICircleFactory interface {
		Create(name *CircleName, owner *User) (*Circle, error)
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
		now              time.Time
	}

	CircleJoinCommand struct {
		userId   string
		circleId string
	}

	CircleFullSpecification struct {
		repo IUserRepository
	}

	CircleGetRecommendResult struct {
		circles []Circle
	}

	CircleRecommendSpecification struct {
		executeDateTime time.Time
	}
)

func (c *CircleName) Equals(other CircleName) bool {
	return c.V == other.V
}

func NewCircleService(repo ICircleRepository) CircleService {
	return CircleService{repo: repo}
}

func (s *CircleService) Exist(circle *Circle) bool {
	duplicated, _ := s.repo.FindByName(circle.name)
	return duplicated.name.V != ""
}

func NewCircle(id *CircleId, name *CircleName, owner *UserId, users []UserId) (Circle, bool) {
	if id == nil {
		return Circle{}, false
	}
	if name == nil {
		return Circle{}, false
	}
	if owner == nil {
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

func NewCircleApplicationService(circleFactory ICircleFactory, circleRepository ICircleRepository, circleService CircleService, userRepository IUserRepository, now time.Time) CircleApplicationService {
	return CircleApplicationService{
		circleFactory:    circleFactory,
		circleRepository: circleRepository,
		circleService:    circleService,
		userRepository:   userRepository,
		now:              time.Now(),
	}
}

func (cas *CircleApplicationService) Create(command CircleCreateCommand) bool {

	// TX Starts

	// find owner's user id
	ownerId, _ := NewUserId(command.userId)
	owner, err := cas.userRepository.Find(&ownerId)
	if err != nil {
		return false
	}

	name, _ := NewCircleName(command.name)
	circle, _ := cas.circleFactory.Create(&name, owner)

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

	user, err := cas.userRepository.Find(&memberId)

	if err != nil {
		return false
	}

	circleId, _ := NewCircleId(command.circleId)
	circle, err := cas.circleRepository.FindById(circleId)
	if err != nil {
		return false
	}

	cfs := NewCircleFullSpecification(cas.userRepository)
	if cfs.IsSatisfiedBy(circle) {
		return false
	}

	// This violates Law of Demeter (See List 12.2 & Chap 12.1.2)
	// circle.members = append(circle.members, memberId)
	if !circle.Join(user) {
		return false
	}

	cas.circleRepository.Save(circle)
	return true
	// TX Ends

}

func (cas *CircleApplicationService) GetRecommend() CircleGetRecommendResult {
	recommendCircleSpec := NewCircleRecommendSpecification(cas.now)

	circles, _ := cas.circleRepository.FindAll()
	recommendCircles := []Circle{}
	for _, circle := range circles {
		if recommendCircleSpec.IsSatisfiedBy(circle) {
			recommendCircles = append(recommendCircles, circle)
		}
		if len(recommendCircles) >= 10 {
			break
		}
	}
	return CircleGetRecommendResult{circles: recommendCircles}
}

func (c *Circle) Join(member *User) bool {
	if member == nil {
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

func NewCircleFullSpecification(repo IUserRepository) CircleFullSpecification {
	return CircleFullSpecification{repo: repo}
}

func NewCircleRecommendSpecification(executeDateTime time.Time) CircleRecommendSpecification {
	return CircleRecommendSpecification{executeDateTime: executeDateTime}
}

func (cfs *CircleFullSpecification) IsSatisfiedBy(circle *Circle) bool {
	owner, _ := cfs.repo.Find(circle.owner)
	upperLimit := 30
	if owner.IsPremium() {
		upperLimit = 50
	}

	return circle.CountMembers() >= upperLimit
}

func (crs *CircleRecommendSpecification) IsSatisfiedBy(circle Circle) bool {
	if circle.CountMembers() < 10 {
		return false
	}

	return circle.created.Before(crs.executeDateTime.AddDate(0, -1, 0))
}
