package repos

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"reflect"
	"strings"
	models "test-manager/usecase_models/boiler"
	"time"
)

type EndpointStatsRepository interface {
	Write(ctx context.Context, options WriteEndpointStatsOptions) error
	WriteBulk(ctx context.Context, options []WriteEndpointStatsOptions) error
	Read(ctx context.Context, selectors []string, filters Filters, loads []string, limit int, offset int, total bool) (int64, models.EndpointStatSlice, error)
	GetLastNSessionsByEndpointId(ctx context.Context, n int, endpointId int) ([]string, error)
	GetSessionSuccessions(ctx context.Context, filters Filters) ([]EndpointSessionSuccessions, error)
	CleanUpResponseBodies(ctx context.Context, time time.Time) error
	VacuumEndpointStats(ctx context.Context) error
}

type endpointStatsRepository struct {
	db *sql.DB
}

func NewEndpointStatsRepository(db *sql.DB) EndpointStatsRepository {
	return &endpointStatsRepository{db: db}
}

type WriteEndpointStatsOptions struct {
	Time             time.Time `json:"time"`
	ProjectId        int       `json:"project_id"`
	SessionId        string    `json:"session_id"`
	EndpointName     string    `json:"endpoint_name"`
	EndpointId       int       `json:"endpoint_id"`
	IsHeartBeat      bool      `json:"is_heart_beat"`
	Url              string    `json:"url"`
	DatacenterId     int       `json:"datacenter_id"`
	Success          int       `json:"success"`
	ResponseTime     float64   `json:"response_time"`
	ResponseTimes    string    `json:"response_times"`
	ResponseBodies   string    `json:"response_bodies"`
	ResponseHeaders  string    `json:"response_headers"`
	ResponseStatuses string    `json:"response_statuses"`
}

func (e *endpointStatsRepository) Write(ctx context.Context, options WriteEndpointStatsOptions) error {
	endpointStat := models.EndpointStat{
		Time:             options.Time,
		SessionID:        options.SessionId,
		ProjectID:        options.ProjectId,
		EndpointName:     null.NewString(options.EndpointName, true),
		EndpointID:       options.EndpointId,
		URL:              null.NewString(options.Url, true),
		DatacenterID:     options.DatacenterId,
		IsHeartBeat:      options.IsHeartBeat,
		Success:          options.Success,
		ResponseTime:     options.ResponseTime,
		ResponseTimes:    null.NewString(options.ResponseTimes, true),
		ResponseBodies:   null.NewBytes([]byte(options.ResponseBodies), true),
		ResponseHeaders:  null.NewString(options.ResponseHeaders, true),
		ResponseStatuses: null.NewString(options.ResponseStatuses, true),
	}
	return endpointStat.Insert(ctx, e.db, boil.Infer())
}

func (e *endpointStatsRepository) WriteBulk(ctx context.Context, options []WriteEndpointStatsOptions) error {
	txn, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer txn.Rollback()

	stmt, err := txn.Prepare(pq.CopyIn("endpoint_stats",
		models.EndpointStatColumns.Time,
		models.EndpointStatColumns.SessionID,
		models.EndpointStatColumns.ProjectID,
		models.EndpointStatColumns.EndpointName,
		models.EndpointStatColumns.EndpointID,
		models.EndpointStatColumns.URL,
		models.EndpointStatColumns.DatacenterID,
		models.EndpointStatColumns.IsHeartBeat,
		models.EndpointStatColumns.Success,
		models.EndpointStatColumns.ResponseTime,
		models.EndpointStatColumns.ResponseTimes,
		models.EndpointStatColumns.ResponseBodies,
		models.EndpointStatColumns.ResponseHeaders,
		models.EndpointStatColumns.ResponseStatuses))
	if err != nil {
		return err
	}

	for _, option := range options {
		_, err = stmt.Exec(
			option.Time,
			option.SessionId,
			option.ProjectId,
			null.NewString(option.EndpointName, true),
			option.EndpointId,
			null.NewString(option.Url, true),
			option.DatacenterId,
			option.IsHeartBeat,
			option.Success,
			option.ResponseTime,
			null.NewString(option.ResponseTimes, true),
			null.NewBytes([]byte(option.ResponseBodies), true),
			null.NewString(option.ResponseHeaders, true),
			null.NewString(option.ResponseStatuses, true),
		)
		if err != nil {
			log.Error(err)
			continue
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (e *endpointStatsRepository) Read(ctx context.Context, selectors []string, filters Filters, loads []string, limit int, offset int, total bool) (int64, models.EndpointStatSlice, error) {
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

	tc := int64(0)
	var err error
	if total {

		tc, err = models.EndpointStats(append(qmQuery, qm.Select("time"))...).Count(ctx, e.db)
		if err != nil {
			return 0, models.EndpointStatSlice{}, err
		}
		qmQuery = append(qmQuery, qm.Limit(limit))
		qmQuery = append(qmQuery, qm.Offset(offset))
	}

	qmQuery = append(qmQuery, qm.OrderBy("time DESC"))
	for _, load := range loads {
		qmQuery = append(qmQuery, qm.Load(load))
	}
	selectQ := strings.Join(selectors, ",")
	if selectQ == "" {
		selectQ = "*"
	}
	qmQuery = append(qmQuery, qm.Select(selectQ))

	res, err := models.EndpointStats(qmQuery...).All(ctx, e.db)
	return tc, res, err
}

func (e *endpointStatsRepository) GetLastNSessionsByEndpointId(ctx context.Context, n int, endpointId int) (res []string, err error) {
	rows, err := queries.Raw("select session_id from (select session_id, max(time) as tt from endpoint_stats where endpoint_id = $1 group by session_id order by tt desc limit $2) as b", endpointId, n).QueryContext(ctx, e.db)
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

type EndpointSessionSuccessions struct {
	SessionId  string    `json:"session_id"`
	EndpointId int       `json:"endpoint_id"`
	MinTime    time.Time `json:"min_time"`
	MaxTime    time.Time `json:"max_time"`
	Success    bool      `json:"success"`
}

func (e *endpointStatsRepository) GetSessionSuccessions(ctx context.Context, filters Filters) ([]EndpointSessionSuccessions, error) {
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
		"endpoint_id",
		"min(time) as min_time",
		"max(time) as max_time",
		"bit_and(success) as success"))
	qmQuery = append(qmQuery, qm.OrderBy("max_time desc"))
	qmQuery = append(qmQuery, qm.GroupBy("session_id, endpoint_id"))

	var response []EndpointSessionSuccessions
	err := models.EndpointStats(qmQuery...).Bind(ctx, e.db, &response)
	return response, err
}

func (e *endpointStatsRepository) CleanUpResponseBodies(ctx context.Context, time time.Time) error {
	fmt.Println("clean up started")
	query := queries.Raw("update endpoint_stats set response_bodies = null where time < $1 and response_bodies is not null;", time)
	_, err := query.ExecContext(ctx, e.db)
	if err != nil {
		return err
	}
	fmt.Println("clean up successfully finished")
	return nil
}

func (e *endpointStatsRepository) VacuumEndpointStats(ctx context.Context) error {
	fmt.Println("vacuum started")
	query := queries.Raw("VACUUM verbose ANALYZE endpoint_stats;")
	_, err := query.ExecContext(ctx, e.db)
	if err != nil {
		return err
	}
	fmt.Println("vacuum successfully finished")
	return nil
}
