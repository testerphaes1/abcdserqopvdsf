package influx

import (
	"context"
	"database/sql"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"strconv"
	"strings"
	"time"
)

type EndpointReportRepository interface {
	WriteEndpointReport(ctx context.Context, options WriteEndpointReportOptions) error
	EndpointReport(ctx context.Context, timeFrame string, enableAggregate bool, aggregateYield string, aggregateTimeframe string, tags map[string]string, fields []string, enablePivot bool, limit int, offset int) (err error, res []interface{})
	ReadEndpointReportByProject(ctx context.Context, projectId int, pipelineId int, timeFrame string, fields []string) (err error, res []interface{})
}

type endpointReportRepository struct {
	bucket   string
	writeAPI api.WriteAPIBlocking
	queryAPI api.QueryAPI

	db *sql.DB
}

func NewEndpointReportRepository(bucket string, writeAPI api.WriteAPIBlocking, queryAPI api.QueryAPI, db *sql.DB) EndpointReportRepository {
	return &endpointReportRepository{
		bucket:   bucket,
		writeAPI: writeAPI,
		queryAPI: queryAPI,
		db:       db,
	}
}

type WriteEndpointReportOptions struct {
	ProjectId       int
	EndpointName    string
	PipeLineId      int
	Url             string
	DatacenterTitle string
	Success         int
	ResponseTime    float64
	ResponseBody    string
	ResponseHeader  string
	ResponseStatus  int
}

func (r *endpointReportRepository) WriteEndpointReport(ctx context.Context, options WriteEndpointReportOptions) error {
	p := influxdb2.NewPoint("endpoint",
		map[string]string{
			"project_id":       strconv.Itoa(options.ProjectId),
			"endpoint_name":    options.EndpointName,
			"pipeline_id":      strconv.Itoa(options.PipeLineId),
			"url":              options.Url,
			"datacenter_title": options.DatacenterTitle,
		},
		map[string]interface{}{
			"success":         options.Success,
			"response_time":   options.ResponseTime,
			"response_body":   options.ResponseBody,
			"response_header": options.ResponseHeader,
			"response_status": options.ResponseStatus,
		},
		time.Now())
	err := r.writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		return err
	}
	return nil
}

func (r *endpointReportRepository) EndpointReport(ctx context.Context, startTime string,
	enableAggregate bool, aggregateYield string, aggregateTimeframe string, tags map[string]string, fields []string, enablePivot bool, limit int, offset int) (err error, res []interface{}) {
	var tagsQ []string
	for key, value := range tags {
		tagsQ = append(tagsQ, fmt.Sprintf("r[\"%s\"] == \"%s\"", key, value))
	}
	var fieldsQ []string
	for _, value := range fields {
		fieldsQ = append(fieldsQ, fmt.Sprintf("r._field == \"%s\"", value))
	}

	tagsQuery := ""
	if len(tagsQ) != 0 {
		tagsQuery = fmt.Sprintf("|> filter(fn: (r) => %s)", strings.Join(tagsQ, " and "))
	}
	fieldsQuery := ""
	if len(fieldsQ) != 0 {
		fieldsQuery = fmt.Sprintf("|> filter(fn: (r) => %s)", strings.Join(fieldsQ, " or "))
	}

	aggQuery := ""
	if enableAggregate {
		aggQuery = fmt.Sprintf(`|> aggregateWindow(every: %s, fn: %s, createEmpty: true)
                                       |> yield(name: "%s")`, aggregateTimeframe, aggregateYield, aggregateYield)
	}

	limitQuery := ""
	if limit != 0 {
		limitQuery = fmt.Sprintf(`|> limit(n: %d, offset: %d)`, limit, offset)
	}

	pivotQuery := ""
	if enablePivot {
		pivotQuery = `|> pivot(rowKey: ["_time"], columnKey: ["_field"], valueColumn: "_value")`
	}

	query := fmt.Sprintf(`from(bucket:"%s")
		 |> range(start: -%s)
		 |> filter(fn: (r) => r["_measurement"] == "endpoint")
		 %s
		 %s
		 %s
  	     %s
		 %s`, r.bucket, startTime, tagsQuery, fieldsQuery, aggQuery, pivotQuery, limitQuery)

	result, err := r.queryAPI.Query(context.Background(), query)

	if err == nil {
		for result.Next() {
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

func (r *endpointReportRepository) ReadEndpointReportByProject(ctx context.Context, projectId int, pipelineId int, timeFrame string, fields []string) (err error, res []interface{}) {
	fieldsQuery := ""
	for _, value := range fields {
		fieldsQuery = fieldsQuery + fmt.Sprintf("or r[\"_field\"] == \"%s\"", value)
	}
	query := ""
	if fieldsQuery != "" {
		fieldsQuery = fieldsQuery[2:]
		if pipelineId != 0 {
			query = fmt.Sprintf(`from(bucket:"%s")
		 |> range(start: -%s) 
		 |> filter(fn: (r) => r["_measurement"] == "endpoint")
		 |> filter(fn: (r) => %s)
		 |> filter(fn: (r) => r.project_id == "%s")
		 |> filter(fn: (r) => r.pipeline_id == "%s")`, r.bucket, timeFrame, fieldsQuery, strconv.Itoa(projectId), strconv.Itoa(pipelineId))
		} else {
			query = fmt.Sprintf(`from(bucket:"%s")
		 |> range(start: -%s) 
		 |> filter(fn: (r) => r["_measurement"] == "endpoint")
		 |> filter(fn: (r) => %s)
		 |> filter(fn: (r) => r.project_id == "%s")`, r.bucket, timeFrame, fieldsQuery, strconv.Itoa(projectId))
		}
	} else {
		if pipelineId != 0 {
			query = fmt.Sprintf(`from(bucket:"%s")
		 |> range(start: -%s) 
		 |> filter(fn: (r) => r["_measurement"] == "endpoint")
		 |> filter(fn: (r) => r.project_id == "%s")
		 |> filter(fn: (r) => r.pipeline_id == "%s")`, r.bucket, timeFrame, strconv.Itoa(projectId), strconv.Itoa(pipelineId))
		} else {
			query = fmt.Sprintf(`from(bucket:"%s")
		 |> range(start: -%s) 
		 |> filter(fn: (r) => r["_measurement"] == "endpoint")
		 |> filter(fn: (r) => r.project_id == "%s")`, r.bucket, timeFrame, strconv.Itoa(projectId))
		}
	}

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
