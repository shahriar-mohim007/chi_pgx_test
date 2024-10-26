package repository

import (
	"github.com/gofrs/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `db:"id"`         // UUID for each user
	Name      string    `db:"name"`       // User's full name
	Email     string    `db:"email"`      // User's email (must be unique)
	Password  string    `db:"password"`   // Hashed password for security
	IsActive  bool      `db:"is_active"`  // Indicates if the user is activated
	CreatedAt time.Time `db:"created_at"` // Timestamp of when the user was created
	UpdatedAt time.Time `db:"updated_at"` // Timestamp of the last update
}

type Contact struct {
	ID        uuid.UUID `json:"id" db:"id"`                 // Unique ID for each contact
	UserID    uuid.UUID `json:"user_id" db:"user_id"`       // Foreign key to users table
	Phone     string    `json:"phone" db:"phone"`           // Contact's phone number
	Street    string    `json:"street" db:"street"`         // Street address
	City      string    `json:"city" db:"city"`             // City
	State     string    `json:"state" db:"state"`           // State
	ZipCode   string    `json:"zip_code" db:"zip_code"`     // Zip code
	Country   string    `json:"country" db:"country"`       // Country
	CreatedAt time.Time `json:"created_at" db:"created_at"` // Created timestamp
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // Updated timestamp
}

type ContactWithUserResponse struct {
	ContactID uuid.UUID `json:"contact_id"`
	Phone     string    `json:"phone"`
	Street    string    `json:"street"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	ZipCode   string    `json:"zip_code"`
	Country   string    `json:"country"`

	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
}
