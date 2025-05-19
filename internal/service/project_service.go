package service

type IProjectService interface {
	IsValidProject(projectId uint) (bool, error)
}

type ProjectService struct {
}

func NewProjectService() *ProjectService {
	return &ProjectService{}
}

func (p *ProjectService) IsValidProject(projectId uint) (bool, error) {

	// TODO: Access the DB adn validate the projectId
	// Returning based ont he default project for now.
	if projectId == 1 {
		return true, nil
	}
	return false, ErrInvalidProject
}
