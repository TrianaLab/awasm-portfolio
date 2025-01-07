package handlers

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/storage"
	"awasm-portfolio/internal/ui"
	"fmt"
)

func Get(resourceManager *storage.ResourceManager, resource, name, namespace string, allNamespaces bool, formatter ui.Formatter) string {
	resources, err := resourceManager.GetAll(resource, namespace, allNamespaces)
	if err != nil {
		return formatter.Error(err.Error())
	}

	if name != "" {
		res, ok := resources[name]
		if !ok {
			return formatter.Error(fmt.Sprintf("Resource '%s' not found", name))
		}
		return formatter.FormatDetails(res.GetFields())
	}

	// Convert map[string]models.ResourceBase to map[string]interface{}
	convertedResources := make(map[string]interface{}, len(resources))
	for k, v := range resources {
		convertedResources[k] = v
	}

	return formatter.FormatTable(convertedResources)
}

func Describe(resourceManager *storage.ResourceManager, resource, name, namespace string, formatter ui.Formatter) string {
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
	err := resourceManager.Delete(resource, name, namespace)
	if err != nil {
		return formatter.Error(err.Error())
	}

	return formatter.Success(fmt.Sprintf("Resource '%s' deleted successfully", name))
}
