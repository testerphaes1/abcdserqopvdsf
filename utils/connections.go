package utils

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/hibiken/asynq"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func PostgresConnection(host, port, user, pass, database, sslmode string, maxOpenConns, maxIdleConns int) (*sql.DB, error) {
	connString := postgresURI(host, port, user, pass, database, sslmode)
	log.Println("postgres options -> " + connString)
	conn, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("error in openning postgres connection: %w", err)
	}

	conn.SetMaxOpenConns(maxOpenConns)
	conn.SetMaxIdleConns(maxIdleConns)

	dbContext, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = conn.PingContext(dbContext)
	if err != nil {
		return nil, fmt.Errorf("error in pinging postgres database: %w", err)
	}
	return conn, nil
}

func postgresURI(host, port, user, pass, database, sslmode string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, database, sslmode)
}

func CreateRedisConnection(ctx context.Context, host, port string, db int, timeout time.Duration) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%s", host, port),
		DialTimeout: timeout,
		DB:          db,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	if pong != "PONG" {
		return nil, fmt.Errorf("expected PONG, got %s", pong)
	}

	return client, nil
}

func CreateSentinelConnection(ctx context.Context, host, port string, db int, timeout time.Duration) (*redis.Client, error) {
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName: "mymaster",
		SentinelAddrs: []string{
			fmt.Sprintf("%s:%s", host, port),
		},
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	if pong != "PONG" {
		panic(err)
	}

	return client, nil
}

func CreateInfluxDBConnection(ctx context.Context, token string, url string, org string, bucket string) (influxdb2.Client, api.WriteAPIBlocking, api.QueryAPI, error) {
	client := influxdb2.NewClient(url, token)
	p, err := client.Ping(ctx)
	if err != nil || !p {
		return nil, nil, nil, err
	}
	h, err := client.Health(ctx)
	if err != nil || h.Status != domain.HealthCheckStatusPass {
		return nil, nil, nil, err
	}

	writeAPI := client.WriteAPIBlocking(org, bucket)
	queryAPI := client.QueryAPI(org)

	return client, writeAPI, queryAPI, nil
}

func AsynqInspector(redisConf *redis.Options) *asynq.Inspector {
	return asynq.NewInspector(asynq.RedisClientOpt{
		Addr:        redisConf.Addr,
		DialTimeout: redisConf.DialTimeout,
		Username:    redisConf.Username,
		Password:    redisConf.Password,
	})
}
