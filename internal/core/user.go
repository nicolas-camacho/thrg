package core

import "github.com/google/uuid"

type UserLookupModel struct {
	ID       uuid.UUID
	Username string
}
