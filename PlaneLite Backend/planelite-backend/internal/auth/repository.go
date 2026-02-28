package auth

// Auth has no own persistence; it uses user.Service for signup/login and issues JWTs.
// This file exists to match the mandated handler → service → repository structure;
// the "repository" for auth is the user package.
