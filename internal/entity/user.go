package entity

// User represents a user.
type User struct {
	ID             int
	Email          string
	HashedPassword string
	FirstName      string
	LastName       string
	ProfileIMG     string
	AuthCode       string
	CreatedAt      string
	LastLogin      string
	LastLoginIP    string
}

// Identity represents an authenticated user identity.
type Identity interface {
	// GetID returns the user ID.
	GetID() int
	// GetFirstName returns the user first name.
	GetFirstName() string
	// GetLastName returns the user last name.
	GetLastName() string
	// GetName returns the user email.
	GetEmail() string
	// GetProfile returns the user profile picture URL.
	GetProfile() string
}

// GetID returns the user ID.
func (u User) GetID() int {
	return u.ID
}

// GetName returns the user first name.
func (u User) GetFirstName() string {
	return u.FirstName
}

// GetName returns the user last name.
func (u User) GetLastName() string {
	return u.LastName
}

// GetEmail returns the user email.
func (u User) GetEmail() string {
	return u.Email
}

// GetProfile returns the user profile picture URL.
func (u User) GetProfile() string {
	return u.ProfileIMG
}
