package entity

import (
	"encoding/json"
)

// File represents a file record.
type FileNode struct {
	ID        string      `json:"id"`
	Parent_ID *int        `json:"parent_id"`
	Name      string      `json:"name"`
	IsDir     bool        `json:"is_dir"`
	Children  []*FileNode `json:"children,omitempty"`
	CreatedAt string      `json:"created_at"`
	OwnerID   int         `json:"owner_id"`
	Path      string      `json:"path"`
}

// Convert to JSON
func (node *FileNode) ToJSON() (string, error) {
	jsonData, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// AddChild adds a file or directory to a directory
func (dir *FileNode) AddChild(node *FileNode) {
	if dir.IsDir {
		dir.Children = append(dir.Children, node)
	}
}
