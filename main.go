package main

import "fmt"

type (
	CircleId struct {
		V string
	}

	CircleName struct {
		V string
	}

	User struct {
		V string
	}

	Circle struct {
		id    CircleId
		name  CircleName
		user  User
		users []User
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
)

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

func NewUser(v string) (User, bool) {
	if len(v) == 0 {
		return User{}, false
	}
	if len(v) < 3 {
		return User{}, false
	}

	return User{V: v}, true
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
	if user.V == "" {
		return Circle{}, false
	}

	if len(users) == 0 {
		return Circle{}, false
	}

	return Circle{
		id:    id,
		name:  name,
		user:  user,
		users: users,
	}, true
}

func main() {
	u1, _ := NewUser("user1")
	var u2 User
	var c1 Circle
	fmt.Println(u1)
	fmt.Println(u2)
	fmt.Println(c1)

}
