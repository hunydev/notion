package notion

import "fmt"

const (
	UserTypePerson = "pserson"
	UserTypeBot    = "bot"
)

type User struct {
	ID string

	JSON JSON
}

func NewUser(ID string) *User {
	user := &User{
		ID: ID,
		JSON: JSON{
			"object": "user",
			"id":     ID,
		},
	}

	return user
}

func (user *User) Object() string {
	return "user"
}

func (user *User) Type() (string, error) {
	t := user.JSON.GetString("type")
	switch t {
	case UserTypePerson:
		return UserTypePerson, nil
	case UserTypeBot:
		return UserTypeBot, nil
	}

	return "", fmt.Errorf("Unknown UserType[%s]", t)
}

func (user *User) IsPerson() bool {
	t, _ := user.Type()
	return t == UserTypePerson
}

func (user *User) IsBot() bool {
	t, _ := user.Type()
	return t == UserTypeBot
}

func (user *User) Name() string {
	return user.JSON.GetString("name")
}

func (user *User) AvatarURL() string {
	return user.JSON.GetString("avatar_url")
}

func (user *User) Email() string {
	if !user.IsPerson() {
		return ""
	}

	if person, ok := user.JSON.GetJSON("person"); ok {
		return person.GetString("email")
	}

	return ""
}

func (user *User) String() string {
	return user.JSON.String()
}
