package repos

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"test-manager/usecase_models"
	models "test-manager/usecase_models/boiler"
	"time"
)

type NetCatRepository interface {
	UpdateNetCat(ctx context.Context, NetCat models.NetCat) error
	GetNetCats(ctx context.Context, projectId int) (netCatUseCase []*usecase_models.NetCats, err error)
	GetNetCat(ctx context.Context, id int) (netCatUseCase *usecase_models.NetCats, err error)
	GetActiveNetCats(ctx context.Context) (netCatUseCase []*usecase_models.NetCats, err error)
	SaveNetCat(ctx context.Context, NetCat models.NetCat) (int, error)
	DeleteNetcat(ctx context.Context, netcatId int) error
}

type netCatRepository struct {
	db *sql.DB
}

func NewNetCatRepository(db *sql.DB) NetCatRepository {
	return &netCatRepository{db: db}
}

func (r *netCatRepository) SaveNetCat(ctx context.Context, netCat models.NetCat) (int, error) {
	err := netCat.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return 0, err
	}
	return netCat.ID, nil
}

func (r *netCatRepository) UpdateNetCat(ctx context.Context, netCat models.NetCat) error {
	_, err := netCat.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}

func (r *netCatRepository) GetNetCats(ctx context.Context, projectId int) (netCatUseCase []*usecase_models.NetCats, err error) {
	netCats, err := models.NetCats(models.NetCatWhere.ProjectID.EQ(projectId)).All(ctx, r.db)
	if err != nil {
		return []*usecase_models.NetCats{}, err
	}

	for _, value := range netCats {
		var netCat usecase_models.NetCats
		err := json.Unmarshal(value.Data.JSON, &netCat)
		if err != nil {
			log.Println(err.Error())
		}
		netCat.Scheduling.PipelineId = value.ID
		netCatUseCase = append(netCatUseCase, &netCat)
	}
	return netCatUseCase, nil
}

func (r *netCatRepository) GetNetCat(ctx context.Context, id int) (netCatUseCase *usecase_models.NetCats, err error) {
	var netcat models.NetCat
	err = models.NetCats(models.NetCatWhere.ID.EQ(id)).Bind(ctx, r.db, &netcat)
	if err != nil {
		return &usecase_models.NetCats{}, err
	}

	err = json.Unmarshal(netcat.Data.JSON, &netCatUseCase)
	if err != nil {
		log.Println(err.Error())
	}
	netCatUseCase.Scheduling.PipelineId = netcat.ID

	return netCatUseCase, nil
}

func (r *netCatRepository) GetActiveNetCats(ctx context.Context) (netCatUseCase []*usecase_models.NetCats, err error) {
	netCats, err := models.NetCats(qm.Where("data->'scheduling'->>'is_active' = ? and data->'scheduling'->>'end_at' > ?", "true", time.Now().Format("2006-01-02 15:04:05"))).All(ctx, r.db)
	if err != nil {
		return []*usecase_models.NetCats{}, err
	}

	for _, value := range netCats {
		var netCat usecase_models.NetCats
		err := json.Unmarshal(value.Data.JSON, &netCat)
		if err != nil {
			log.Println(err.Error())
		}
		netCat.Scheduling.PipelineId = value.ID
		netCatUseCase = append(netCatUseCase, &netCat)
	}
	return netCatUseCase, nil
}

func (r *netCatRepository) DeleteNetcat(ctx context.Context, netcatId int) error {
	netcat := models.NetCat{ID: netcatId}
	_, err := netcat.Delete(ctx, r.db, false)
	if err != nil {
		return err
	}
	return nil
}
