package repos

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	models "test-manager/usecase_models/boiler"
)

type ProjectsRepository interface {
	UpdateProjects(ctx context.Context, Project models.Project, blackList ...string) error
	UpdateProjectsWithRules(ctx context.Context, Project models.Project, blackList ...string) error
	GetProjects(ctx context.Context, accountId int) ([]*models.Project, error)
	GetProject(ctx context.Context, id int) (models.Project, error)
	GetProjectWithLoads(ctx context.Context, id int) (models.Project, error)
	SaveProjects(ctx context.Context, Project models.Project) (int, error)
	GetProjectsInMembers(ctx context.Context, email string) ([]*models.Project, error)
}

type projectsRepository struct {
	db *sql.DB
}

func NewProjectsRepository(db *sql.DB) ProjectsRepository {
	return &projectsRepository{db: db}
}

func (r *projectsRepository) SaveProjects(ctx context.Context, Project models.Project) (int, error) {
	err := Project.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return 0, err
	}
	return Project.ID, nil
}

func (r *projectsRepository) UpdateProjects(ctx context.Context, Project models.Project, blackList ...string) error {
	_, err := Project.Update(ctx, r.db, boil.Blacklist(blackList...))
	if err != nil {
		return err
	}
	return nil
}

func (r *projectsRepository) UpdateProjectsWithRules(ctx context.Context, Project models.Project, blackList ...string) error {
	_, err := Project.Update(ctx, r.db, boil.Blacklist(blackList...))
	if err != nil {
		return err
	}
	return nil
}

func (r *projectsRepository) GetProjects(ctx context.Context, accountId int) ([]*models.Project, error) {
	Project, err := models.Projects(models.ProjectWhere.AccountID.EQ(accountId)).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return Project, nil
}

func (r *projectsRepository) GetProject(ctx context.Context, id int) (models.Project, error) {
	Project, err := models.Projects(models.ProjectWhere.ID.EQ(id)).One(ctx, r.db)
	if err != nil {
		return models.Project{}, err
	}
	return *Project, nil
}

func (r *projectsRepository) GetProjectsInMembers(ctx context.Context, email string) ([]*models.Project, error) {
	Project, err := models.Projects(qm.Where("members @> $1", fmt.Sprintf(`[{"email": "%s"}]`, email))).All(ctx, r.db)
	if err != nil {
		return []*models.Project{}, err
	}
	return Project, nil

	//rows, err := r.db.Query("SELECT * FROM projects WHERE members @> $1")
	//if err != nil {
	//	return nil, err
	//}
	//defer rows.Close()
	//
	//projects := make([]*models.Project, 0)
	//// Iterate through the results
	//for rows.Next() {
	//	var project *models.Project
	//	err := rows.Scan(&project)
	//	if err != nil {
	//		return nil, err
	//	}
	//	projects = append(projects, project)
	//}
	//return projects, nil
}

func (r *projectsRepository) GetProjectWithLoads(ctx context.Context, id int) (models.Project, error) {
	Project, err := models.Projects(models.ProjectWhere.ID.EQ(id), qm.Load(models.ProjectRels.Account)).One(ctx, r.db)
	if err != nil {
		return models.Project{}, err
	}
	return *Project, nil
}
