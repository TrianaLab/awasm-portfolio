// Package preload seeds the in-memory repository with the portfolio's
// static content. To customize the portfolio, edit the slices below.
package preload

import (
	"time"

	"github.com/TrianaLab/awasm-portfolio/internal/models"
	"github.com/TrianaLab/awasm-portfolio/internal/models/types"
	"github.com/TrianaLab/awasm-portfolio/internal/repository"
)

const (
	namespaceName = "default"
	ownerName     = "eduardo-diaz"
)

// PreloadData populates repo with the namespace, the aggregate Resume,
// and every child resource the resume aggregates.
func PreloadData(repo *repository.InMemoryRepository) {
	now := time.Now()
	owner := models.OwnerReference{Kind: "resume", Name: ownerName, Namespace: namespaceName}
	mk := func(kind, name string) models.Meta {
		return models.Meta{
			Kind:              kind,
			Name:              name,
			Namespace:         namespaceName,
			OwnerRef:          owner,
			CreationTimestamp: now,
		}
	}

	namespace := &types.Namespace{Meta: models.Meta{Kind: "namespace", Name: namespaceName}}
	basics := buildBasics(mk)
	work := buildWork(mk)
	education := buildEducation(mk)
	volunteer := buildVolunteer(mk)
	certificates := buildCertificates(mk)
	skills := buildSkills(mk)
	languages := buildLanguages(mk)
	interests := buildInterests(mk)

	resume := &types.Resume{
		Meta: models.Meta{
			Kind:              "resume",
			Name:              "main-resume",
			Namespace:         namespaceName,
			CreationTimestamp: now,
		},
		Basics:       *basics,
		Work:         work,
		Volunteer:    volunteer,
		Education:    education,
		Certificates: certificates,
		Skills:       skills,
		Languages:    languages,
		Interests:    interests,
	}

	resources := []models.Resource{namespace, resume, basics}
	for i := range work {
		resources = append(resources, &work[i])
	}
	for i := range volunteer {
		resources = append(resources, &volunteer[i])
	}
	for i := range education {
		resources = append(resources, &education[i])
	}
	for i := range certificates {
		resources = append(resources, &certificates[i])
	}
	for i := range skills {
		resources = append(resources, &skills[i])
	}
	for i := range languages {
		resources = append(resources, &languages[i])
	}
	for i := range interests {
		resources = append(resources, &interests[i])
	}

	for _, resource := range resources {
		if _, err := repo.Create(resource); err != nil {
			panic(err)
		}
	}
}

func buildBasics(mk func(kind, name string) models.Meta) *types.Basics {
	return &types.Basics{
		Meta:     mk("basics", "basics-eduardo-diaz"),
		FullName: "Eduardo Díaz",
		Label:    "Platform Engineer",
		Email:    "edudiazasencio@gmail.com",
		Url:      "https://edudiaz.dev",
		Phone:    "+34 622287557",
		Summary:  "Platform Engineer focused on cloud-native infrastructure, WebAssembly, and operating-systems development. I design and operate the platforms — Kubernetes, observability, networking, CI/CD — that other engineering teams build on, with an emphasis on reliability, developer experience, and long-term maintainability. I sustain a steady cadence of upstream open-source contributions (KEDA, Docker, Artifact Hub, container2wasm, Spin, CloudTTY) and run TrianaLab, where I design and ship original cloud-native tooling (pacto, awasm-portfolio, remake).",
		Location: types.Location{
			PostalCode:  "41010",
			City:        "Sevilla",
			Region:      "Andalucía",
			CountryCode: "ES",
		},
		Profiles: []types.Profile{
			{Network: "LinkedIn", Username: "eduardo-diaz-asencio", Url: "https://www.linkedin.com/in/eduardo-diaz-asencio/"},
			{Network: "GitHub", Username: "edu-diaz", Url: "https://github.com/edu-diaz"},
		},
	}
}

