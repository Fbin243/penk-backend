package test

import (
	"log"
	"testing"
)

func TestUserFlow(t *testing.T) {
	p := &PineLine{
		ctx: &Context{
			"testingT": t,
		},
		firstNode: NewFirstNode(getUserInfo, updateUser, createNewCharacter),
	}

	err := p.Exec()
	log.Print(err)
}
