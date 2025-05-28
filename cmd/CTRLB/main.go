package main

import (
	"context"
	"github/Bharatjawa2/CtrlB_Assignment/internal/Storage/sqlite"
	"github/Bharatjawa2/CtrlB_Assignment/internal/config"
	"github/Bharatjawa2/CtrlB_Assignment/internal/http/Handlers/courses"
	"github/Bharatjawa2/CtrlB_Assignment/internal/http/Handlers/student"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main(){
	// load config
	cfg:=config.MustLoad()
	// database setup
	storage,dberr:=sqlite.New(cfg)
	if dberr!=nil{
		log.Fatal(dberr)
	}
	slog.Info("Storage intialized", slog.String("env",cfg.Env),slog.String("version","1.0.0"))

	// setup router
	router:=http.NewServeMux()

	// Student
		router.HandleFunc("POST /api/students",student.Register(storage))
		router.HandleFunc("GET /api/students/{id}",student.GetById(storage))

	// Courses
		router.HandleFunc("POST /api/courses",courses.CreateCourse(storage))

	// Enrollment
		
	// setup server

	server:=http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}

	slog.Info("Server Started!!! %s",slog.String("address",cfg.HTTPServer.Addr))

	done:=make(chan os.Signal,1) // channel
	signal.Notify(done,os.Interrupt,syscall.SIGINT,syscall.SIGTERM)

	go func(){ // gorountine 
		err := server.ListenAndServe()
	if err!=nil{
		log.Fatal("failed to Start sever")
	}
	}()

	<- done

	slog.Info("Shutting down the server")

	ctx, cancel:=context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	err:=server.Shutdown(ctx) // Take Time , to solve use context
	if err!=nil{
		slog.Error("Failed to ShutDown Server",slog.String("error",err.Error()))
	}

	slog.Info("Shutting shutdown successfully!!!")
}