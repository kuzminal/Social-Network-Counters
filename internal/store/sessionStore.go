package store

import (
	"SocialNetCounters/models"
)

type SessionStore interface {
	Store

	LoadSession(id string) (models.UserSession, error)
	CreateSession(m *models.UserSession) (models.UserSession, error)
}
