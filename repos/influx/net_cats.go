package influx

import (
	"context"
	"database/sql"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"strconv"
	"time"
)

type NetCatsReportRepository interface {
	WriteNetCatsReport(ctx context.Context, projectId int, url string, success int) error
	ReadNetCatsReportByProject(ctx context.Context, projectId int, url string, timeFrame string, fields []string) (err error, res []interface{})
}

type netCatsReportRepository struct {
	bucket   string
	writeAPI api.WriteAPIBlocking
	queryAPI api.QueryAPI

	db *sql.DB
}

func NewNetCatsReportRepository(bucket string, writeAPI api.WriteAPIBlocking, queryAPI api.QueryAPI, db *sql.DB) NetCatsReportRepository {
	return &netCatsReportRepository{
		bucket:   bucket,
		writeAPI: writeAPI,
		queryAPI: queryAPI,
		db:       db,
	}
}

func (r *netCatsReportRepository) WriteNetCatsReport(ctx context.Context, projectId int, url string, success int) error {
	p := influxdb2.NewPoint("net_cats",
		map[string]string{"project_id": strconv.Itoa(projectId), "url": url},
		map[string]interface{}{"success": success},
		time.Now())
	err := r.writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		return err
	}
	return nil
}

func (r *netCatsReportRepository) ReadNetCatsReportByProject(ctx context.Context, projectId int, url string, timeFrame string, fields []string) (err error, res []interface{}) {
	fieldsQuery := ""
	for _, value := range fields {
		fieldsQuery = fieldsQuery + fmt.Sprintf("or r[\"_field\"] == \"%s\"", value)
	}
	query := ""
	if fieldsQuery != "" {
		fieldsQuery = fieldsQuery[2:]
		if url != "" {
			query = fmt.Sprintf(`from(bucket:"%s")
		 |> range(start: -%s) 
		 |> filter(fn: (r) => r["_measurement"] == "net_cats")
		 |> filter(fn: (r) => %s)
		 |> filter(fn: (r) => r.project_id == "%s")
		 |> filter(fn: (r) => r.url == "%s")
		 |> aggregateWindow(every: 10s, fn: last, createEmpty: false)
		 |> yield(name: "last")`, r.bucket, timeFrame, fieldsQuery, strconv.Itoa(projectId), url)
		} else {
			query = fmt.Sprintf(`from(bucket:"%s")
		 |> range(start: -%s) 
		 |> filter(fn: (r) => r["_measurement"] == "net_cats")
		 |> filter(fn: (r) => %s)
		 |> filter(fn: (r) => r.project_id == "%s")
		 |> aggregateWindow(every: 10s, fn: last, createEmpty: false)
		 |> yield(name: "last")`, r.bucket, timeFrame, fieldsQuery, strconv.Itoa(projectId))
		}
	} else {
		if url != "" {
			query = fmt.Sprintf(`from(bucket:"%s")
		 |> range(start: -%s) 
		 |> filter(fn: (r) => r["_measurement"] == "net_cats")
		 |> filter(fn: (r) => r.project_id == "%s")
		 |> filter(fn: (r) => r.url == "%s")
		 |> aggregateWindow(every: 10s, fn: last, createEmpty: false)
		 |> yield(name: "last")`, r.bucket, timeFrame, strconv.Itoa(projectId), url)
		} else {
			query = fmt.Sprintf(`from(bucket:"%s")
		 |> range(start: -%s) 
		 |> filter(fn: (r) => r["_measurement"] == "net_cats")
		 |> filter(fn: (r) => r.project_id == "%s")
		 |> aggregateWindow(every: 10s, fn: last, createEmpty: false)
		 |> yield(name: "last")`, r.bucket, timeFrame, strconv.Itoa(projectId))
		}
	}

	fmt.Println(query)
	result, err := r.queryAPI.Query(context.Background(), query)

	if err == nil {
		for result.Next() {
			fmt.Printf("value: %v\n", result.Record().Values())
			res = append(res, result.Record().Values())
		}
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
			return err, res
		}
	} else {
		return err, res
	}
	return nil, res
}
