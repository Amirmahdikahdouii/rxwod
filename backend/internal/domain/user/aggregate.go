package user

import "time"

type User struct {
	id           UserID
	email        Email
	passwordHash PasswordHash
	displayName  DisplayName
	createdAt    time.Time
	updatedAt    time.Time
}

func NewUser(id UserID, email Email, passwordHash PasswordHash, displayName DisplayName, now time.Time) (User, error) {
	if err := validateEmail(email); err != nil {
		return User{}, err
	}
	if passwordHash == "" {
		return User{}, ErrPasswordHashEmpty
	}
	if err := validateDisplayName(displayName); err != nil {
		return User{}, err
	}
	return User{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		displayName:  displayName,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func ReconstructUser(id UserID, email Email, passwordHash PasswordHash, displayName DisplayName, createdAt time.Time, updatedAt time.Time) (User, error) {
	user, err := NewUser(id, email, passwordHash, displayName, createdAt)
	if err != nil {
		return User{}, err
	}
	user.updatedAt = updatedAt
	return user, nil
}

func (u User) ID() UserID {
	return u.id
}

func (u User) Email() Email {
	return u.email
}

func (u User) PasswordHash() PasswordHash {
	return u.passwordHash
}

func (u User) DisplayName() DisplayName {
	return u.displayName
}

func (u User) CreatedAt() time.Time {
	return u.createdAt
}

func (u User) UpdatedAt() time.Time {
	return u.updatedAt
}
