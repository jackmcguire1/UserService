package user

import (
	"testing"

	"github.com/jackmcguire1/UserService/pkg/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUser(t *testing.T) {
	user := &User{
		FirstName:   "John",
		LastName:    "Doe",
		CountryCode: "US",
	}
	mockRepo := &MockRepository{}

	mockRepo.On("GetUser", mock.Anything).Return(user, nil)
	svc, err := NewService(&Resources{
		Repo: mockRepo,
	})
	assert.NoError(t, err)

	resp, err := svc.GetUser("100249558")
	assert.NoError(t, err)
	assert.Equal(t, resp.FirstName, "John")
}

func TestPutUser(t *testing.T) {
	user := &User{
		FirstName:   "John",
		LastName:    "Doe",
		CountryCode: "us",
		Email:       "jack@blah.com",
	}

	mockRepo := &MockRepository{}

	mockRepo.On("PutUser", user).Return(nil)
	mockRepo.On("GetUserByEmail", mock.Anything).Return(nil, utils.ErrNotFound)
	svc, err := NewService(&Resources{
		Repo: mockRepo,
	})
	assert.NoError(t, err)

	user, err = svc.PutUser(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Saved)
	assert.Equal(t, "US", user.CountryCode)
}

func TestPutInvalidRequest(t *testing.T) {
	user := &User{
		ID:          "100249558",
		FirstName:   "John",
		LastName:    "Doe",
		CountryCode: "asas",
		Email:       "@b",
	}

	svc, err := NewService(&Resources{})
	assert.NoError(t, err)

	user, err = svc.PutUser(user)
	assert.ErrorIs(t, err, utils.ValidationErr)
}

func TestDeleteUser(t *testing.T) {
	mockRepo := &MockRepository{}

	mockRepo.On("DeleteUser", mock.Anything).Return(nil)
	svc, err := NewService(&Resources{
		Repo: mockRepo,
	})
	assert.NoError(t, err)

	err = svc.DeleteUser("100249558")
	assert.NoError(t, err)
}

func TestGetUsersByCountry(t *testing.T) {
	users := []*User{
		&User{
			FirstName:   "Bob",
			LastName:    "Ballmer",
			CountryCode: "GB",
		},
		&User{
			FirstName:   "Garry",
			LastName:    "Stevens",
			CountryCode: "GB",
		},
	}

	mockRepo := &MockRepository{}

	mockRepo.On("GetUsersByCountry", "GB").Return(users, nil)

	svc, err := NewService(&Resources{
		Repo: mockRepo,
	})
	assert.NoError(t, err)

	resp, err := svc.GetUsersByCountry("GB")
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	assert.Len(t, resp, 2)
}
