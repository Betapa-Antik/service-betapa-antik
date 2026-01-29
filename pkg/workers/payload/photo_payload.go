package payload

import "github.com/google/uuid"

type PhotoFile struct {
	Filename string `json:"filename"`
	Data     []byte `json:"data"`
}

type PhotoUploadPayload struct {
	UserID uuid.UUID   `json:"user_id"`
	Folder string      `json:"folder"` // e.g., "betapa_antik/foto_admin"
	Files  []PhotoFile `json:"files"`  // support single or multiple files
}