func buildWork(mk func(kind, name string) models.Meta) []types.Work {
	return []types.Work{
		{
			Meta:      mk("work", "work-mlops-emergenceai"),
			Company:   "Emergence AI",
			Position:  "Platform Engineer",
			URL:       "https://emergence.ai",
			StartDate: "2024-07-29",
			Summary:   "Design and operate the AI platform on GKE. Lead workstreams across cluster provisioning (Terraform, Crossplane), service mesh and ingress (Istio), event-driven autoscaling (KEDA), policy and admission control (Kyverno), and end-to-end observability (Prometheus, OpenTelemetry). Focus on making the platform safe, scalable, and self-service for the AI engineering teams that depend on it.",
		},
		{
			Meta:      mk("work", "work-prodeng-appian"),
			Company:   "Appian Corporation",
			Position:  "Product Software Engineer - Kubernetes Team",
			URL:       "https://appian.com",
			StartDate: "2024-02-01",
			EndDate:   "2024-07-29",
			Summary:   "Built platform services on Kubernetes that let product teams ship elastic, high-impact changes safely across cloud and self-managed environments. Owned data-platform primitives (lifecycle, retention, analytics APIs) that reduced toil for downstream service teams, and drove the broader push to make Appian more Kubernetes-native.",
		},
		{
			Meta:      mk("work", "work-ssoleng-appian"),
			Company:   "Appian Corporation",
			Position:  "Senior Solution Engineer – Infrastructure Team",
			URL:       "https://appian.com",
			StartDate: "2023-10-01",
			EndDate:   "2024-02-01",
			Summary:   "Senior engineer on the team handling escalated Kubernetes and cloud incidents for enterprise customers. Diagnosed and resolved complex platform issues across observability, networking, and automation, and fed recurring failure modes back into platform-level improvements.",
		},
		{
			Meta:      mk("work", "work-soleng-appian"),
			Company:   "Appian Corporation",
			Position:  "Solution Engineer – Infrastructure Team",
			URL:       "https://appian.com",
			StartDate: "2022-10-01",
			EndDate:   "2023-10-01",
			Summary:   "Supported enterprise customers running Appian on Kubernetes across cloud and on-prem deployments. Owned observability tuning, performance investigations, and cluster-level troubleshooting, partnering with product and engineering to harden infrastructure reliability and delivery.",
		},
		{
			Meta:      mk("work", "work-asoleng-appian"),
			Company:   "Appian Corporation",
			Position:  "Associate Solution Engineer",
			URL:       "https://appian.com",
			StartDate: "2021-11-01",
			EndDate:   "2022-10-01",
			Summary:   "Joined the infrastructure support team focused on networking, storage, and observability for Appian's self-managed Kubernetes deployments. Built deep operational fluency in containerized environments and customer-incident response.",
		},
		{
			Meta:      mk("work", "work-qaeng-solera"),
			Company:   "Solera Holdings",
			Position:  "Software QA Engineer",
			URL:       "https://solera.com",
			StartDate: "2020-08-01",
			EndDate:   "2021-11-01",
			Summary:   "Partnered with developers and architects to bake testability into applications at design time, defining and enforcing tagging, naming, and component standards. Authored test plans, automated regression suites, and integrated them into the team's CI workflows.",
		},
	}
}

func buildEducation(mk func(kind, name string) models.Meta) []types.Education {
	return []types.Education{
		{
			Meta:        mk("education", "education-miar-viu"),
			Institution: "Universidad Internacional de Valencia",
			URL:         "https://www.universidadviu.com/es/master-inteligencia-artificial",
			Area:        "Artificial Intelligence",
			StudyType:   "M.Eng",
			StartDate:   "2021-09-01",
			EndDate:     "2022-06-30",
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
			Meta:        mk("education", "education-gitt-us"),
			Institution: "Universidad de Sevilla",
			URL:         "https://etsi.us.es/en/studies-and-qualifications/degrees/degree-in-telecommunications-technology-engineering",
			Area:        "Telecommunications",
			StudyType:   "B.Eng",
			StartDate:   "2021-09-01",
			EndDate:     "2022-06-01",
			Courses: []string{
				"Physics", "Mathematics", "Statistics", "Theory of Circuits",
				"Electronics", "Operating Systems", "Telecommunication Network Management",
				"Software Engineering", "Security", "Traffic Engineering",
				"Advanced Network Architecture", "Database Design",
				"Network Planning and Simulation", "Telematics Projects",
				"Advanced Telematic Services",
			},
		},
	}
}

