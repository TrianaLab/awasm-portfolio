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
		OwnerRef:  ownerRef,
		Name:      "eduardo-diaz",
		Namespace: namespace.Name,
		FullName:  "Eduardo Díaz",
		Label:     "Software Infrastructure Engineer",
		Email:     "edudiazasencio@gmail.com",
		Url:       "https://edudiaz.dev",
		Phone:     "+34 622287557",
		Summary:   "I'm a Software Infrastructure Engineer with strong interest in open-source and cloud-native technologies. I design and build scalable, robust infrastructures and actively contribute to projects that drive innovation. Most of my work and open-source contributions are managed through my personal organization, TrianaLab.",
		Location: types.Location{
			PostalCode:  "41010",
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
			Name:              "mlops-emergenceai",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Position:          "Machine Learning Operations Engineer",
			URL:               "https://emergence.ai",
			StartDate:         "2024-07",
			EndDate:           "",
			Summary:           "As an MLOps Engineer at Emergence, I’m actively involved in building scalable and efficient AI infrastructure utilizing technologies like Terraform, Crossplane, Prometheus, Istio, Keda, Kyverno, etc. My responsibilities include provisioning GKE clusters, orchestrating additional infrastructure components, ensuring robust observability and implementing advanced networking capabilities.",
		},
		{
			OwnerRef:          ownerRef,
			Name:              "prodeng-appian",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Position:          "Product Software Engineer - Kubernetes Team",
			URL:               "https://appian.com",
			StartDate:         "2024-02",
			EndDate:           "2024-07",
			Summary:           "Provide services that support elastic scale and allow frequent, reliable, high-impact changes to the deployed products. Reduce friction and toil surrounding data when creating new product services and features, including data lifecycle management, data retention, data analytics and providing easy-to-use APIs. Make Appian more Kubernetes-native both in cloud and self-managed environments.",
		},
		{
			OwnerRef:          ownerRef,
			Name:              "ssoleng-appian",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Position:          "Senior Solution Engineer – Infrastructure Team",
			URL:               "https://appian.com",
			StartDate:         "2023-10",
			EndDate:           "2024-02",
			Summary:           "Provide technical support globally, address critical challenges and mentor newcomers. With a focus on data analysis and creative solutions, my role emphasizes effective troubleshooting and a comprehensive understanding of the platform’s inner infrastructure.",
		},
		{
			OwnerRef:          ownerRef,
			Name:              "soleng-appian",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Position:          "Solution Engineer – Infrastructure Team",
			URL:               "https://appian.com",
			StartDate:         "2022-10",
			EndDate:           "2023-10",
			Summary:           "Hands-on support for global customers utilizing the Appian platform. With strong problem-solving skills, proficiency in Kubernetes and expertise in web services, programming, and Linux, I contribute to ongoing customer relationships by delivering effective solutions.",
		},
		{
			OwnerRef:          ownerRef,
			Name:              "asoleng-appian",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Position:          "Associate Solution Engineer",
			URL:               "https://appian.com",
			StartDate:         "2021-11",
			EndDate:           "2022-10",
			Summary:           "Tackle complex technical challenges, providing creative solutions and offering world-class support to customers globally. Proficient in troubleshooting, data analytics, and collaboration with internal teams, contribute to resolve Appian installations for both self-managed and Appian Cloud environments.",
		},
		{
			OwnerRef:          ownerRef,
			Name:              "qaeng-solera",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Position:          "Associate Solution Engineer",
			URL:               "https://solera.com",
			StartDate:         "2020-08",
			EndDate:           "2021-11",
			Summary:           "Interface with developers and system architects to ensure applications are designed to be testable while ensuring tags, object ID’s, component and page name standards are in place. Create test plans and test cases based on defined stories. Automate those test cases and incorporate them to correspondent test suites.",
		},
	}

	education := []types.Education{
		{
			OwnerRef:          ownerRef,
			Name:              "miar-viu",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Institution:       "Universidad Internacional de Valencia",
			URL:               "https://www.universidadviu.com/es/master-inteligencia-artificial",
			Area:              "Artificial Intelligence",
			StudyType:         "M.Eng",
			StartDate:         "2021-09",
			EndDate:           "2022-06",
			Courses: []string{
				"Python Fundamentals",
				"Mathematics applied to Artificial Intelligence",
				"Optimization Algorithms",
				"Supervised Learning",
				"Fuzzy Reasoning",
				"Unsupervised Learning",
				"Neural Networks and Deep Learning",
				"Reinforcement Learning",
				"Master's Thesis",
			},
		},
		{
			OwnerRef:          ownerRef,
			Name:              "gitt-us",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Institution:       "Universidad de Sevilla",
			URL:               "https://www.universidadviu.com/es/master-inteligencia-artificial",
			Area:              "Telecommunications",
			StudyType:         "B.Eng",
			StartDate:         "2021-09",
			EndDate:           "2022-06",
			Courses: []string{
				"Physics",
				"Mathematics",
				"Statistics",
				"Theory of Circuits",
				"Electronics",
				"Operating Systems",
				"Telecommunication Network Management",
				"Software Engineering",
				"Security",
				"Traffic Engineering",
				"Advanced Network Architecture",
				"Database Design",
				"Network Planning and Simulation",
				"Telematics Projects",
				"Advanced Telematic Services",
			},
		},
	}

	volunteer := []types.Volunteer{
		{
			OwnerRef:          ownerRef,
			Name:              "jesse-chart",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Organization:      "Jesse",
			Position:          "Open-source contributor",
			URL:               "https://jesse.trade/",
			StartDate:         "2024-12",
			EndDate:           "",
			Summary:           "Created a Kubernetes helm chart for the Jesse AI trading bot.",
		},
		{
			OwnerRef:          ownerRef,
			Name:              "trianalab-awasmportfolio",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Organization:      "TrianaLab",
			Position:          "Owner",
			URL:               "https://edudiaz.dev",
			StartDate:         "2025-01",
			EndDate:           "",
			Summary:           "A webassembly powered application that shows developer’s portoflio using kubernetes-like commands.",
		},
		{
			OwnerRef:          ownerRef,
			Name:              "keda-docs",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Organization:      "Keda",
			Position:          "Open-source contributor",
			URL:               "https://github.com/kedacore/keda-docs/commits/main/?author=edu-diaz",
			StartDate:         "2025-02",
			EndDate:           "2025-02",
			Summary:           "Improved GCP PubSub scaler documentation from KEDA, a Kubernetes-based event driven autoscaler from the Cloud Native Computing Foundation.",
		},
		{
			OwnerRef:          ownerRef,
			Name:              "cloudtty-bug",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Organization:      "CloudTTY",
			Position:          "Open-source contributor",
			URL:               "https://github.com/cloudtty/cloudtty/commits/main/?author=edu-diaz",
			StartDate:         "2024-02",
			EndDate:           "2024-02",
			Summary:           "Fixed a bug from CloudTTY, a Kubernetes cloudshell operator from the Cloud Native Computing Foundation.",
		},
		{
			OwnerRef:          ownerRef,
			Name:              "spin-docs",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Organization:      "Fermyon",
			Position:          "Open-source contributor",
			URL:               "https://github.com/fermyon/developer/commits/main/?author=edu-diaz",
			StartDate:         "2024-11",
			EndDate:           "2024-11",
			Summary:           "Improved documenation from Spin, a developer tool for building WebAssembly microservices and web applications from Fermyon.",
		},
		{
			OwnerRef:          ownerRef,
			Name:              "container2wasm-feat",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Organization:      "container2wasm",
			Position:          "Open-source contributor",
			URL:               "https://github.com/container2wasm/container2wasm/commits/main/?author=edu-diaz",
			StartDate:         "2025-03",
			EndDate:           "2025-03",
			Summary:           "Added Apple Silicon build support to container2wasm, a container to WASM image converter that enables to run the container on web assembly.",
		},
	}

	certificates := []types.Certificate{
		{
			OwnerRef:          ownerRef,
			Name:              "tlf-cka",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Certificate:       "CKA: Certified Kubernetes Administrator",
			Date:              "2022-08",
			Issuer:            "The Linux Foundation",
			URL:               "https://www.credly.com/badges/f1c5619d-f6a1-4988-8393-5f9a21455736/linked_in_profile",
		},
		{
			OwnerRef:          ownerRef,
			Name:              "tlf-cks",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Certificate:       "CKS: Certified Kubernetes Security Specialist",
			Date:              "2023-10",
			Issuer:            "The Linux Foundation",
			URL:               "https://www.credly.com/badges/9e2a89df-4283-4502-9834-7b11b05bb152/linked_in_profile",
		},
		{
			OwnerRef:          ownerRef,
			Name:              "confluent-cf",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Certificate:       "Confluent Fundamentals Accreditation",
			Date:              "2023-04",
			Issuer:            "Confluent",
			URL:               "https://www.credential.net/901b00ba-1188-4eb9-9e24-38d2ee067166#acc.61VbzwAd",
		},
	}

	awards := []types.Award{
		{
			OwnerRef:          ownerRef,
			Name:              "impact-appian",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Title:             "Impact Award",
			Date:              "2022",
			Awarder:           "Appian Corporation",
		},
	}

	publications := []types.Publication{
		{
			OwnerRef:          ownerRef,
			Name:              "trianalab-homelab",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Publication:       "Step by step guide: K3s on Raspberry Pi 5 cluster",
			Publisher:         "TrianaLab",
			ReleaseDate:       "2024-01",
			URL:               "https://edudiaz.trianalab.net/2024/01/12/step-by-step-guide-k3s-on-raspberry-pi-5-cluster/",
			Summary:           "Guide to create a 3 raspberry pi K3s kubernetes cluster in your home to host your workloads.",
		},
	}

	skills := []types.Skill{
		{
			OwnerRef:          ownerRef,
			Name:              "devops-tools",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Skill:             "DevOps Tools",
			Level:             "Expert",
			Keywords:          []string{"Jenkins", "GitHub Actions", "ArgoCD", "Crossplane", "Ansible", "Helm"},
		},
		{
			OwnerRef:          ownerRef,
			Name:              "cloud-platforms",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Skill:             "Cloud Platforms",
			Level:             "Expert",
			Keywords:          []string{"GCP", "AWS", "Cloudflare"},
		},
		{
			OwnerRef:          ownerRef,
			Name:              "containerization",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Skill:             "Containerization and Orchestration",
			Level:             "Expert",
			Keywords:          []string{"Kubernetes", "Helm", "Docker", "Docker Compose", "Operators"},
		},
		{
			OwnerRef:          ownerRef,
			Name:              "service-mesh-observability",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Skill:             "Service Mesh and Observability",
			Level:             "Advanced",
			Keywords:          []string{"Istio", "Prometheus", "Grafana", "Open Telemetry"},
		},
		{
			OwnerRef:          ownerRef,
			Name:              "programming-languages",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Skill:             "Programming Languages",
			Level:             "Advanced",
			Keywords:          []string{"Go", "Java", "Python", "C", "Bash", "WebAssembly (Wasm)"},
		},
		{
			OwnerRef:          ownerRef,
			Name:              "networking",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Skill:             "Networking",
			Level:             "Expert",
			Keywords:          []string{"Networking Fundamentals", "Network Security"},
		},
		{
			OwnerRef:          ownerRef,
			Name:              "other-skills",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Skill:             "Other Skills",
			Level:             "Expert",
			Keywords:          []string{"Observability", "Infrastructure as Code", "CI/CD Pipelines", "Microservices Architecture"},
		},
	}

	languages := []types.Language{
		{
			OwnerRef:          ownerRef,
			Name:              "spanish",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Language:          "Spanish",
			Fluency:           "Native",
		},
		{
			OwnerRef:          ownerRef,
			Name:              "english",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Language:          "English",
			Fluency:           "Advanced - C1",
		},
		{
			OwnerRef:          ownerRef,
			Name:              "chinese",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Language:          "Mandarin",
			Fluency:           "Basic - HSK2",
		},
	}

	interests := []types.Interest{
		{
			OwnerRef:          ownerRef,
			Name:              "cloud-native",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Interest:          "Cloud Native",
			Keywords:          []string{"Open Source", "Kubernetes", "WebAssembly", "Microservices"},
		},
		{
			OwnerRef:          ownerRef,
			Name:              "mountaineering",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Interest:          "Mountaineering",
			Keywords:          []string{"Mountains", "Hiking", "Trekking"},
		},
		{
			OwnerRef:          ownerRef,
			Name:              "lego-architecture",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Interest:          "Legos",
			Keywords:          []string{"Lego", "Architectural design", "Creative builds"},
		},
		{
			OwnerRef:          ownerRef,
			Name:              "culinary-adventures",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Interest:          "Gastronomy",
			Keywords:          []string{"Cooking", "Pucherito", "Mediterranean gastronomy"},
		},
	}

	references := []types.Reference{
		{
			OwnerRef:          ownerRef,
			Name:              "john-doe",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Person:            "Pepito",
			Reference:         "Viva er beti.",
		},
	}

	projects := []types.Project{
		{
			OwnerRef:          ownerRef,
			Name:              "trianalab",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Project:           "TrianaLab",
			Description:       "TrianaLab is my personal GitHub organization where I create my open source projects and contributions.",
			StartDate:         "2024-12",
			EndDate:           "",
			URL:               "https://github.com/TrianaLab",
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
		basics,
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
