package preload

import (
	"github.com/TrianaLab/awasm-portfolio/internal/models"
	"github.com/TrianaLab/awasm-portfolio/internal/models/types"
	"github.com/TrianaLab/awasm-portfolio/internal/repository"
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
		Label:     "Platform Engineer",
		Email:     "edudiazasencio@gmail.com",
		Url:       "https://edudiaz.dev",
		Phone:     "+34 622287557",
		Summary:   "Platform Engineer focused on cloud-native infrastructure, WebAssembly, and operating-systems development. I design and operate the platforms — Kubernetes, observability, networking, CI/CD — that other engineering teams build on, with an emphasis on reliability, developer experience, and long-term maintainability. I sustain a steady cadence of upstream open-source contributions (KEDA, Docker, Artifact Hub, container2wasm, Spin, CloudTTY) and run TrianaLab, where I design and ship original cloud-native tooling (pacto, awasm-portfolio, remake).",
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
			Position:          "Platform Engineer",
			URL:               "https://emergence.ai",
			StartDate:         "2024-07-29",
			Summary:           "Design and operate the AI platform on GKE. Lead workstreams across cluster provisioning (Terraform, Crossplane), service mesh and ingress (Istio), event-driven autoscaling (KEDA), policy and admission control (Kyverno), and end-to-end observability (Prometheus, OpenTelemetry). Focus on making the platform safe, scalable, and self-service for the AI engineering teams that depend on it.",
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
			Summary:           "Built platform services on Kubernetes that let product teams ship elastic, high-impact changes safely across cloud and self-managed environments. Owned data-platform primitives (lifecycle, retention, analytics APIs) that reduced toil for downstream service teams, and drove the broader push to make Appian more Kubernetes-native.",
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
			Summary:           "Senior engineer on the team handling escalated Kubernetes and cloud incidents for enterprise customers. Diagnosed and resolved complex platform issues across observability, networking, and automation, and fed recurring failure modes back into platform-level improvements.",
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
			Summary:           "Supported enterprise customers running Appian on Kubernetes across cloud and on-prem deployments. Owned observability tuning, performance investigations, and cluster-level troubleshooting, partnering with product and engineering to harden infrastructure reliability and delivery.",
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
			Summary:           "Joined the infrastructure support team focused on networking, storage, and observability for Appian's self-managed Kubernetes deployments. Built deep operational fluency in containerized environments and customer-incident response.",
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
			Summary:           "Partnered with developers and architects to bake testability into applications at design time, defining and enforcing tagging, naming, and component standards. Authored test plans, automated regression suites, and integrated them into the team's CI workflows.",
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
			URL:               "https://etsi.us.es/en/studies-and-qualifications/degrees/degree-in-telecommunications-technology-engineering",
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
			Name:              "volunteer-trianalab-pacto",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Organization:      "TrianaLab: pacto",
			Position:          "Owner",
			URL:               "https://github.com/TrianaLab/pacto",
			StartDate:         "2026-03-03",
			Summary:           "An OCI-distributed contract system for cloud-native services. Pacto pairs a CLI, dashboard, and Kubernetes operator so teams can describe a service's operational behavior once — interfaces, dependencies, runtime semantics, configuration, scaling — then validate it, diff it, distribute it via OCI registries, and verify alignment against running workloads.",
		},
		{
			Kind:              "volunteer",
			OwnerRef:          ownerRef,
			Name:              "volunteer-remake",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Organization:      "TrianaLab: remake",
			Position:          "Owner",
			URL:               "https://github.com/TrianaLab/remake",
			StartDate:         "2025-05-23",
			Summary:           "A lightweight CLI tool that lets you package and share Makefiles as OCI artifacts.",
		},
		{
			Kind:              "volunteer",
			OwnerRef:          ownerRef,
			Name:              "volunteer-trianalab-awasmportfolio",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Organization:      "TrianaLab: awasm-portfolio",
			Position:          "Owner",
			URL:               "https://edudiaz.dev",
			StartDate:         "2025-01-19",
			Summary:           "A WebAssembly-powered application that exposes the developer's portfolio through kubectl-style commands.",
		},
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
			Summary:           "Created a Kubernetes Helm chart for the Jesse AI trading bot.",
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
			Summary:           "Improved the GCP Pub/Sub scaler documentation for KEDA, the CNCF Kubernetes event-driven autoscaler.",
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
			Summary:           "Fixed a bug in CloudTTY, the CNCF Kubernetes cloud-shell operator.",
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
			Summary:           "Improved documentation for Spin, Fermyon's developer tool for building WebAssembly microservices and web applications.",
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
			Summary:           "Added Apple Silicon build support to container2wasm, a tool that converts container images into WebAssembly so they can run in the browser.",
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
			Summary:           "Improved the Docker Compose OCI artifacts documentation.",
		},
		{
			Kind:              "volunteer",
			OwnerRef:          ownerRef,
			Name:              "volunteer-artifacthub-mermaid",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Organization:      "Artifact Hub",
			Position:          "Open-source contributor",
			URL:               "https://github.com/artifacthub/hub/commits/main/?author=edu-diaz",
			StartDate:         "2026-04-13",
			Summary:           "Added Mermaid diagram rendering to package README files in Artifact Hub, a CNCF project for finding, installing, and publishing cloud-native packages.",
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
			Name:              "interest-open-source",
			Namespace:         namespace.Name,
			CreationTimestamp: timestamp,
			Interest:          "Cloud Native",
			Keywords:          []string{"Open Source", "Kubernetes", "WebAssembly", "OSDev"},
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
			Keywords:          []string{"Cooking", "Mediterranean gastronomy"},
		},
	}

	resume := &types.Resume{
		Kind:              "resume",
		Name:              "main-resume",
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
