package common

type Role string

const (
	RoleAdmin          Role = "ADMIN"
	RoleProjectManager Role = "PROJECT_MANAGER"
	RoleUser           Role = "USER"
)
