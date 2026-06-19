package authz

import "testing"

func TestRolePermissionMatrix(t *testing.T) {
	tests := []struct {
		name       string
		role       Role
		permission Permission
		want       bool
	}{
		{name: "owner manages members", role: RoleOwner, permission: PermissionMemberInviteCoach, want: true},
		{name: "owner updates member role", role: RoleOwner, permission: PermissionMemberUpdateRole, want: true},
		{name: "owner removes member", role: RoleOwner, permission: PermissionMemberRemove, want: true},
		{name: "coach cannot update member role", role: RoleCoach, permission: PermissionMemberUpdateRole, want: false},
		{name: "coach cannot remove member", role: RoleCoach, permission: PermissionMemberRemove, want: false},
		{name: "owner updates wod", role: RoleOwner, permission: PermissionWODUpdate, want: true},
		{name: "coach creates wod", role: RoleCoach, permission: PermissionWODCreate, want: true},
		{name: "coach updates wod", role: RoleCoach, permission: PermissionWODUpdate, want: true},
		{name: "coach publishes wod", role: RoleCoach, permission: PermissionWODPublish, want: true},
		{name: "athlete reads wod", role: RoleAthlete, permission: PermissionWODRead, want: true},
		{name: "athlete cannot create wod", role: RoleAthlete, permission: PermissionWODCreate, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasPermission(tt.role, tt.permission); got != tt.want {
				t.Fatalf("HasPermission(%q, %q) = %v, want %v", tt.role, tt.permission, got, tt.want)
			}
		})
	}
}
