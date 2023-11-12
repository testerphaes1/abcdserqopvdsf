package repos

import (
	"context"
	"database/sql"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	models "test-manager/usecase_models/boiler"
)

type PackagesRepository interface {
	UpdatePackages(ctx context.Context, Package models.Package) error
	GetPackages(ctx context.Context) ([]*models.Package, error)
	GetPackage(ctx context.Context, id int) (models.Package, error)
	SavePackages(ctx context.Context, Package models.Package) (int, error)
}

type packagesRepository struct {
	db *sql.DB
}

func NewPackagesRepository(db *sql.DB) PackagesRepository {
	return &packagesRepository{db: db}
}

func (r *packagesRepository) SavePackages(ctx context.Context, Package models.Package) (int, error) {
	err := Package.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return 0, err
	}
	return Package.ID, nil
}

func (r *packagesRepository) UpdatePackages(ctx context.Context, Package models.Package) error {
	_, err := Package.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}

func (r *packagesRepository) GetPackages(ctx context.Context) ([]*models.Package, error) {
	Package, err := models.Packages(qm.OrderBy("id")).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return Package, nil
}

func (r *packagesRepository) GetPackage(ctx context.Context, id int) (models.Package, error) {
	Package, err := models.Packages(models.PackageWhere.ID.EQ(id)).One(ctx, r.db)
	if err != nil {
		return models.Package{}, err
	}
	return *Package, nil
}
