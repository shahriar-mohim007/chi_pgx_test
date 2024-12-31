package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strings"
	"sync"
	"time"
)

type PgxRepository struct {
	db *pgxpool.Pool
}

var (
	once       sync.Once
	repository *PgxRepository
)

func NewPgRepository(databaseUrl string) (*PgxRepository, error) {
	var onceErr error // Local error variable to avoid race conditions.
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// Define connection pool configuration
		config, err := pgxpool.ParseConfig(databaseUrl)
		if err != nil {
			onceErr = fmt.Errorf("invalid database URL: %w", err)
			log.Error().Err(err).Msg("Failed to parse database configuration")
			return
		}

		// Customize pool settings
		config.MaxConns = 1000 // Maximum number of connections
		config.MinConns = 2    // Minimum number of idle connections
		config.MaxConnLifetime = 30 * time.Minute
		config.MaxConnIdleTime = 5 * time.Second
		config.HealthCheckPeriod = 1 * time.Minute

		// Create connection pool
		db, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			onceErr = fmt.Errorf("failed to create connection pool: %w", err)
			log.Error().Err(err).Msg("Database Connection Error")
			return
		}

		// Ping the database to ensure connectivity
		if err := db.Ping(ctx); err != nil {
			onceErr = fmt.Errorf("failed to ping database: %w", err)
			log.Error().Err(err).Msg("Database Ping Error")
			db.Close()
			return
		}

		// Assign the initialized pool to the repository
		repository = &PgxRepository{db: db}
		log.Info().Msg("Database connection pool successfully initialized")
	})
	if repository != nil {
		go monitorPoolStats(repository.db) // Use repository.db for monitoring
	}
	return repository, onceErr
}

// Close cleans up the database connection pool when the application shuts down.
func (repo *PgxRepository) Close() {
	if repo.db != nil {
		repo.db.Close()
		log.Info().Msg("Database connection pool closed")
	}
}

func monitorPoolStats(pool *pgxpool.Pool) {
	for {
		stats := pool.Stat()

		fmt.Printf("Max Connections: %d\n", stats.MaxConns())
		fmt.Printf("Total Connections: %d\n", stats.TotalConns())
		fmt.Printf("Idle Connections: %d\n", stats.IdleConns())
		fmt.Printf("Acquired Connections: %d\n", stats.AcquiredConns())
		time.Sleep(10 * time.Second) // Adjust interval as needed
	}
}
func (repo *PgxRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := `SELECT id, name, email, password, is_active FROM users WHERE email = $1`
	err := repo.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.IsActive)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *PgxRepository) CreateUser(ctx context.Context, user *User) error {
	query := `INSERT INTO users (id, name, email, password, is_active,created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id`
	err := repo.db.QueryRow(ctx, query, user.ID, user.Name, user.Email, user.Password, user.IsActive).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PgxRepository) ActivateUserByID(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET is_active = TRUE WHERE id = $1`
	result, err := repo.db.Exec(ctx, query, userID)
	rowsAffected := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %s not found", userID)
	}
	return nil
}

func (repo *PgxRepository) GetAllContacts(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Contact, error) {
	query := `
		SELECT id, phone, street, city, state, zip_code, country 
		FROM contacts 
		WHERE user_id = $1 
		LIMIT $2 OFFSET $3`
	rows, err := repo.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []Contact
	for rows.Next() {
		var contact Contact
		err := rows.Scan(&contact.ID, &contact.Phone, &contact.Street, &contact.City, &contact.State, &contact.ZipCode, &contact.Country)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	return contacts, rows.Err()
}

func (repo *PgxRepository) CreateContact(ctx context.Context, contact *Contact) error {
	query := `
        INSERT INTO contacts 
        (id, user_id, phone, street, city, state, zip_code, country, created_at, updated_at) 
        VALUES 
        ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
    `
	_, err := repo.db.Exec(
		ctx, query,
		contact.ID, contact.UserID, contact.Phone, contact.Street, contact.City, contact.State, contact.ZipCode, contact.Country,
	)
	return err
}

func (repo *PgxRepository) GetContactByID(ctx context.Context, contactID uuid.UUID) (*ContactWithUserResponse, error) {
	query := `
       SELECT
           contacts.id AS contact_id,
           contacts.phone,
           contacts.street,
           contacts.city,
           contacts.state,
           contacts.zip_code,
           contacts.country,
           users.name AS user_name,
           users.email AS user_email
       FROM
           contacts
       JOIN
           users ON contacts.user_id = users.id
       WHERE
           contacts.id = $1;
   `

	var response ContactWithUserResponse
	err := repo.db.QueryRow(ctx, query, contactID).Scan(
		&response.ContactID,
		&response.Phone,
		&response.Street,
		&response.City,
		&response.State,
		&response.ZipCode,
		&response.Country,
		&response.UserName,
		&response.UserEmail,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no contact found with ID: %s", contactID)
		}
		return nil, err
	}
	return &response, nil
}

func (repo *PgxRepository) PatchContactByID(ctx context.Context, contactID uuid.UUID, contact *Contact) error {

	var queryParts []string
	var args []interface{}
	argID := 1

	if contact.Phone != "" {
		queryParts = append(queryParts, fmt.Sprintf("phone = $%d", argID))
		args = append(args, contact.Phone)
		argID++
	}
	if contact.Street != "" {
		queryParts = append(queryParts, fmt.Sprintf("street = $%d", argID))
		args = append(args, contact.Street)
		argID++
	}
	if contact.City != "" {
		queryParts = append(queryParts, fmt.Sprintf("city = $%d", argID))
		args = append(args, contact.City)
		argID++
	}
	if contact.State != "" {
		queryParts = append(queryParts, fmt.Sprintf("state = $%d", argID))
		args = append(args, contact.State)
		argID++
	}
	if contact.ZipCode != "" {
		queryParts = append(queryParts, fmt.Sprintf("zip_code = $%d", argID))
		args = append(args, contact.ZipCode)
		argID++
	}
	if contact.Country != "" {
		queryParts = append(queryParts, fmt.Sprintf("country = $%d", argID))
		args = append(args, contact.Country)
		argID++
	}

	if len(queryParts) == 0 {
		return fmt.Errorf("no fields provided to update")
	}

	query := fmt.Sprintf("UPDATE contacts SET %s WHERE id = $%d", strings.Join(queryParts, ", "), argID)
	args = append(args, contactID)

	_, err := repo.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PgxRepository) DeleteContactByID(ctx context.Context, contactID uuid.UUID) error {
	query := `
       DELETE FROM contacts
       WHERE id = $1;
   `

	result, err := repo.db.Exec(ctx, query, contactID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (repo *PgxRepository) GetContactsCount(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM contacts 
		WHERE user_id = $1`

	var count int
	err := repo.db.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
