package acl

import (
	"time"

	"github.com/bww/go-util/v1/uuid"
)

type Authorization struct {
	Id          uuid.UUID `json:"id" db:"id,pk"`
	Key         string    `json:"api_key" db:"key"`
	Secret      string    `json:"api_secret" db:"secret"`
	Description string    `json:"description,omitempty" db:"description"`
	Policies    []Policy  `json:"policies"`
	Active      bool      `json:"active" db:"active"`
	Created     time.Time `json:"created_at" db:"created_at"`
}