func buildVolunteer(mk func(kind, name string) models.Meta) []types.Volunteer {
	return []types.Volunteer{
		{
			Meta:         mk("volunteer", "volunteer-trianalab-pacto"),
			Organization: "TrianaLab: pacto",
			Position:     "Owner",
			URL:          "https://github.com/TrianaLab/pacto",
			StartDate:    "2026-03-03",
			Summary:      "An OCI-distributed contract system for cloud-native services. Pacto pairs a CLI, dashboard, and Kubernetes operator so teams can describe a service's operational behavior once — interfaces, dependencies, runtime semantics, configuration, scaling — then validate it, diff it, distribute it via OCI registries, and verify alignment against running workloads.",
		},
		{
			Meta:         mk("volunteer", "volunteer-remake"),
			Organization: "TrianaLab: remake",
			Position:     "Owner",
			URL:          "https://github.com/TrianaLab/remake",
			StartDate:    "2025-05-23",
			Summary:      "A lightweight CLI tool that lets you package and share Makefiles as OCI artifacts.",
		},
		{
			Meta:         mk("volunteer", "volunteer-trianalab-awasmportfolio"),
			Organization: "TrianaLab: awasm-portfolio",
			Position:     "Owner",
			URL:          "https://edudiaz.dev",
			StartDate:    "2025-01-19",
			Summary:      "A WebAssembly-powered application that exposes the developer's portfolio through kubectl-style commands.",
		},
		{
			Meta:         mk("volunteer", "volunteer-jesse-chart"),
			Organization: "Jesse",
			Position:     "Open-source contributor",
			URL:          "https://jesse.trade/",
			StartDate:    "2024-12-29",
			Summary:      "Created a Kubernetes Helm chart for the Jesse AI trading bot.",
		},
		{
			Meta:         mk("volunteer", "volunteer-keda-docs"),
			Organization: "Keda",
			Position:     "Open-source contributor",
			URL:          "https://github.com/kedacore/keda-docs/commits/main/?author=edu-diaz",
			StartDate:    "2025-02-06",
			Summary:      "Improved the GCP Pub/Sub scaler documentation for KEDA, the CNCF Kubernetes event-driven autoscaler.",
		},
		{
			Meta:         mk("volunteer", "volunteer-cloudtty-bug"),
			Organization: "CloudTTY",
			Position:     "Open-source contributor",
			URL:          "https://github.com/cloudtty/cloudtty/commits/main/?author=edu-diaz",
			StartDate:    "2024-02-01",
			Summary:      "Fixed a bug in CloudTTY, the CNCF Kubernetes cloud-shell operator.",
		},
		{
			Meta:         mk("volunteer", "volunteer-spin-docs"),
			Organization: "Fermyon",
			Position:     "Open-source contributor",
			URL:          "https://github.com/fermyon/developer/commits/main/?author=edu-diaz",
			StartDate:    "2024-11-27",
			Summary:      "Improved documentation for Spin, Fermyon's developer tool for building WebAssembly microservices and web applications.",
		},
		{
			Meta:         mk("volunteer", "volunteer-container2wasm-feat"),
			Organization: "container2wasm",
			Position:     "Open-source contributor",
			URL:          "https://github.com/container2wasm/container2wasm/commits/main/?author=edu-diaz",
			StartDate:    "2025-03-08",
			Summary:      "Added Apple Silicon build support to container2wasm, a tool that converts container images into WebAssembly so they can run in the browser.",
		},
		{
			Meta:         mk("volunteer", "volunteer-docker-docs"),
			Organization: "Docker",
			Position:     "Open-source contributor",
			URL:          "https://github.com/docker/docs/commits/main/?author=edu-diaz",
			StartDate:    "2025-05-01",
			Summary:      "Improved the Docker Compose OCI artifacts documentation.",
		},
		{
			Meta:         mk("volunteer", "volunteer-artifacthub-mermaid"),
			Organization: "Artifact Hub",
			Position:     "Open-source contributor",
			URL:          "https://github.com/artifacthub/hub/commits/master/?author=edu-diaz",
			StartDate:    "2026-04-13",
			Summary:      "Added Mermaid diagram rendering to package README files in Artifact Hub, a CNCF project for finding, installing, and publishing cloud-native packages.",
		},
	}
}

