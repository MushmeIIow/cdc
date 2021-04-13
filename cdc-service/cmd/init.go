package main

import (
	"fmt"
	"time"

	"github.com/jackc/pgx"
	"github.com/mushmellow/cdc-service/config"
	"github.com/mushmellow/cdc-service/listener"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// logger log levels.
const (
	warningLoggerLevel = "warning"
	errorLoggerLevel   = "error"
	fatalLoggerLevel   = "fatal"
	infoLoggerLevel    = "info"
)

// initLogger init logrus preferences.
func initLogger(cfg config.LoggerCfg) {
	logrus.SetReportCaller(cfg.Caller)
	if !cfg.HumanReadable {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	var level logrus.Level
	switch cfg.Level {
	case warningLoggerLevel:
		level = logrus.WarnLevel
	case errorLoggerLevel:
		level = logrus.ErrorLevel
	case fatalLoggerLevel:
		level = logrus.FatalLevel
	case infoLoggerLevel:
		level = logrus.InfoLevel
	default:
		level = logrus.DebugLevel
	}
	logrus.SetLevel(level)
}

// initPgxConnections initialise db and replication connections.
func initPgxConnections(cfg config.DatabaseCfg) (pgConn *pgx.Conn, rConnection *pgx.ReplicationConn, err error) {
	pgxConf := pgx.ConnConfig{
		// TODO logger
		LogLevel: pgx.LogLevelInfo,
		Logger:   pgxLogger{},
		Host:     cfg.Host,
		Port:     cfg.Port,
		Database: cfg.Name,
		User:     cfg.User,
		Password: cfg.Password,
	}

	for i := 1; i < 5; i++ {
		logrus.Infoln("Trying to conect to DB. Attempt:", i)
		pgConn, err = pgx.Connect(pgxConf)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(i) * time.Second)
	}

	if err != nil {
		return nil, nil, errors.Wrap(err, listener.ErrPostgresConnection)
	}

	rConnection, err = pgx.ReplicationConnect(pgxConf)
	if err != nil {
		return nil, nil, fmt.Errorf("%v: %w", listener.ErrReplicationConnection, err)
	}

	return
}

type pgxLogger struct{}

func (l pgxLogger) Log(level pgx.LogLevel, msg string, data map[string]interface{}) {
	logrus.Debugln(msg)
}
