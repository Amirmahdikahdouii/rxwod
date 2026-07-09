package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	appauthz "github.com/rxwod/backend/internal/application/authz"
	appgym "github.com/rxwod/backend/internal/application/gym"
	domainauthz "github.com/rxwod/backend/internal/domain/authz"
	domaingym "github.com/rxwod/backend/internal/domain/gym"
	"github.com/rxwod/backend/internal/domain/user"
)

type GymRepository struct {
	db *DB
}

func NewGymRepository(db *DB) *GymRepository {
	return &GymRepository{db: db}
}

func (r *GymRepository) CreateGymWithOwner(ctx context.Context, aggregate domaingym.Gym, ownerMembership domaingym.Membership) error {
	tx, err := r.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		INSERT INTO gyms (id, name, owner_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, aggregate.ID().String(), aggregate.Name(), aggregate.OwnerID().String(), aggregate.CreatedAt(), aggregate.UpdatedAt()); err != nil {
		return fmt.Errorf("insert gym: %w", err)
	}
	if _, err := tx.Exec(ctx, `
		INSERT INTO gym_memberships (id, gym_id, user_id, role, status, invited_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, ownerMembership.ID().String(), ownerMembership.GymID().String(), ownerMembership.UserID().String(), ownerMembership.Role(), ownerMembership.Status(), nil, ownerMembership.CreatedAt(), ownerMembership.UpdatedAt()); err != nil {
		return fmt.Errorf("insert owner membership: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (r *GymRepository) ListForUser(ctx context.Context, userID user.UserID) ([]appgym.WorkspaceDTO, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT g.id, g.name, gm.role
		FROM gym_memberships gm
		JOIN gyms g ON g.id = gm.gym_id
		WHERE gm.user_id = $1 AND gm.status = 'active'
		ORDER BY g.created_at DESC
	`, userID.String())
	if err != nil {
		return nil, fmt.Errorf("list gyms for user: %w", err)
	}
	defer rows.Close()

	var workspaces []appgym.WorkspaceDTO
	for rows.Next() {
		var workspace appgym.WorkspaceDTO
		if err := rows.Scan(&workspace.ID, &workspace.Name, &workspace.Role); err != nil {
			return nil, fmt.Errorf("scan workspace: %w", err)
		}
		workspaces = append(workspaces, workspace)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate workspaces: %w", err)
	}
	return workspaces, nil
}

func (r *GymRepository) FindByID(ctx context.Context, gymID domaingym.GymID) (domaingym.Gym, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, name, owner_id, created_at, updated_at
		FROM gyms
		WHERE id = $1
	`, gymID.String())
	return scanGym(row)
}

func (r *GymRepository) FindActiveMembership(ctx context.Context, gymID domaingym.GymID, userID user.UserID) (domaingym.Membership, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, gym_id, user_id, role, status, invited_by, created_at, updated_at
		FROM gym_memberships
		WHERE gym_id = $1 AND user_id = $2 AND status = 'active'
	`, gymID.String(), userID.String())
	membership, err := scanMembership(row)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return domaingym.Membership{}, appauthz.ErrActiveMembershipMissing
		}
		return domaingym.Membership{}, err
	}
	return membership, nil
}

func (r *GymRepository) FindMembership(ctx context.Context, gymID domaingym.GymID, userID user.UserID) (domaingym.Membership, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, gym_id, user_id, role, status, invited_by, created_at, updated_at
		FROM gym_memberships
		WHERE gym_id = $1 AND user_id = $2
	`, gymID.String(), userID.String())
	membership, err := scanMembership(row)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return domaingym.Membership{}, appgym.ErrMemberNotFound
		}
		return domaingym.Membership{}, err
	}
	return membership, nil
}

func (r *GymRepository) FindMember(ctx context.Context, gymID domaingym.GymID, userID user.UserID) (appgym.MemberDTO, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT gm.id, u.id, u.email, u.display_name, gm.role, gm.status
		FROM gym_memberships gm
		JOIN users u ON u.id = gm.user_id
		WHERE gm.gym_id = $1 AND gm.user_id = $2
	`, gymID.String(), userID.String())
	var member appgym.MemberDTO
	if err := row.Scan(&member.MembershipID, &member.UserID, &member.Email, &member.DisplayName, &member.Role, &member.Status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return appgym.MemberDTO{}, appgym.ErrMemberNotFound
		}
		return appgym.MemberDTO{}, fmt.Errorf("scan gym member: %w", err)
	}
	return member, nil
}

