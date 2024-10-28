package mocks

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"go_chi_pgx/repository"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*repository.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*repository.User), args.Error(1)
}

func (m *MockRepository) CreateUser(ctx context.Context, user *repository.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) ActivateUserByID(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepository) GetAllContacts(ctx context.Context, userID uuid.UUID, limit, offset int) ([]repository.Contact, error) {
	args := m.Called(ctx, userID, limit, offset)

	// Ensure args.Get(0) is a non-nil and correct type ([]repository.Contact)
	if contacts, ok := args.Get(0).([]repository.Contact); ok && contacts != nil {
		return contacts, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) CreateContact(ctx context.Context, contact *repository.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *MockRepository) GetContactByID(ctx context.Context, contactID uuid.UUID) (*repository.ContactWithUserResponse, error) {
	args := m.Called(ctx, contactID)
	if contact, ok := args.Get(0).(*repository.ContactWithUserResponse); ok {
		return contact, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) PatchContactByID(ctx context.Context, contactID uuid.UUID, contact *repository.Contact) error {
	args := m.Called(ctx, contactID, contact)
	return args.Error(0)
}

func (m *MockRepository) DeleteContactByID(ctx context.Context, contactID uuid.UUID) error {
	args := m.Called(ctx, contactID)
	return args.Error(0)
}

func (m *MockRepository) GetContactsCount(ctx context.Context, userID uuid.UUID) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}
