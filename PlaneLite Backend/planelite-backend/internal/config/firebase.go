package config

// Firebase holds Firebase/Admin SDK config if needed later (e.g. push, auth).
// Kept as stub per file structure; no hardcoded secrets.
type Firebase struct {
	ProjectID string
	Enabled   bool
}

// LoadFirebase loads Firebase config from env. Optional.
func LoadFirebase() Firebase {
	return Firebase{
		ProjectID: getEnv("FIREBASE_PROJECT_ID", ""),
		Enabled:   getEnv("FIREBASE_ENABLED", "false") == "true",
	}
}
