package authz

type Role string

const (
	RoleOwner   Role = "owner"
	RoleCoach   Role = "coach"
	RoleAthlete Role = "athlete"
)

type Permission string

const (
	PermissionGymRead             Permission = "gym:read"
	PermissionGymManage           Permission = "gym:manage"
	PermissionMemberInviteCoach   Permission = "member:invite_coach"
	PermissionMemberInviteAthlete Permission = "member:invite_athlete"
	PermissionMemberList          Permission = "member:list"
	PermissionMemberUpdateRole    Permission = "member:update_role"
	PermissionMemberRemove        Permission = "member:remove"
	PermissionWODCreate           Permission = "wod:create"
	PermissionWODRead             Permission = "wod:read"
	PermissionWODUpdate           Permission = "wod:update"
	PermissionWODPublish          Permission = "wod:publish"
	PermissionWODDelete           Permission = "wod:delete"
	PermissionWODResultSubmit     Permission = "wod_result:submit"
	PermissionWODResultRead       Permission = "wod_result:read"
)

var rolePermissions = map[Role]map[Permission]struct{}{
	RoleOwner: {
		PermissionGymRead:             {},
		PermissionGymManage:           {},
		PermissionMemberInviteCoach:   {},
		PermissionMemberInviteAthlete: {},
		PermissionMemberList:          {},
		PermissionMemberUpdateRole:    {},
		PermissionMemberRemove:        {},
		PermissionWODCreate:           {},
		PermissionWODRead:             {},
		PermissionWODUpdate:           {},
		PermissionWODPublish:          {},
		PermissionWODDelete:           {},
		PermissionWODResultSubmit:     {},
		PermissionWODResultRead:       {},
	},
	RoleCoach: {
		PermissionWODCreate:       {},
		PermissionWODRead:         {},
		PermissionWODUpdate:       {},
		PermissionWODPublish:      {},
		PermissionWODResultSubmit: {},
		PermissionWODResultRead:   {},
	},
	RoleAthlete: {
		PermissionWODRead:         {},
		PermissionWODResultSubmit: {},
		PermissionWODResultRead:   {},
	},
}

func HasPermission(role Role, permission Permission) bool {
	permissions, ok := rolePermissions[role]
	if !ok {
		return false
	}
	_, ok = permissions[permission]
	return ok
}

func IsValidRole(role Role) bool {
	_, ok := rolePermissions[role]
	return ok
}