func (r *GymRepository) DeleteMembership(ctx context.Context, gymID domaingym.GymID, userID user.UserID) error {
	tag, err := r.db.pool.Exec(ctx, `
		DELETE FROM gym_memberships
		WHERE gym_id = $1 AND user_id = $2
	`, gymID.String(), userID.String())
	if err != nil {
		return fmt.Errorf("delete gym membership: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return appgym.ErrMemberNotFound
	}
	return nil
}

func (r *GymRepository) ListMembers(ctx context.Context, gymID domaingym.GymID, filter appgym.ListMembersFilter) (appgym.ListMembersResult, error) {
	var total int
	if err := r.db.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM gym_memberships gm
		WHERE gm.gym_id = $1
	`, gymID.String()).Scan(&total); err != nil {
		return appgym.ListMembersResult{}, fmt.Errorf("count gym members: %w", err)
	}

	offset := (filter.Page - 1) * filter.Limit
	rows, err := r.db.pool.Query(ctx, `
		SELECT gm.id, u.id, u.email, u.display_name, gm.role, gm.status
		FROM gym_memberships gm
		JOIN users u ON u.id = gm.user_id
		WHERE gm.gym_id = $1
		ORDER BY gm.created_at ASC
		LIMIT $2 OFFSET $3
	`, gymID.String(), filter.Limit, offset)
	if err != nil {
		return appgym.ListMembersResult{}, fmt.Errorf("list gym members: %w", err)
	}
	defer rows.Close()

	var members []appgym.MemberDTO
	for rows.Next() {
		var member appgym.MemberDTO
		if err := rows.Scan(&member.MembershipID, &member.UserID, &member.Email, &member.DisplayName, &member.Role, &member.Status); err != nil {
			return appgym.ListMembersResult{}, fmt.Errorf("scan gym member: %w", err)
		}
		members = append(members, member)
	}
	if err := rows.Err(); err != nil {
		return appgym.ListMembersResult{}, fmt.Errorf("iterate gym members: %w", err)
	}
	return appgym.ListMembersResult{Items: members, Total: total}, nil
}

func (r *GymRepository) FindUserByEmail(ctx context.Context, email user.Email) (user.User, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, display_name, email_verified_at, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email)
	aggregate, err := scanUser(row)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return user.User{}, appgym.ErrInviteeNotFound
		}
		return user.User{}, err
	}
	return aggregate, nil
}

func (r *GymRepository) FindUserByID(ctx context.Context, userID user.UserID) (user.User, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, display_name, email_verified_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`, userID.String())
	aggregate, err := scanUser(row)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return user.User{}, appgym.ErrInviteeNotFound
		}
		return user.User{}, err
	}
	return aggregate, nil
}

func (r *GymRepository) UpsertMembership(ctx context.Context, membership domaingym.Membership) error {
	invitedBy := membership.InvitedBy()
	_, err := r.db.pool.Exec(ctx, `
		INSERT INTO gym_memberships (id, gym_id, user_id, role, status, invited_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (gym_id, user_id) DO UPDATE SET
			role = EXCLUDED.role,
			status = EXCLUDED.status,
			invited_by = EXCLUDED.invited_by,
			updated_at = EXCLUDED.updated_at
	`, membership.ID().String(), membership.GymID().String(), membership.UserID().String(), membership.Role(), membership.Status(), userIDString(invitedBy), membership.CreatedAt(), membership.UpdatedAt())
	if err != nil {
		return fmt.Errorf("upsert gym membership: %w", err)
	}
	return nil
}

func (r *GymRepository) SaveInvitation(ctx context.Context, invitation domaingym.Invitation, tokenHash string) error {
	_, err := r.db.pool.Exec(ctx, `
		INSERT INTO gym_invitations (id, gym_id, email, role, status, invited_by, expires_at, created_at, updated_at, token_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (gym_id, email, role) WHERE status = 'pending'
		DO UPDATE SET
			invited_by = EXCLUDED.invited_by,
			expires_at = EXCLUDED.expires_at,
			updated_at = EXCLUDED.updated_at,
			token_hash = EXCLUDED.token_hash
	`, invitation.ID().String(), invitation.GymID().String(), invitation.Email(), invitation.Role(), invitation.Status(), invitation.InvitedBy().String(), invitation.ExpiresAt(), invitation.CreatedAt(), invitation.UpdatedAt(), tokenHash)
	if err != nil {
		return fmt.Errorf("save gym invitation: %w", err)
	}
	return nil
}

func (r *GymRepository) FindPendingInvitationsByEmail(ctx context.Context, email user.Email) ([]domaingym.Invitation, error) {
	rows, err := r.db.pool.Query(ctx, `
		SELECT id, gym_id, email, role, status, invited_by, expires_at, created_at, updated_at
		FROM gym_invitations
		WHERE email = $1 AND status = 'pending'
		ORDER BY created_at ASC
	`, email)
	if err != nil {
		return nil, fmt.Errorf("find pending invitations: %w", err)
	}
	defer rows.Close()

	var invitations []domaingym.Invitation
	for rows.Next() {
		invitation, err := scanInvitation(rows)
		if err != nil {
			return nil, err
		}
		invitations = append(invitations, invitation)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate invitations: %w", err)
	}
	return invitations, nil
}

