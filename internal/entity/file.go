package entity

import (
	"encoding/json"
)

// File represents a file record.
type FileNode struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	CategoryName string `json:"category_name"`
	OwnerID      int    `json:"owner_id"`
	CreatedAt    string `json:"created_at"`
	Size         *int   `json:"size"`
}

// Convert to JSON
func (node *FileNode) ToJSON() (string, error) {
	jsonData, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
