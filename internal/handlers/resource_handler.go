package handlers

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/storage"
	"awasm-portfolio/internal/ui"
	"fmt"
)

func Get(resourceManager *storage.ResourceManager, resource, name, namespace string, allNamespaces bool, formatter ui.Formatter) string {
	// Validate namespace unless listing across all namespaces
	if !allNamespaces && resource != "namespace" {
		if err := resourceManager.ValidateNamespace(namespace); err != nil {
			return formatter.Error(fmt.Sprintf("Cannot perform action 'get' for resource '%s' because namespace '%s' does not exist", resource, namespace))
		}
	}

	resources, err := resourceManager.GetAll(resource, namespace, allNamespaces)
	if err != nil {
		return formatter.Error(err.Error())
	}

	if len(resources) == 0 {
		return formatter.Error(fmt.Sprintf("No resources found for '%s'", resource))
	}

	// Pass resources directly to FormatTable
	return formatter.FormatTable(resources)
}

func Describe(resourceManager *storage.ResourceManager, resource, name, namespace string, formatter ui.Formatter) string {
	// Validate namespace
	if resource != "namespace" {
		if err := resourceManager.ValidateNamespace(namespace); err != nil {
			return formatter.Error(fmt.Sprintf("Cannot perform action 'describe' for resource '%s' because namespace '%s' does not exist", resource, namespace))
		}
	}

	resources, err := resourceManager.GetAll(resource, namespace, false)
	if err != nil {
		return formatter.Error(err.Error())
	}

	res, ok := resources[name]
	if !ok {
		return formatter.Error(fmt.Sprintf("Resource '%s' not found", name))
	}

	return formatter.FormatDetails(res.GetFields())
}

func Create(resourceManager *storage.ResourceManager, resource, name, namespace string, data map[string]interface{}, formatter ui.Formatter) string {
	// Validate namespace unless creating a namespace
	if resource != "namespace" {
		if err := resourceManager.ValidateNamespace(namespace); err != nil {
			return formatter.Error(fmt.Sprintf("Cannot perform action 'create' for resource '%s' because namespace '%s' does not exist", resource, namespace))
		}
	}

	resourceObj := models.ResourceBase{
		Name:       name,
		Namespace:  namespace,
		Data:       data,
		Namespaced: resourceManager.IsNamespaced(resource),
	}

	err := resourceManager.Create(resource, resourceObj)
	if err != nil {
		return formatter.Error(err.Error())
	}

	return formatter.Success(fmt.Sprintf("Resource '%s' created successfully", name))
}

func Delete(resourceManager *storage.ResourceManager, resource, name, namespace string, formatter ui.Formatter) string {
	// Validate namespace
	if resource != "namespace" {
		if err := resourceManager.ValidateNamespace(namespace); err != nil {
			return formatter.Error(fmt.Sprintf("Cannot perform action 'delete' for resource '%s' because namespace '%s' does not exist", resource, namespace))
		}
	}

	err := resourceManager.Delete(resource, name, namespace)
	if err != nil {
		return formatter.Error(err.Error())
	}

	return formatter.Success(fmt.Sprintf("Resource '%s' deleted successfully", name))
}
