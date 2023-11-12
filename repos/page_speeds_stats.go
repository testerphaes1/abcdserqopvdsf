package repos

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"reflect"
	"strings"
	models "test-manager/usecase_models/boiler"
	"time"
)

type PageSpeedStatsRepository interface {
	Write(ctx context.Context, time time.Time, options WritePageSpeedStatsOptions) error
	Read(ctx context.Context, selectors []string, filters Filters, loads []string) (models.PageSpeedsStatSlice, error)
	GetLastNSessionsByPageSpeedId(ctx context.Context, n int, pageSpeedId int) (res []string, err error)
	GetSessionSuccessions(ctx context.Context, filters Filters) ([]PageSpeedSessionSuccessions, error)
}

type pageSpeedStatsRepository struct {
	db *sql.DB
}

func NewPageSpeedStatsRepository(db *sql.DB) PageSpeedStatsRepository {
	return &pageSpeedStatsRepository{db: db}
}

type WritePageSpeedStatsOptions struct {
	ProjectId    int
	PageSpeedId  int
	IsHeartBeat  bool
	Url          string
	DatacenterId int
	Success      int
}

func (e *pageSpeedStatsRepository) Write(ctx context.Context, time time.Time, options WritePageSpeedStatsOptions) error {
	PageSpeedStat := models.PageSpeedsStat{
		Time:         time,
		ProjectID:    options.ProjectId,
		PagespeedID:  options.PageSpeedId,
		URL:          null.NewString(options.Url, true),
		DatacenterID: options.DatacenterId,
		IsHeartBeat:  options.IsHeartBeat,
		Success:      options.Success,
	}
	return PageSpeedStat.Insert(ctx, e.db, boil.Infer())
}

func (e *pageSpeedStatsRepository) Read(ctx context.Context, selectors []string, filters Filters, loads []string) (models.PageSpeedsStatSlice, error) {
	var qmQuery []qm.QueryMod
	for _, load := range loads {
		qmQuery = append(qmQuery, qm.Load(load))
	}
	selectQ := strings.Join(selectors, ",")
	if selectQ == "" {
		selectQ = "*"
	}
	qmQuery = append(qmQuery, qm.Select(selectQ))

	qmQuery = append(qmQuery, qm.OrderBy("time DESC"))

	for _, filter := range filters {
		qmQuery = append(qmQuery, qm.Where(fmt.Sprintf("%s %s ?", filter.Field, filter.Op), filter.Value))
	}

	return models.PageSpeedsStats(qmQuery...).All(ctx, e.db)
}

func (e *pageSpeedStatsRepository) GetLastNSessionsByPageSpeedId(ctx context.Context, n int, pageSpeedId int) (res []string, err error) {
	rows, err := queries.Raw("select session_id from (select session_id, max(time) as tt from page_speeds_stats where pagespeed_id = $1 group by session_id order by tt desc limit $2) as b", pageSpeedId, n).QueryContext(ctx, e.db)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		var temp string
		err = rows.Scan(&temp)
		if err != nil {
			return res, err
		}
		res = append(res, temp)
	}
	return res, err
}

type PageSpeedSessionSuccessions struct {
	SessionId   string    `json:"session_id"`
	PageSpeedId int       `json:"pagespeed_id"`
	MinTime     time.Time `json:"min_time"`
	MaxTime     time.Time `json:"max_time"`
	Success     bool      `json:"success"`
}

func (e *pageSpeedStatsRepository) GetSessionSuccessions(ctx context.Context, filters Filters) ([]PageSpeedSessionSuccessions, error) {
	var qmQuery []qm.QueryMod

	for _, filter := range filters {
		if filter.Op == FilterOpIn {
			var out []interface{}
			rv := reflect.ValueOf(filter.Value)
			if rv.Kind() == reflect.Slice {
				for i := 0; i < rv.Len(); i++ {
					out = append(out, rv.Index(i).Interface())
				}
			} else {
				continue
			}
			qmQuery = append(qmQuery, qm.WhereIn(fmt.Sprintf("%s %s ?", filter.Field, filter.Op), out...))
		} else {
			qmQuery = append(qmQuery, qm.Where(fmt.Sprintf("%s %s ?", filter.Field, filter.Op), filter.Value))
		}
	}

	qmQuery = append(qmQuery, qm.Select(
		"session_id",
		"pagespeed_id",
		"min(time) as min_time",
		"max(time) as max_time",
		"bit_and(success) as success"))
	qmQuery = append(qmQuery, qm.OrderBy("max_time desc"))
	qmQuery = append(qmQuery, qm.GroupBy("session_id, pagespeed_id"))

	var response []PageSpeedSessionSuccessions
	err := models.PageSpeedsStats(qmQuery...).Bind(ctx, e.db, &response)
	return response, err
}
