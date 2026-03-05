package handlers

import (
	"github.com/GunarsK-rpg/public-api/internal/cache"
	"github.com/GunarsK-rpg/public-api/internal/repository"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	repo  repository.Repository
	cache *cache.Cache
}

// New creates a new Handler instance.
func New(repo repository.Repository, cache *cache.Cache) *Handler {
	return &Handler{repo: repo, cache: cache}
}
