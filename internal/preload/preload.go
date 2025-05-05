package preload

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/models/types"
	"awasm-portfolio/internal/repository"
	"reflect"
	"time"
)

func PreloadData(repo *repository.InMemoryRepository) {
	timestamp := time.Now()

	namespace := &types.Namespace{
		Kind: "namespace",
		Name: "default",
	}

	ownerRef := models.OwnerReference{
		Kind:      "resume",
		Name:      "eduardo-diaz",
		Namespace: namespace.Name,
	}

	basics := &types.Basics{
		OwnerRef: ownerRef,
		Name:     "Eduardo Díaz",
		Label:    "Machine Learning Operations Engineer",
		Email:    "edudiazasencio@gmail.com",
		Url:      "https://github.com/edu-diaz",
		Summary:  "MLOps Engineer with expertise in Kubernetes, Cloud Infrastructure and Machine Learning",
		Location: types.Location{
			City:        "Sevilla",
			Region:      "Andalucía",
			CountryCode: "ES",
		},
		Profiles: []types.Profile{
			{
				Network:  "LinkedIn",
				Username: "eduardo-diaz-asencio",
				Url:      "https://www.linkedin.com/in/eduardo-diaz-asencio/",
			},
			{
				Network:  "GitHub",
				Username: "edu-diaz",
				Url:      "https://github.com/edu-diaz",
			},
		},
	}

	work := []types.Work{
		{
			OwnerRef:          ownerRef,
			Name:              "emergence-ai-mlops",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Position:          "Machine Learning Operations Engineer",
			URL:               "https://emergence.ai",
			StartDate:         "2024-07",
			EndDate:           "",
			Summary:           "As an MLOps Engineer at Emergence, I'm actively involved in building scalable and efficient AI infrastructure utilizing technologies like Terraform, Crossplane, Prometheus, Istio, Keda, Kyverno, etc.",
			Highlights:        []string{"Provisioning GKE clusters", "Implementing advanced networking capabilities"},
		},
	}

	volunteer := []types.Volunteer{
		{
			OwnerRef:          ownerRef,
			Name:              "sevilla-kubernetes",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Position:          "Community Leader",
			URL:               "https://community.cncf.io/sevilla/",
			StartDate:         "2023-01",
			EndDate:           "",
			Summary:           "Leading the Kubernetes and Cloud Native community in Sevilla",
			Highlights:        []string{"Organizing monthly meetups", "Growing the community from 0 to 100+ members"},
		},
	}

	education := []types.Education{
		{
			OwnerRef:          ownerRef,
			Name:              "masters-ai",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			URL:               "",
			Area:              "Artificial Intelligence",
			StudyType:         "M.Eng",
			StartDate:         "2021",
			EndDate:           "2022",
			Score:             "9.5",
			Courses:           []string{"Machine Learning", "Deep Learning", "Natural Language Processing"},
		},
	}

	awards := []types.Award{
		{
			OwnerRef:          ownerRef,
			Name:              "best-student-award",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Date:              "2022",
			Awarder:           "Universidad Internacional de Valencia",
			Summary:           "Awarded for achieving the highest GPA in the Master's program",
		},
	}

	certificates := []types.Certificate{
		{
			OwnerRef:          ownerRef,
			Name:              "cka",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Date:              "2023",
			Issuer:            "Cloud Native Computing Foundation",
			URL:               "https://www.credly.com/badges/f1c5619d-f6a1-4988-8393-5f9a21455736/linked_in_profile",
		},
	}

	publications := []types.Publication{
		{
			OwnerRef:          ownerRef,
			Name:              "container2wasm-paper",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Publisher:         "arXiv",
			ReleaseDate:       "2024",
			URL:               "https://arxiv.org/container2wasm",
			Summary:           "Research paper about converting container images to WebAssembly modules",
		},
	}

	skills := []types.Skill{
		{
			OwnerRef:          ownerRef,
			Name:              "kubernetes",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Level:             "Expert",
			Keywords:          []string{"Container Orchestration", "Cloud Native", "DevOps"},
		},
	}

	languages := []types.Language{
		{
			OwnerRef:          ownerRef,
			Name:              "spanish",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Fluency:           "Native",
		},
	}

	interests := []types.Interest{
		{
			OwnerRef:          ownerRef,
			Name:              "cloud-native",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Keywords:          []string{"Kubernetes", "WebAssembly", "Microservices"},
		},
	}

	references := []types.Reference{
		{
			OwnerRef:          ownerRef,
			Name:              "john-doe",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Reference:         "Eduardo is an exceptional engineer with deep knowledge in cloud technologies and machine learning infrastructure.",
		},
	}

	projects := []types.Project{
		{
			OwnerRef:          ownerRef,
			Name:              "container2wasm",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Description:       "Container to WASM image converter that enables to run the container on WASM",
			StartDate:         "2023-06",
			EndDate:           "",
			URL:               "https://github.com/container2wasm/container2wasm",
			Highlights:        []string{"Core contributor", "Integration with Docker Desktop"},
		},
	}

	resume := &types.Resume{
		Kind:              "resume",
		Name:              "eduardo-diaz",
		Namespace:         namespace.Name,
		CreationTimestamp: timestamp,
		Basics:            *basics,
		Work:              work,
		Volunteer:         volunteer,
		Education:         education,
		Awards:            awards,
		Certificates:      certificates,
		Publications:      publications,
		Skills:            skills,
		Languages:         languages,
		Interests:         interests,
		References:        references,
		Projects:          projects,
	}

	// List all resources (both individual and slices)
	allResourcesRaw := []interface{}{
		namespace,
		resume,
		work,
		volunteer,
		education,
		awards,
		certificates,
		publications,
		skills,
		languages,
		interests,
		references,
		projects,
	}

	var resources []models.Resource
	for _, r := range allResourcesRaw {
		val := reflect.ValueOf(r)

		if val.Kind() == reflect.Slice {
			for i := 0; i < val.Len(); i++ {
				item := val.Index(i).Addr().Interface()
				if res, ok := item.(models.Resource); ok {
					resources = append(resources, res)
				}
			}
		} else {
			if res, ok := r.(models.Resource); ok {
				resources = append(resources, res)
			}
		}
	}

	for _, resource := range resources {
		_, err := repo.Create(resource)
		if err != nil {
			panic(err)
		}
	}
}