func buildCertificates(mk func(kind, name string) models.Meta) []types.Certificate {
	return []types.Certificate{
		{
			Meta:        mk("certificate", "certificate-tlf-cka"),
			Certificate: "CKA: Certified Kubernetes Administrator",
			Date:        "2022-08-26",
			Issuer:      "The Linux Foundation",
			URL:         "https://www.credly.com/badges/f1c5619d-f6a1-4988-8393-5f9a21455736/linked_in_profile",
		},
		{
			Meta:        mk("certificate", "certificate-tlf-cks"),
			Certificate: "CKS: Certified Kubernetes Security Specialist",
			Date:        "2023-10-06",
			Issuer:      "The Linux Foundation",
			URL:         "https://www.credly.com/badges/9e2a89df-4283-4502-9834-7b11b05bb152/linked_in_profile",
		},
		{
			Meta:        mk("certificate", "certificate-confluent-cf"),
			Certificate: "Confluent Fundamentals Accreditation",
			Date:        "2023-04-03",
			Issuer:      "Confluent",
			URL:         "https://www.credential.net/901b00ba-1188-4eb9-9e24-38d2ee067166#acc.61VbzwAd",
		},
	}
}

func buildSkills(mk func(kind, name string) models.Meta) []types.Skill {
	return []types.Skill{
		{
			Meta:     mk("skill", "skill-devops-tools"),
			Skill:    "DevOps Tools",
			Level:    "Expert",
			Keywords: []string{"Jenkins", "GitHub Actions", "ArgoCD", "Crossplane", "Ansible", "Helm"},
		},
		{
			Meta:     mk("skill", "skill-cloud-platforms"),
			Skill:    "Cloud Platforms",
			Level:    "Expert",
			Keywords: []string{"GCP", "AWS", "Cloudflare"},
		},
		{
			Meta:     mk("skill", "skill-containerization-technologies"),
			Skill:    "Containerization and Orchestration",
			Level:    "Expert",
			Keywords: []string{"Kubernetes", "Helm", "Docker", "Docker Compose", "Operators"},
		},
		{
			Meta:     mk("skill", "skill-service-mesh-observability"),
			Skill:    "Service Mesh and Observability",
			Level:    "Advanced",
			Keywords: []string{"Istio", "Prometheus", "Grafana", "Open Telemetry"},
		},
		{
			Meta:     mk("skill", "skill-programming-languages"),
			Skill:    "Programming Languages",
			Level:    "Advanced",
			Keywords: []string{"Go", "Java", "Python", "C", "Bash", "WebAssembly (Wasm)"},
		},
		{
			Meta:     mk("skill", "skill-networking"),
			Skill:    "Networking",
			Level:    "Expert",
			Keywords: []string{"Networking Fundamentals", "Network Security"},
		},
		{
			Meta:     mk("skill", "skill-other-skills"),
			Skill:    "Other Skills",
			Level:    "Expert",
			Keywords: []string{"Observability", "Infrastructure as Code", "CI/CD Pipelines", "Microservices Architecture"},
		},
	}
}

func buildLanguages(mk func(kind, name string) models.Meta) []types.Language {
	return []types.Language{
		{Meta: mk("language", "language-spanish"), Language: "Spanish", Fluency: "Native"},
		{Meta: mk("language", "language-english"), Language: "English", Fluency: "Advanced - C1"},
		{Meta: mk("language", "language-chinese"), Language: "Mandarin", Fluency: "Basic - HSK2"},
	}
}

func buildInterests(mk func(kind, name string) models.Meta) []types.Interest {
	return []types.Interest{
		{
			Meta:     mk("interest", "interest-open-source"),
			Interest: "Cloud Native",
			Keywords: []string{"Open Source", "Kubernetes", "WebAssembly", "OSDev"},
		},
		{
			Meta:     mk("interest", "interest-mountaineering"),
			Interest: "Mountaineering",
			Keywords: []string{"Mountains", "Hiking", "Trekking"},
		},
		{
			Meta:     mk("interest", "interest-lego-architecture"),
			Interest: "Legos",
			Keywords: []string{"Lego", "Architectural design", "Creative builds"},
		},
		{
			Meta:     mk("interest", "interest-culinary-adventures"),
			Interest: "Gastronomy",
			Keywords: []string{"Cooking", "Mediterranean gastronomy"},
		},
	}
}
