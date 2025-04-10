package storage

import (
	"github.com/cowpeatechnology/slg-game-server/proto"
)

// Storage defines the interface for data storage operations
type Storage interface {
	// GetPlayer retrieves a player's data by ID
	GetPlayer(id string) (*proto.PlayerData, error)
	// SavePlayer saves a new player's data
	SavePlayer(player *proto.PlayerData) error
	// UpdatePlayer updates an existing player's data
	UpdatePlayer(player *proto.PlayerData) error
	// DeletePlayer deletes a player's data
	DeletePlayer(id string) error
	// Close closes the storage connection
	Close() error
}

// StorageFactory defines the interface for creating storage instances
type StorageFactory interface {
	// CreateStorage creates a new storage instance
	CreateStorage() (Storage, error)
}
