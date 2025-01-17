package service

import (
	"awasm-portfolio/internal/repository"
	"awasm-portfolio/internal/ui"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type GetService struct {
	repo *repository.InMemoryRepository
	cmd  *cobra.Command
}

func NewGetService(repo *repository.InMemoryRepository, cmd *cobra.Command) *GetService {
	return &GetService{
		repo: repo,
		cmd:  cmd,
	}
}

func (s *GetService) GetResources(kind string, name string, namespace string) (string, error) {
	if name != "" && namespace == "" {
		return "", fmt.Errorf("a resource cannot be retrieved by name across all namespaces")
	}

	// Retrieve the output format from the command flags
	outputFormat, _ := s.cmd.Flags().GetString("output")
	outputFormat = strings.ToLower(outputFormat)

	// Retrieve resources from the repository
	resources, err := s.repo.List(kind, name, namespace)
	if err != nil {
		return "", err
	}

	// Filter namespaces when kind is "all" and namespace is specified
	if strings.ToLower(kind) == "all" {
		for i := 0; i < len(resources); {
			if resources[i].GetKind() == "namespace" {
				resources = append(resources[:i], resources[i+1:]...)
			} else {
				i++
			}
		}
	}

	return ui.FormatTable(resources, outputFormat), nil
}
