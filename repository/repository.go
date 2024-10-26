package repository

import (
	"context"
	"github.com/gofrs/uuid"
)

// Repository defines the methods for user and contact management.

type Repository interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	ActivateUserByID(ctx context.Context, userID uuid.UUID) error
	GetAllContacts(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Contact, error)
	CreateContact(ctx context.Context, contact *Contact) error
	GetContactByID(ctx context.Context, contactID uuid.UUID) (*ContactWithUserResponse, error)
	PatchContact(ctx context.Context, contactID uuid.UUID, contact *Contact) error
	DeleteContactByID(ctx context.Context, contactID uuid.UUID) error
	GetContactsCount(ctx context.Context, userID uuid.UUID) (int, error)
}
