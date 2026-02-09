package events

type EventType string

const (
	NoteUpdated    EventType = "NOTE_UPDATED"
	NoteRestored   EventType = "NOTE_RESTORED"
	SyncConflict   EventType = "SYNC_CONFLICT"
	LoginFailed   EventType = "LOGIN_FAILED"
	TokenRefreshed EventType = "TOKEN_REFRESHED"
)
