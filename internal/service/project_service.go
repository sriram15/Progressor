package service

type ProjectService interface {
	IsValidProject(projectId uint) (bool, error)
}

type projectService struct {
}

func NewProjectService() ProjectService {
	return &projectService{}
}

func (p *projectService) IsValidProject(projectId uint) (bool, error) {

	// TODO: Access the DB adn validate the projectId
	// Returning based ont he default project for now.
	if projectId == 1 {
		return true, nil
	}
	return false, ErrInvalidProject
}
