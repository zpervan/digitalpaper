package tests

import (
	"digitalpaper/backend/core"
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ValidatorTestSuite struct {
	suite.Suite
	validator core.Validator

	minValue         int
	maxValue         int
	minValidString   string
	maxValidString   string
	minInvalidString string
	maxInvalidString string
	validMail        string
	invalidMails     []string
}

func (ts *ValidatorTestSuite) SetupSuite() {
	fmt.Println("Setting up validator test suite...")

	ts.minValue = 3
	ts.maxValue = 5
	ts.minValidString = "abc"
	ts.maxValidString = "abcde"
	ts.minInvalidString = "ab"
	ts.maxInvalidString = "abcdef"
	ts.validMail = "somemail@mail.com"

	ts.invalidMails = append(ts.invalidMails, "wrong@")
	ts.invalidMails = append(ts.invalidMails, "error")
	ts.invalidMails = append(ts.invalidMails, "error@.com")

	fmt.Println("Setting up validator test suite... COMPLETE")
}

func (ts *ValidatorTestSuite) SetupTest() {
	ts.validator = core.Validator{}
	ts.validator.Tag = "test"
}

func (ts *ValidatorTestSuite) Test_GivenValidMinStringLengths_WhenCheckingMinLengths_ThenNoErrorIsReturned() {
	ts.validator.MinLength(&ts.minValidString, ts.minValue)
	ts.Assertions.True(ts.validator.IsValid())
	ts.Assertions.Nil(ts.validator.ValidationResults)
}

func (ts *ValidatorTestSuite) Test_GivenInvalidMinStringLengths_WhenCheckingMinLengths_ThenErrorIsReturned() {
	ts.validator.MinLength(&ts.minInvalidString, ts.minValue)

	ts.Assertions.False(ts.validator.IsValid())
	ts.Assertions.NotNil(ts.validator.ValidationResults)

	expectedErrorMessage := fmt.Sprintf("value is too short - should contain at least %d characters", ts.minValue)
	actualErrorMessage := ts.validator.ValidationResults[0].Error
	ts.Equal(expectedErrorMessage, actualErrorMessage)
}

func (ts *ValidatorTestSuite) Test_GivenValidMaxStringLengths_WhenCheckingMaxLengths_ThenNoErrorIsReturned() {
	ts.validator.MaxLength(&ts.maxValidString, ts.maxValue)
	ts.Assertions.True(ts.validator.IsValid())
	ts.Assertions.Nil(ts.validator.ValidationResults)
}

func (ts *ValidatorTestSuite) Test_GivenInvalidMaxStringLengths_WhenCheckingMaxLengths_ThenErrorIsReturned() {
	ts.validator.MaxLength(&ts.maxInvalidString, ts.maxValue)

	ts.Assertions.False(ts.validator.IsValid())
	ts.Assertions.NotNil(ts.validator.ValidationResults)

	expectedErrorMessage := fmt.Sprintf("value is too long - should contain at least %d characters", ts.maxValue)
	actualErrorMessage := ts.validator.ValidationResults[0].Error
	ts.Equal(expectedErrorMessage, actualErrorMessage)
}

func (ts *ValidatorTestSuite) Test_GivenValidMailAddress_WhenMatchingMailPattern_ThenNoErrorIsReturned() {
	ts.validator.MatchesPattern(&ts.validMail, core.StandardMailAddressPattern)
	ts.Assertions.True(ts.validator.IsValid())
	ts.Assertions.Nil(ts.validator.ValidationResults)
}

func (ts *ValidatorTestSuite) Test_GivenInvalidMailAddress_WhenMatchingMailPattern_ThenErrorIsReturned() {
	for _, invalidMail := range ts.invalidMails {
		ts.validator.MatchesPattern(&invalidMail, core.StandardMailAddressPattern)
	}

	ts.Assertions.False(ts.validator.IsValid())
	ts.Assertions.NotNil(ts.validator.ValidationResults)

	for i, _ := range ts.validator.ValidationResults {
		expectedErrorMessage := "value is not valid"
		actualErrorMessage := ts.validator.ValidationResults[i].Error
		ts.Equal(expectedErrorMessage, actualErrorMessage)
	}
}

func TestRunValidator(t *testing.T) {
	suite.Run(t, new(ValidatorTestSuite))
}
