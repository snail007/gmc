package session

type Store interface {
	Load(sessionID string) (session *Session, isExits bool)
	Save(session *Session) (err error)
	Delete(sessionID string) (err error)
}
