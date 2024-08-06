package repositories

type CreateRepositoryRequest struct {
	Name        string         `json:"name" validate:"required,min=1,max=30"`
	Description string         `json:"description" validate:"max=1000"`
	Type        RepositoryType `json:"type" validate:"required,oneof=public private"`
	Public      bool           `json:"public"`
}
