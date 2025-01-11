package preload

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/models/types"
	"awasm-portfolio/internal/repository"
)

func PreloadData(repo *repository.InMemoryRepository) {
	// Preload namespaces
	repo.Create(&types.Namespace{
		Name: "default",
	})
	repo.Create(&types.Namespace{
		Name: "dev",
	})
	repo.Create(&types.Namespace{
		Name: "test",
	})

	// Preload profiles with owned resources
	repo.Create(&types.Profile{
		Name:      "john-doe",
		Namespace: "default",
	})

	repo.Create(&types.Certifications{
		Name:      "john-doe-certifications",
		Namespace: "default",
		OwnerRef: models.OwnerReference{
			Kind:      "profile",
			Name:      "john-doe",
			Namespace: "default",
		},
	})

	repo.Create(&types.Contact{
		Name:      "john-doe-contact",
		Namespace: "default",
		OwnerRef: models.OwnerReference{
			Kind:      "profile",
			Name:      "john-doe",
			Namespace: "default",
		},
	})

	repo.Create(&types.Profile{
		Name:      "jane-doe",
		Namespace: "dev",
	})

	repo.Create(&types.Certifications{
		Name:      "jane-doe-certifications",
		Namespace: "dev",
		OwnerRef: models.OwnerReference{
			Kind:      "profile",
			Name:      "jane-doe",
			Namespace: "dev",
		},
	})

	repo.Create(&types.Contact{
		Name:      "jane-doe-contact",
		Namespace: "dev",
		OwnerRef: models.OwnerReference{
			Kind:      "profile",
			Name:      "jane-doe",
			Namespace: "dev",
		},
	})

	repo.Create(&types.Profile{
		Name:      "dev-user",
		Namespace: "test",
	})

	repo.Create(&types.Certifications{
		Name:      "dev-user-certifications",
		Namespace: "test",
		OwnerRef: models.OwnerReference{
			Kind:      "profile",
			Name:      "dev-user",
			Namespace: "test",
		},
	})

	repo.Create(&types.Contact{
		Name:      "dev-user-contact",
		Namespace: "test",
		OwnerRef: models.OwnerReference{
			Kind:      "profile",
			Name:      "dev-user",
			Namespace: "test",
		},
	})
}
