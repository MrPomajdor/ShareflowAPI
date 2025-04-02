package entity

// User represents a user.
type User struct {
	ID             int
	Email          string
	HashedPassword string
	FirstName      string
	LastName       string
	AuthCode       string
	CreatedAt      string
	LastLogin      string
	LastLoginIP    string
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
	return u.FirstName
}

// GetEmail returns the user email.
func (u User) GetEmail() string {
	return u.Email
}
