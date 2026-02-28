package task

import (
	"planelite-backend/internal/common"
)

// CanCreateTask: only PROJECT_MANAGER (and ADMIN) can create tasks.
func CanCreateTask(role common.Role) bool {
	return role == common.RoleAdmin || role == common.RoleProjectManager
}

// CanUpdateTaskStatusOrPriority: developers/users can update status & priority; admin full.
func CanUpdateTaskStatusOrPriority(role common.Role) bool {
	return role == common.RoleAdmin || role == common.RoleProjectManager || role == common.RoleUser
}

// CanUpdateTaskFull: admin and PROJECT_MANAGER can update all fields.
func CanUpdateTaskFull(role common.Role) bool {
	return role == common.RoleAdmin || role == common.RoleProjectManager
}
