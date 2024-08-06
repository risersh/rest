package features

type GroupType string
type FeatureType string

type Group struct {
	ID          GroupType `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Features    []Feature `json:"features"`
}

type Feature struct {
	Group       GroupType   `json:"group"`
	ID          FeatureType `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
}

var (
	Repository    GroupType   = "repository"
	DockerBuilder FeatureType = "repository.dockerbuilder"
)

var Groups = []Group{
	{
		ID:          Repository,
		Name:        "Repository",
		Description: "Repository features",
		Features: []Feature{
			{
				ID:          DockerBuilder,
				Name:        "Docker Builder",
				Description: "Builds Docker images from the repository.",
			},
		},
	},
}

func GetGroup(name string) *Group {
	for _, group := range Groups {
		if group.ID == GroupType(name) {
			return &group
		}
	}
	return nil
}

func GetFeature(name string) *Feature {
	for _, group := range Groups {
		for _, feature := range group.Features {
			if feature.ID == FeatureType(name) {
				feature.Group = group.ID
				return &feature
			}
		}
	}

	return nil
}
