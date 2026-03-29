package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// Session represents a user session
type Session struct {
	Token     string
	Username  string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// SessionManager manages user sessions
type SessionManager struct {
	sessions map[string]*Session
	mu       sync.RWMutex
	duration time.Duration
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
		duration: 24 * time.Hour, // 24 hours expiry
	}
}

// CreateSession creates a new session for the given username
func (sm *SessionManager) CreateSession(username string) (*Session, error) {
	token := generateToken()
	session := &Session{
		Token:     token,
		Username:  username,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(sm.duration),
	}

	sm.mu.Lock()
	sm.sessions[token] = session
	sm.mu.Unlock()

	// Clean up expired sessions periodically
	go sm.cleanupExpiredSessions()

	return session, nil
}

// ValidateSession validates a session token and returns the session if valid
func (sm *SessionManager) ValidateSession(token string) (*Session, bool) {
	sm.mu.RLock()
	session, exists := sm.sessions[token]
	sm.mu.RUnlock()

	if !exists {
		return nil, false
	}

	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		sm.DeleteSession(token)
		return nil, false
	}

	return session, true
}

// DeleteSession removes a session
func (sm *SessionManager) DeleteSession(token string) {
	sm.mu.Lock()
	delete(sm.sessions, token)
	sm.mu.Unlock()
}

// GetSessionCount returns the number of active sessions
func (sm *SessionManager) GetSessionCount() int {
	sm.mu.RLock()
	count := len(sm.sessions)
	sm.mu.RUnlock()
	return count
}

// cleanupExpiredSessions removes all expired sessions
func (sm *SessionManager) cleanupExpiredSessions() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	now := time.Now()
	for token, session := range sm.sessions {
		if now.After(session.ExpiresAt) {
			delete(sm.sessions, token)
		}
	}
}

// generateToken generates a secure random token
func generateToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based token if crypto/rand fails
		return hex.EncodeToString([]byte(time.Now().String()))
	}
	return hex.EncodeToString(bytes)
}

// SetSessionDuration allows customization of session duration
func (sm *SessionManager) SetSessionDuration(duration time.Duration) {
	sm.mu.Lock()
	sm.duration = duration
	sm.mu.Unlock()
}
