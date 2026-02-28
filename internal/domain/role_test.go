package domain

import "testing"

func TestRoleHasPermission(t *testing.T) {
	tests := []struct {
		name     string
		perms    []Permission
		resource string
		write    bool
		want     bool
	}{
		{
			"read granted",
			[]Permission{{Resource: "events", Read: true, Write: false}},
			"events", false, true,
		},
		{
			"write granted",
			[]Permission{{Resource: "events", Read: true, Write: true}},
			"events", true, true,
		},
		{
			"write denied",
			[]Permission{{Resource: "events", Read: true, Write: false}},
			"events", true, false,
		},
		{
			"resource not found",
			[]Permission{{Resource: "events", Read: true, Write: true}},
			"schedule", false, false,
		},
		{
			"empty permissions",
			nil,
			"events", false, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Role{Permissions: tt.perms}
			if got := r.HasPermission(tt.resource, tt.write); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoleEntityKind(t *testing.T) {
	r := &Role{}
	if kind := r.EntityKind(); kind != "Role" {
		t.Errorf("expected 'Role', got %q", kind)
	}
}
