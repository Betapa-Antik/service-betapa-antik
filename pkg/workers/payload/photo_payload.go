package payload

import "github.com/google/uuid"

type PhotoFile struct {
	Filename string `json:"filename"`
	Path     string `json:"path"` // ⬅️ PATH FILE, BUKAN BYTE
}
type PhotoUploadPayload struct {
	ID     uuid.UUID   `json:"id"`
	Folder string      `json:"folder"` // e.g., "betapa_antik/foto_admin"
	Files  []PhotoFile `json:"files"`  // support single or multiple files
}
