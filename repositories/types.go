package repositories

type RegistryType string

const (
	PublicRegistry  RegistryType = "public"
	PrivateRegistry RegistryType = "private"
)

type RepositoryType string

const (
	PublicRepository  RepositoryType = "public"
	PrivateRepository RepositoryType = "private"
)
