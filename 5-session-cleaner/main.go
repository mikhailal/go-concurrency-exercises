//////////////////////////////////////////////////////////////////////
//
// Given is a SessionManager that stores session information in
// memory. The SessionManager itself is working, however, since we
// keep on adding new sessions to the manager our program will
// eventually run out of memory.
//
// Your task is to implement a session cleaner routine that runs
// concurrently in the background and cleans every session that
// hasn't been updated for more than 5 seconds (of course usually
// session times are much longer).
//
// Note that we expect the session to be removed anytime between 5 and
// 7 seconds after the last update. Also, note that you have to be
// very careful in order to prevent race conditions.
//

package main

import (
	"errors"
	"log"
	"sync"
	"time"
)

// SessionManager keeps track of all sessions from creation, updating
// to destroying.
type SessionManager struct {
	sessions map[string]Session
	mt       sync.Mutex
}

// Session stores the session's data
type Session struct {
	Data        map[string]interface{}
	last_active time.Time
}

// NewSessionManager creates a new sessionManager
func NewSessionManager() *SessionManager {
	m := &SessionManager{
		sessions: make(map[string]Session),
	}

	go m.Autoclean()
	return m
}

// CreateSession creates a new session and returns the sessionID
func (m *SessionManager) CreateSession() (string, error) {
	defer m.mt.Unlock()
	m.mt.Lock()
	sessionID, err := MakeSessionID()
	if err != nil {
		return "", err
	}

	m.sessions[sessionID] = Session{
		Data:        make(map[string]interface{}),
		last_active: time.Now(),
	}

	return sessionID, nil
}

// DeleteSession deletes a new session and returns the sessionID
func (m *SessionManager) DeleteSession(sessionID string) error {
	_, ok := m.sessions[sessionID]
	if ok {
		delete(m.sessions, sessionID)
	}

	return nil
}

const poll_time = 250 * time.Millisecond
const timeout = 25 * poll_time

func (m *SessionManager) Autoclean() {
	for {
		m.mt.Lock()
		for key, val := range m.sessions {
			if time.Now().Sub(val.last_active) > timeout {
				m.DeleteSession(key)
			}
		}
		m.mt.Unlock()
		time.Sleep(poll_time)
	}
}

// ErrSessionNotFound returned when sessionID notgg listed in
// SessionManager
var ErrSessionNotFound = errors.New("SessionID does not exists")

// GetSessionData returns data related to session if sessionID is
// found, errors otherwise
func (m *SessionManager) GetSessionData(sessionID string) (map[string]interface{}, error) {
	//	defer m.mt.Unlock()
	m.mt.Lock()
	session, ok := m.sessions[sessionID]
	m.mt.Unlock()
	if !ok {
		return nil, ErrSessionNotFound
	}
	return session.Data, nil
}

// UpdateSessionData overwrites the old session data with the new one
func (m *SessionManager) UpdateSessionData(sessionID string, data map[string]interface{}) error {
	m.mt.Lock()
	_, ok := m.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	// Hint: you should renew expiry of the session here
	m.sessions[sessionID] = Session{
		Data:        data,
		last_active: time.Now(),
	}
	m.mt.Unlock()
	return nil
}

func main() {
	// Create new sessionManager and new session
	m := NewSessionManager()
	sID, err := m.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Created new session with ID", sID)

	// Update session data
	data := make(map[string]interface{})
	data["website"] = "longhoang.de"

	err = m.UpdateSessionData(sID, data)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Update session data, set website to longhoang.de")

	// Retrieve data from manager again
	updatedData, err := m.GetSessionData(sID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Get session data:", updatedData)

}
