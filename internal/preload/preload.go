package preload

import (
	"awasm-portfolio/internal/models"
	"awasm-portfolio/internal/models/types"
	"awasm-portfolio/internal/repository"
)

func PreloadData(repo *repository.InMemoryRepository) {
	name := "eduardo"
	namespace := &types.Namespace{
		Kind: "namespace",
		Name: name + "-namespace",
	}
	ownerRef := models.OwnerReference{
		Kind:      (&types.Profile{}).GetKind(),
		Name:      name + "-profile",
		Namespace: namespace.Name,
	}
	certifications := &types.Certifications{
		Name:      name + "-certifications",
		Namespace: namespace.Name,
		OwnerRef:  ownerRef,
		Certifications: []types.Certification{
			{
				Description: "Certified Kubernetes Administrator",
				Link:        "https://www.credly.com/badges/f1c5619d-f6a1-4988-8393-5f9a21455736/linked_in_profile",
			},
			{
				Description: "Certified Kubernetes Security Specialist",
				Link:        "https://www.credly.com/badges/9e2a89df-4283-4502-9834-7b11b05bb152/linked_in_profile",
			},
		},
	}
	contact := &types.Contact{
		Name:      name + "-contact",
		Namespace: namespace.Name,
		OwnerRef:  ownerRef,
		Email:     "edudiazasencio@gmail.com",
		Linkedin:  "https://www.linkedin.com/in/eduardo-diaz-asencio/",
		Github:    "https://github.com/edu-diaz",
	}
	contributions := &types.Contributions{
		Name:      name + "-contributions",
		Namespace: namespace.Name,
		OwnerRef:  ownerRef,
		Contributions: []types.Contribution{
			{
				Project:     "Spin",
				Description: "Developer tool for building WebAssembly microservices and web applications from Fermyon",
				Link:        "https://github.com/fermyon/developer/graphs/contributors",
			},
			{
				Project:     "CloudTTY",
				Description: "A Kubernetes cloudshell operator in the Cloud Native Computing Foundation",
				Link:        "https://cloudtty.github.io/cloudtty/#contributors",
			},
			{
				Project:     "Jesse",
				Description: "A Kubernetes helm chart for the Jesse AI trading bot",
				Link:        "https://github.com/TrianaLab/jesse-chart",
			},
			{
				Project:     "awasm-portfolio",
				Description: "This portfolio, isn't it awesome?",
				Link:        "https://github.com/TrianaLab/awasm-portfolio",
			},
		},
	}
	education := &types.Education{
		Name:      name + "-education",
		Namespace: namespace.Name,
		OwnerRef:  ownerRef,
		Courses: []types.Course{
			{
				Title:       "Artificial Intelligence - M.Eng",
				Institution: "Universidad Internacional de Valencia",
				Duration:    "2021 - 2022",
			},
			{
				Title:       "Telecommunications Engineering - B.Eng",
				Institution: "Universidad de Sevilla",
				Duration:    "2016 - 2021",
			},
		},
	}
	experience := &types.Experience{
		Name:      name + "-experience",
		Namespace: namespace.Name,
		OwnerRef:  ownerRef,
		Jobs: []types.Job{
			{
				Title:       "Machine Learning Operations Engineer",
				Description: "As an MLOps Engineer at Emergence, I’m actively involved in building scalable and efficient AI infrastructure utilizing technologies like Terraform, Crossplane, Prometheus, and Istio. My responsibilities include provisioning GKE clusters, orchestrating additional infrastructure components, ensuring robust observability and implementing advanced networking capabilities.",
				Company:     "Emergence AI",
				Duration:    "July 2024 - Now",
			},
			{
				Title:       "Product Software Engineer - Kubernetes Team",
				Description: "Provide services that support elastic scale and allow frequent, reliable, high-impact changes to the deployed products. Reduce friction and toil surrounding data when creating new product services and features, including data lifecycle management, data retention, data analytics and providing easy-to-use APIs. Make Appian more Kubernetes-native both in cloud and self-managed environments.",
				Company:     "Appian Corporation",
				Duration:    "February 2024 - July 2024",
			},
			{
				Title:       "Senior Solution Engineer - Infrastructure Team",
				Description: "Provide technical support globally, address critical challenges and mentor newcomers. With a focus on data analysis and creative solutions, my role emphasizes effective troubleshooting and a comprehensive understanding of the platform’s inner infrastructure.",
				Company:     "Appian Corporation",
				Duration:    "October 2023 - February 2024",
			},
			{
				Title:       "Solution Engineer - Infrastructure Team",
				Description: "Hands-on support for global customers utilizing the Appian platform. With strong problem-solving skills, proficiency in Kubernetes and expertise in web services, programming, and Linux, I contribute to ongoing customer relationships by delivering effective solutions.",
				Company:     "Appian Corportation",
				Duration:    "October 2022 - October 2023",
			},
			{
				Title:       "Associate Solution Engineer",
				Description: "Tackle complex technical challenges, providing creative solutions and offering world-class support to customers globally. Proficient in troubleshooting, data analytics, and collaboration with internal teams, contribute to resolve Appian installations for both self-managed and Appian Cloud environments.",
				Company:     "Appian Corporation",
				Duration:    "November 2021 - October 2022",
			},
			{
				Title:       "Software QA Automation Engineer",
				Description: "Interface with developers and system architects to ensure applications are designed to be testable while ensuring tags, object ID’s, component and page name standards are in place. Create test plans and test cases based on defined stories. Automate those test cases and incorporate them to correspondent test suites.",
				Company:     "Solera Inc.",
				Duration:    "August 2020 - November 2021",
			},
		},
	}
	profile := &types.Profile{
		Name:           name + "-profile",
		Namespace:      namespace.Name,
		Contributions:  *contributions,
		Contact:        *contact,
		Certifications: *certifications,
		Education:      *education,
		Experience:     *experience,
	}

	resources := []models.Resource{
		namespace,
		certifications,
		contact,
		contributions,
		education,
		experience,
		profile,
	}

	// Iterate over the resources and create them
	for _, resource := range resources {
		repo.Create(resource)
	}
}
