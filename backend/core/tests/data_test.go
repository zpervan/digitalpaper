package tests

import (
	"digitalpaper/backend/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GivenEmptyUserData_WhenCheckingEmptiness_ThenEmptinessCheckResultIsTrue(t *testing.T) {
	emptyUserData := core.User{}

	assert.True(t, emptyUserData.IsEmpty(), "user data should be empty, but isn't")
}

func Test_GivenPopulatedUserData_WhenCheckingEmptiness_ThenEmptinessCheckResultIsFalse(t *testing.T) {
	fullPopulatedUserData := core.User{
		Id:       "weird_id",
		Username: "weirdUsername",
		Name:     "weird name",
		Surname:  "weird surname",
		Mail:     "weird@mail.com",
		Password: "weirdy",
	}

	assert.False(t, fullPopulatedUserData.IsEmpty(), "user data should non-empty, but isn't")

	partialPopulatedUserData := core.User{
		Id:       "",
		Username: "weirdUsername",
		Name:     "",
		Surname:  "weird surname",
		Mail:     "weird@mail.com",
		Password: "",
	}

	assert.False(t, partialPopulatedUserData.IsEmpty(), "user data should be non-empty but isn't")
}
