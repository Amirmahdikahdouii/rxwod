package user

import "time"

type User struct {
	id              UserID
	email           Email
	passwordHash    PasswordHash
	displayName     DisplayName
	emailVerifiedAt *time.Time
	createdAt       time.Time
	updatedAt       time.Time
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

func ReconstructUser(id UserID, email Email, passwordHash PasswordHash, displayName DisplayName, emailVerifiedAt *time.Time, createdAt time.Time, updatedAt time.Time) (User, error) {
	user, err := NewUser(id, email, passwordHash, displayName, createdAt)
	if err != nil {
		return User{}, err
	}
	user.emailVerifiedAt = emailVerifiedAt
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

func (u User) EmailVerifiedAt() *time.Time {
	return u.emailVerifiedAt
}

func (u User) IsEmailVerified() bool {
	return u.emailVerifiedAt != nil
}

func (u User) MarkEmailVerified(now time.Time) User {
	u.emailVerifiedAt = &now
	u.updatedAt = now
	return u
}
