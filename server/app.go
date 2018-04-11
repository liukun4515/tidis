//
// app.go
// Copyright (C) 2018 YanMing <yming0221@gmail.com>
//
// Distributed under terms of the MIT license.
//

package server

import (
	"net"
	"sync"

	"github.com/YongMan/go/log"
	"github.com/YongMan/tedis/config"
	"github.com/YongMan/tedis/tedis"
)

type App struct {
	conf *config.Config

	listener net.Listener

	// wrapper and manager for db instance
	tdb *tedis.Tedis

	quitCh chan bool

	clientWG sync.WaitGroup

	//client map?
}

// initialize an app
func NewApp(conf *config.Config) *App {
	var err error
	app := &App{
		conf: conf,
	}
	app.tdb, err = tedis.NewTedis(conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	app.listener, err = net.Listen("tcp", "0.0.0.0:6666")
	log.Info("server listen in 0.0.0.0:6666")
	if err != nil {
		log.Fatal(err.Error())
	}

	return app
}

func (app *App) Close() error {
	return nil
}

func (app *App) Run() {

	for {
		select {
		case <-app.quitCh:
			return
		default:
			// accept new client connect and perform
			log.Debug("waiting for new connection")
			conn, err := app.listener.Accept()
			if err != nil {
				log.Error(err.Error())
				continue
			}
			// handle conn
			log.Debug("handle new connection")
			ClientHandler(conn, app)
		}
	}
}