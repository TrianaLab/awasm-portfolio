package preload

import (
	"awasm-portfolio/internal/repository"
)

func PreloadData(repo *repository.InMemoryRepository) {
	for _, resource := range resources {
		repo.Create(resource)
	}
}