func (r *GymRepository) FindInvitationPreviewByTokenHash(ctx context.Context, tokenHash string) (appgym.InvitationPreviewDTO, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT gi.gym_id, g.name, gi.email, gi.role, gi.status, gi.expires_at
		FROM gym_invitations gi
		JOIN gyms g ON g.id = gi.gym_id
		WHERE gi.token_hash = $1
	`, tokenHash)
	var preview appgym.InvitationPreviewDTO
	if err := row.Scan(&preview.GymID, &preview.GymName, &preview.Email, &preview.Role, &preview.Status, &preview.ExpiresAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return appgym.InvitationPreviewDTO{}, appgym.ErrInvitationNotFound
		}
		return appgym.InvitationPreviewDTO{}, fmt.Errorf("find invitation preview: %w", err)
	}
	return preview, nil
}

func (r *GymRepository) FindPendingInvitationByTokenHash(ctx context.Context, gymID domaingym.GymID, tokenHash string) (domaingym.Invitation, error) {
	row := r.db.pool.QueryRow(ctx, `
		SELECT id, gym_id, email, role, status, invited_by, expires_at, created_at, updated_at
		FROM gym_invitations
		WHERE gym_id = $1 AND token_hash = $2 AND status = 'pending'
	`, gymID.String(), tokenHash)
	invitation, err := scanInvitation(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domaingym.Invitation{}, appgym.ErrInvitationNotFound
		}
		return domaingym.Invitation{}, err
	}
	return invitation, nil
}

func (r *GymRepository) AcceptInvitationWithMembership(ctx context.Context, invitation domaingym.Invitation, membership domaingym.Membership) error {
	tx, err := r.db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		UPDATE gym_invitations
		SET status = $2, updated_at = $3
		WHERE id = $1
	`, invitation.ID().String(), invitation.Status(), invitation.UpdatedAt()); err != nil {
		return fmt.Errorf("accept invitation: %w", err)
	}
	invitedBy := membership.InvitedBy()
	if _, err := tx.Exec(ctx, `
		INSERT INTO gym_memberships (id, gym_id, user_id, role, status, invited_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (gym_id, user_id) DO UPDATE SET
			role = EXCLUDED.role,
			status = EXCLUDED.status,
			invited_by = EXCLUDED.invited_by,
			updated_at = EXCLUDED.updated_at
	`, membership.ID().String(), membership.GymID().String(), membership.UserID().String(), membership.Role(), membership.Status(), userIDString(invitedBy), membership.CreatedAt(), membership.UpdatedAt()); err != nil {
		return fmt.Errorf("upsert invited membership: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func scanGym(scanner rowScanner) (domaingym.Gym, error) {
	var (
		id        string
		name      string
		ownerID   string
		createdAt time.Time
		updatedAt time.Time
	)
	if err := scanner.Scan(&id, &name, &ownerID, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domaingym.Gym{}, ErrNotFound
		}
		return domaingym.Gym{}, fmt.Errorf("scan gym: %w", err)
	}
	return domaingym.ReconstructGym(domaingym.GymID(id), domaingym.GymName(name), user.UserID(ownerID), createdAt, updatedAt)
}

func scanMembership(scanner rowScanner) (domaingym.Membership, error) {
	var (
		id        string
		gymID     string
		userID    string
		role      domainauthz.Role
		status    domaingym.MembershipStatus
		invitedBy sql.NullString
		createdAt time.Time
		updatedAt time.Time
	)
	if err := scanner.Scan(&id, &gymID, &userID, &role, &status, &invitedBy, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domaingym.Membership{}, ErrNotFound
		}
		return domaingym.Membership{}, fmt.Errorf("scan membership: %w", err)
	}
	var inviter *user.UserID
	if invitedBy.Valid {
		value := user.UserID(invitedBy.String)
		inviter = &value
	}
	return domaingym.ReconstructMembership(
		domaingym.MembershipID(id),
		domaingym.GymID(gymID),
		user.UserID(userID),
		role,
		status,
		inviter,
		createdAt,
		updatedAt,
	)
}

func scanInvitation(scanner rowScanner) (domaingym.Invitation, error) {
	var (
		id        string
		gymID     string
		email     string
		role      domainauthz.Role
		status    domaingym.InvitationStatus
		invitedBy string
		expiresAt time.Time
		createdAt time.Time
		updatedAt time.Time
	)
	if err := scanner.Scan(&id, &gymID, &email, &role, &status, &invitedBy, &expiresAt, &createdAt, &updatedAt); err != nil {
		return domaingym.Invitation{}, fmt.Errorf("scan invitation: %w", err)
	}
	return domaingym.ReconstructInvitation(
		domaingym.InvitationID(id),
		domaingym.GymID(gymID),
		user.Email(email),
		role,
		status,
		user.UserID(invitedBy),
		expiresAt,
		createdAt,
		updatedAt,
	)
}

func userIDString(id *user.UserID) any {
	if id == nil {
		return nil
	}
	return id.String()
}
