package repos

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"strconv"
	"test-manager/cache"
	models "test-manager/usecase_models/boiler"
)

const (
	DatacenterPrefixCacheKey = "datacenter_id:"
	DatacenterIdsCacheKey    = "all_datacenter_ids"
)

type DataCentersRepository interface {
	UpdateDataCenters(ctx context.Context, dataCenter models.Datacenter) error
	GetDataCenter(ctx context.Context, id int) (models.Datacenter, error)
	GetDataCenterWithCache(ctx context.Context, id int) (*models.Datacenter, error)
	GetDataCenters(ctx context.Context) ([]*models.Datacenter, error)
	GetDataCentersWithCache(ctx context.Context) ([]*models.Datacenter, error)
	GetDataCenterByTitle(ctx context.Context, title string) (models.Datacenter, error)
	SaveDataCenters(ctx context.Context, dataCenter models.Datacenter) (int, error)
}

type dataCentersRepository struct {
	cacheRepo cache.Cache
	db        *sql.DB
}

func NewDataCentersRepositoryRepository(
	cacheRepo cache.Cache,
	db *sql.DB) DataCentersRepository {
	return &dataCentersRepository{
		cacheRepo: cacheRepo,
		db:        db,
	}
}

func (r *dataCentersRepository) SaveDataCenters(ctx context.Context, dataCenter models.Datacenter) (int, error) {
	err := dataCenter.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return 0, err
	}

	err = r.cacheRepo.Delete(ctx, DatacenterPrefixCacheKey+strconv.Itoa(dataCenter.ID))
	if err != nil {
		log.Error(err)
	}

	err = r.cacheRepo.Delete(ctx, DatacenterIdsCacheKey)
	if err != nil {
		log.Error(err)
	}
	return dataCenter.ID, nil
}

func (r *dataCentersRepository) UpdateDataCenters(ctx context.Context, dataCenter models.Datacenter) error {
	_, err := dataCenter.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return err
	}

	err = r.cacheRepo.Delete(ctx, DatacenterPrefixCacheKey+strconv.Itoa(dataCenter.ID))
	if err != nil {
		log.Error(err)
	}

	err = r.cacheRepo.Delete(ctx, DatacenterIdsCacheKey)
	if err != nil {
		log.Error(err)
	}
	return nil
}

func (r *dataCentersRepository) GetDataCenter(ctx context.Context, id int) (models.Datacenter, error) {
	datacenter, err := models.Datacenters(models.DatacenterWhere.ID.EQ(id)).One(ctx, r.db)
	if err != nil {
		return models.Datacenter{}, err
	}
	return *datacenter, nil
}

func (r *dataCentersRepository) GetDataCenterWithCache(ctx context.Context, id int) (datacenter *models.Datacenter, err error) {
	datacenterInt, err := r.cacheRepo.Get(ctx, DatacenterPrefixCacheKey+strconv.Itoa(id))
	if err != nil {
		fmt.Println("problem on fetching datacenter from cache: ", err)
		datacenter, err = models.Datacenters(models.DatacenterWhere.ID.EQ(id)).One(ctx, r.db)
		if err != nil {
			return &models.Datacenter{}, err
		}
		jd, err := json.Marshal(datacenter)
		if err == nil {
			err = r.cacheRepo.Set(ctx, DatacenterPrefixCacheKey+strconv.Itoa(datacenter.ID), jd, 0)
			if err != nil {
				fmt.Println("set datacenter cache key error: %v", err)
			}
		}
		return datacenter, nil
	}

	datacenter, ok := datacenterInt.(*models.Datacenter)
	if !ok {
		datacenter, err = models.Datacenters(models.DatacenterWhere.ID.EQ(id)).One(ctx, r.db)
		if err != nil {
			return &models.Datacenter{}, err
		}
		return datacenter, nil
	}

	return datacenter, nil
}

func (r *dataCentersRepository) GetDataCentersWithCache(ctx context.Context) ([]*models.Datacenter, error) {
	cacheRes, err := r.cacheRepo.HGet(ctx, DatacenterIdsCacheKey, "data")
	if err != nil || cacheRes == "" {
		log.Info("cache not found", err)

		datacenters, err := models.Datacenters().All(ctx, r.db)
		if err != nil {
			return []*models.Datacenter{}, err
		}
		jsonData, err := json.Marshal(datacenters)
		// Set the JSON data in Redis using HMSET
		if err == nil {
			err = r.cacheRepo.HSet(ctx, DatacenterIdsCacheKey, "data", jsonData)
			if err != nil {
				fmt.Println("Failed to perform HMSET:", err)
			}
		}

		return datacenters, nil
	}

	var res []*models.Datacenter
	err = json.Unmarshal([]byte(cacheRes), &res)
	if err != nil {
		fmt.Println("Failed to unmarshal data:", err)
		return nil, err
	}
	return res, nil
}

func (r *dataCentersRepository) GetDataCenters(ctx context.Context) ([]*models.Datacenter, error) {
	datacenters, err := models.Datacenters().All(ctx, r.db)
	if err != nil {
		return []*models.Datacenter{}, err
	}
	return datacenters, nil
}

func (r *dataCentersRepository) GetDataCenterByTitle(ctx context.Context, title string) (models.Datacenter, error) {
	datacenter, err := models.Datacenters(models.DatacenterWhere.Title.EQ(title)).One(ctx, r.db)
	if err != nil {
		return models.Datacenter{}, err
	}
	return *datacenter, nil
}
