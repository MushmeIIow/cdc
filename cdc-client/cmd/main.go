package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx"
	"github.com/mushmellow/cdc-client/protobuf"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func main() {
	config, err := pgx.ParseConnectionString(os.Getenv("CONN_STRING"))
	if err != nil {
		logrus.Fatalln("failed to parse config", err)
	}

	logrus.Infoln("Connecting to database", config.Database)
	pgConn, err := pgx.Connect(config)
	if err != nil {
		logrus.Fatalln("failed to connect", err)
	}

	defer pgConn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go insertProto(ctx, pgConn)

	<-sigs
}

func insertProto(ctx context.Context, pgConn *pgx.Conn) {
	data := &protobuf.HelloRequest{Name: "foobar field"}
	a, err := proto.Marshal(data)
	if err != nil {
		logrus.Fatalln("failed to marshal", err)
	}

	insertFunc := func() {
		logrus.Infoln("Inserting...")
		_, err := pgConn.Exec(fmt.Sprintf("INSERT INTO protobuf values(1, '%s');", a))
		if err != nil {
			logrus.Fatalln("failed to execute", err)
		}
		logrus.Infoln("Inserted ", a)
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			insertFunc()
		}
	}
}
