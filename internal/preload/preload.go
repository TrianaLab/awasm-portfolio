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
		Kind:      "basics",
		OwnerRef:  ownerRef,
		Name:      "basics-eduardo-diaz",
		Namespace: namespace.Name,
		FullName:  "Eduardo Díaz",
		Label:     "Software Infrastructure Engineer",
		Email:     "edudiazasencio@gmail.com",
		Url:       "https://edudiaz.dev",
		Phone:     "622287557",
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
			Kind:              "work",
			OwnerRef:          ownerRef,
			Name:              "work-mlops-emergenceai",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Company:           "Emergence AI",
			Position:          "Machine Learning Operations Engineer",
			URL:               "https://emergence.ai",
			StartDate:         "2024-07-29",
			Summary:           "As an MLOps Engineer at Emergence, I’m actively involved in building scalable and efficient AI infrastructure utilizing technologies like Terraform, Crossplane, Prometheus, Istio, Keda, Kyverno, etc. My responsibilities include provisioning GKE clusters, orchestrating additional infrastructure components, ensuring robust observability and implementing advanced networking capabilities.",
		},
		{
			Kind:              "work",
			OwnerRef:          ownerRef,
			Name:              "work-prodeng-appian",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Company:           "Appian Corporation",
			Position:          "Product Software Engineer - Kubernetes Team",
			URL:               "https://appian.com",
			StartDate:         "2024-02-01",
			EndDate:           "2024-07-29",
			Summary:           "Provide services that support elastic scale and allow frequent, reliable, high-impact changes to the deployed products. Reduce friction and toil surrounding data when creating new product services and features, including data lifecycle management, data retention, data analytics and providing easy-to-use APIs. Make Appian more Kubernetes-native both in cloud and self-managed environments.",
		},
		{
			Kind:              "work",
			OwnerRef:          ownerRef,
			Name:              "work-ssoleng-appian",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Company:           "Appian Corporation",
			Position:          "Senior Solution Engineer – Infrastructure Team",
			URL:               "https://appian.com",
			StartDate:         "2023-10-01",
			EndDate:           "2024-02-01",
			Summary:           "Provide technical support globally, address critical challenges and mentor newcomers. With a focus on data analysis and creative solutions, my role emphasizes effective troubleshooting and a comprehensive understanding of the platform’s inner infrastructure.",
		},
		{
			Kind:              "work",
			OwnerRef:          ownerRef,
			Name:              "work-soleng-appian",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Company:           "Appian Corporation",
			Position:          "Solution Engineer – Infrastructure Team",
			URL:               "https://appian.com",
			StartDate:         "2022-10-01",
			EndDate:           "2023-10-01",
			Summary:           "Hands-on support for global customers utilizing the Appian platform. With strong problem-solving skills, proficiency in Kubernetes and expertise in web services, programming, and Linux, I contribute to ongoing customer relationships by delivering effective solutions.",
		},
		{
			Kind:              "work",
			OwnerRef:          ownerRef,
			Name:              "work-asoleng-appian",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Company:           "Appian Corporation",
			Position:          "Associate Solution Engineer",
			URL:               "https://appian.com",
			StartDate:         "2021-11-01",
			EndDate:           "2022-10-01",
			Summary:           "Tackle complex technical challenges, providing creative solutions and offering world-class support to customers globally. Proficient in troubleshooting, data analytics, and collaboration with internal teams, contribute to resolve Appian installations for both self-managed and Appian Cloud environments.",
		},
		{
			Kind:              "work",
			OwnerRef:          ownerRef,
			Name:              "work-qaeng-solera",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Company:           "Solera Holdings",
			Position:          "Software QA Engineer",
			URL:               "https://solera.com",
			StartDate:         "2020-08-01",
			EndDate:           "2021-11-01",
			Summary:           "Interface with developers and system architects to ensure applications are designed to be testable while ensuring tags, object ID’s, component and page name standards are in place. Create test plans and test cases based on defined stories. Automate those test cases and incorporate them to correspondent test suites.",
		},
	}

	education := []types.Education{
		{
			Kind:              "education",
			OwnerRef:          ownerRef,
			Name:              "education-miar-viu",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Institution:       "Universidad Internacional de Valencia",
			URL:               "https://www.universidadviu.com/es/master-inteligencia-artificial",
			Area:              "Artificial Intelligence",
			StudyType:         "M.Eng",
			StartDate:         "2021-09-01",
			EndDate:           "2022-06-30",
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
			Kind:              "education",
			OwnerRef:          ownerRef,
			Name:              "education-gitt-us",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Institution:       "Universidad de Sevilla",
			URL:               "https://www.universidadviu.com/es/master-inteligencia-artificial",
			Area:              "Telecommunications",
			StudyType:         "B.Eng",
			StartDate:         "2021-09-01",
			EndDate:           "2022-06-01",
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
			Kind:              "volunteer",
			OwnerRef:          ownerRef,
			Name:              "volunteer-jesse-chart",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Organization:      "Jesse",
			Position:          "Open-source contributor",
			URL:               "https://jesse.trade/",
			StartDate:         "2024-12-29",
			Summary:           "Created a Kubernetes helm chart for the Jesse AI trading bot.",
		},
		{
			Kind:              "volunteer",
			OwnerRef:          ownerRef,
			Name:              "volunteer-trianalab-awasmportfolio",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Organization:      "TrianaLab",
			Position:          "Owner",
			URL:               "https://edudiaz.dev",
			StartDate:         "2025-01-19",
			Summary:           "A webassembly powered application that shows developer’s portoflio using kubernetes-like commands.",
		},
		{
			Kind:              "volunteer",
			OwnerRef:          ownerRef,
			Name:              "volunteer-keda-docs",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Organization:      "Keda",
			Position:          "Open-source contributor",
			URL:               "https://github.com/kedacore/keda-docs/commits/main/?author=edu-diaz",
			StartDate:         "2025-02-06",
			Summary:           "Improved GCP PubSub scaler documentation from KEDA, a Kubernetes-based event driven autoscaler from the Cloud Native Computing Foundation.",
		},
		{
			Kind:              "volunteer",
			OwnerRef:          ownerRef,
			Name:              "volunteer-cloudtty-bug",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Organization:      "CloudTTY",
			Position:          "Open-source contributor",
			URL:               "https://github.com/cloudtty/cloudtty/commits/main/?author=edu-diaz",
			StartDate:         "2024-02-01",
			Summary:           "Fixed a bug from CloudTTY, a Kubernetes cloudshell operator from the Cloud Native Computing Foundation.",
		},
		{
			Kind:              "volunteer",
			OwnerRef:          ownerRef,
			Name:              "volunteer-spin-docs",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Organization:      "Fermyon",
			Position:          "Open-source contributor",
			URL:               "https://github.com/fermyon/developer/commits/main/?author=edu-diaz",
			StartDate:         "2024-11-27",
			Summary:           "Improved documenation from Spin, a developer tool for building WebAssembly microservices and web applications from Fermyon.",
		},
		{
			Kind:              "volunteer",
			OwnerRef:          ownerRef,
			Name:              "volunteer-container2wasm-feat",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Organization:      "container2wasm",
			Position:          "Open-source contributor",
			URL:               "https://github.com/container2wasm/container2wasm/commits/main/?author=edu-diaz",
			StartDate:         "2025-03-08",
			Summary:           "Added Apple Silicon build support to container2wasm, a container to WASM image converter that enables to run the container on web assembly.",
		},
		{
			Kind:              "volunteer",
			OwnerRef:          ownerRef,
			Name:              "volunteer-docker-docs",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Organization:      "Docker",
			Position:          "Open-source contributor",
			URL:               "https://github.com/docker/docs/commits/main/?author=edu-diaz",
			StartDate:         "2025-05-01",
			Summary:           "Improved docker compose OCI artifacts documentation.",
		},
	}

	certificates := []types.Certificate{
		{
			Kind:              "certificate",
			OwnerRef:          ownerRef,
			Name:              "certificate-tlf-cka",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Certificate:       "CKA: Certified Kubernetes Administrator",
			Date:              "2022-08-26",
			Issuer:            "The Linux Foundation",
			URL:               "https://www.credly.com/badges/f1c5619d-f6a1-4988-8393-5f9a21455736/linked_in_profile",
		},
		{
			Kind:              "certificate",
			OwnerRef:          ownerRef,
			Name:              "certificate-tlf-cks",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Certificate:       "CKS: Certified Kubernetes Security Specialist",
			Date:              "2023-10-06",
			Issuer:            "The Linux Foundation",
			URL:               "https://www.credly.com/badges/9e2a89df-4283-4502-9834-7b11b05bb152/linked_in_profile",
		},
		{
			Kind:              "certificate",
			OwnerRef:          ownerRef,
			Name:              "certificate-confluent-cf",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Certificate:       "Confluent Fundamentals Accreditation",
			Date:              "2023-04-03",
			Issuer:            "Confluent",
			URL:               "https://www.credential.net/901b00ba-1188-4eb9-9e24-38d2ee067166#acc.61VbzwAd",
		},
	}

	skills := []types.Skill{
		{
			Kind:              "skill",
			OwnerRef:          ownerRef,
			Name:              "skill-devops-tools",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Skill:             "DevOps Tools",
			Level:             "Expert",
			Keywords:          []string{"Jenkins", "GitHub Actions", "ArgoCD", "Crossplane", "Ansible", "Helm"},
		},
		{
			Kind:              "skill",
			OwnerRef:          ownerRef,
			Name:              "skill-cloud-platforms",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Skill:             "Cloud Platforms",
			Level:             "Expert",
			Keywords:          []string{"GCP", "AWS", "Cloudflare"},
		},
		{
			Kind:              "skill",
			OwnerRef:          ownerRef,
			Name:              "skill-containerization-technologies",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Skill:             "Containerization and Orchestration",
			Level:             "Expert",
			Keywords:          []string{"Kubernetes", "Helm", "Docker", "Docker Compose", "Operators"},
		},
		{
			Kind:              "skill",
			OwnerRef:          ownerRef,
			Name:              "skill-service-mesh-observability",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Skill:             "Service Mesh and Observability",
			Level:             "Advanced",
			Keywords:          []string{"Istio", "Prometheus", "Grafana", "Open Telemetry"},
		},
		{
			Kind:              "skill",
			OwnerRef:          ownerRef,
			Name:              "skill-programming-languages",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Skill:             "Programming Languages",
			Level:             "Advanced",
			Keywords:          []string{"Go", "Java", "Python", "C", "Bash", "WebAssembly (Wasm)"},
		},
		{
			Kind:              "skill",
			OwnerRef:          ownerRef,
			Name:              "skill-networking",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Skill:             "Networking",
			Level:             "Expert",
			Keywords:          []string{"Networking Fundamentals", "Network Security"},
		},
		{
			Kind:              "skill",
			OwnerRef:          ownerRef,
			Name:              "skill-other-skills",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Skill:             "Other Skills",
			Level:             "Expert",
			Keywords:          []string{"Observability", "Infrastructure as Code", "CI/CD Pipelines", "Microservices Architecture"},
		},
	}

	languages := []types.Language{
		{
			Kind:              "language",
			OwnerRef:          ownerRef,
			Name:              "language-spanish",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Language:          "Spanish",
			Fluency:           "Native",
		},
		{
			Kind:              "language",
			OwnerRef:          ownerRef,
			Name:              "language-english",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Language:          "English",
			Fluency:           "Advanced - C1",
		},
		{
			Kind:              "language",
			OwnerRef:          ownerRef,
			Name:              "language-chinese",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Language:          "Mandarin",
			Fluency:           "Basic - HSK2",
		},
	}

	interests := []types.Interest{
		{
			Kind:              "interest",
			OwnerRef:          ownerRef,
			Name:              "interest-cloud-native",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Interest:          "Cloud Native",
			Keywords:          []string{"Open Source", "Kubernetes", "WebAssembly", "Microservices"},
		},
		{
			Kind:              "interest",
			OwnerRef:          ownerRef,
			Name:              "interest-mountaineering",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Interest:          "Mountaineering",
			Keywords:          []string{"Mountains", "Hiking", "Trekking"},
		},
		{
			Kind:              "interest",
			OwnerRef:          ownerRef,
			Name:              "interest-lego-architecture",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Interest:          "Legos",
			Keywords:          []string{"Lego", "Architectural design", "Creative builds"},
		},
		{
			Kind:              "interest",
			OwnerRef:          ownerRef,
			Name:              "interest-culinary-adventures",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Interest:          "Gastronomy",
			Keywords:          []string{"Cooking", "Pucherito", "Mediterranean gastronomy"},
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
		Certificates:      certificates,
		Skills:            skills,
		Languages:         languages,
		Interests:         interests,
	}

	// List all resources (both individual and slices)
	allResourcesRaw := []interface{}{
		namespace,
		resume,
		basics,
		work,
		volunteer,
		education,
		certificates,
		skills,
		languages,
		interests,
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
