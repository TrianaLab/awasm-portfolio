package service

import (
	"awasm-portfolio/internal/repository"
)

type ResourceService interface {
	CreateResource(kind, name, namespace string) (string, error)
	DeleteResource(kind, name, namespace string) (string, error)
	GetResources(kind, name, namespace string) (string, error) // Fixed: Ensure this method is implemented
	DescribeResource(kind, name, namespace string) (string, error)
}

type ResourceServiceImpl struct {
	createService   *CreateService
	deleteService   *DeleteService
	getService      *GetService
	describeService *DescribeService
}

func NewResourceService(repo *repository.InMemoryRepository) ResourceService {
	return &ResourceServiceImpl{
		createService:   NewCreateService(repo),
		deleteService:   NewDeleteService(repo),
		getService:      NewGetService(repo),
		describeService: NewDescribeService(repo),
	}
}

// CreateResource delegates to the CreateService
func (s *ResourceServiceImpl) CreateResource(kind, name, namespace string) (string, error) {
	return s.createService.CreateResource(kind, name, namespace)
}

// DeleteResource delegates to the DeleteService
func (s *ResourceServiceImpl) DeleteResource(kind, name, namespace string) (string, error) {
	return s.deleteService.DeleteResource(kind, name, namespace)
}

// Fixed: GetResources delegates to the GetService
func (s *ResourceServiceImpl) GetResources(kind, name, namespace string) (string, error) {
	return s.getService.GetResources(kind, name, namespace)
}

// DescribeResource delegates to the DescribeService
func (s *ResourceServiceImpl) DescribeResource(kind, name, namespace string) (string, error) {
	return s.describeService.DescribeResource(kind, name, namespace)
}
