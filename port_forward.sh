#!/bin/bash

oc port-forward service/mysqlcentral 3307:3306 &
MYSQL_PID=$!

oc port-forward service/nats 4223:4222 &
NATS_PID=$!

oc port-forward service/postgrescentral 5433:5432 &
POSTGRES_PID=$!

oc port-forward service/rediscentral 6377:6379 &
REDIS_PID=$!

oc port-forward service/sanjeh 9001:8443 &
SANJEH_PID=$!

oc port-forward service/rasoul 9002:80 &
RASOUL_PID=$!

oc port-forward service/dabir 9003:80 &
DABIR_PID=$!

oc port-forward service/user_banning 9004:80 &
USER_BANNING_PID=$!

terminate() {
  echo -e "\nTerminating connections..."

  kill -SIGTERM $MYSQL_PID
  kill -SIGTERM $NATS_PID
  kill -SIGTERM $POSTGRES_PID
  kill -SIGTERM $REDIS_PID
  kill -SIGTERM $SANJEH_PID
  kill -SIGTERM $RASOUL_PID
  kill -SIGTERM $DABIR_PID
  kill -SIGTERM $USER_BANNING_PID

  echo -e "Finished terminating connections."
}

trap terminate INT TERM
wait
