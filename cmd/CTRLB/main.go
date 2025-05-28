package main

import (
	"context"
	"github/Bharatjawa2/CtrlB_Assignment/internal/Storage/sqlite"
	"github/Bharatjawa2/CtrlB_Assignment/internal/config"
	"github/Bharatjawa2/CtrlB_Assignment/internal/http/Handlers/admin"
	"github/Bharatjawa2/CtrlB_Assignment/internal/http/Handlers/courses"
	"github/Bharatjawa2/CtrlB_Assignment/internal/http/Handlers/enrollment"
	"github/Bharatjawa2/CtrlB_Assignment/internal/http/Handlers/student"
	"github/Bharatjawa2/CtrlB_Assignment/internal/http/middlewares"
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

	// Admin
		router.HandleFunc("POST /api/admin",admin.LoginAdmin(*cfg))
		router.HandleFunc("POST /api/admin/logout",middlewares.AdminMiddleware(cfg.JWTSecret,admin.Logout()))

	// Student
		router.HandleFunc("POST /api/students",student.Register(storage))
		router.HandleFunc("POST /api/students/login",student.LoginStudent(storage,*cfg))
		router.HandleFunc("GET /api/students/{id}",middlewares.AdminMiddleware(cfg.JWTSecret,student.GetById(storage)))
		router.HandleFunc("GET /api/students/all",middlewares.AdminMiddleware(cfg.JWTSecret,student.GetAllStudents(storage)))
		router.HandleFunc("GET /api/students",middlewares.AdminMiddleware(cfg.JWTSecret,student.GetStudentByEmail(storage)))
		router.HandleFunc("PUT /api/students/update",middlewares.StudentMiddleware(cfg.JWTSecret, student.UpdateStudent(storage)))
		router.HandleFunc("POST /api/students/logout",middlewares.StudentMiddleware(cfg.JWTSecret,student.Logout()))

	// Courses
		router.HandleFunc("POST /api/courses",middlewares.AdminMiddleware(cfg.JWTSecret,courses.CreateCourse(storage)))
		router.HandleFunc("GET /api/courses/{id}",courses.GetCourseById(storage))
		router.HandleFunc("GET /api/courses/all",courses.GetAllCourses(storage))
		router.HandleFunc("PUT /api/courses/update/{id}",middlewares.AdminMiddleware(cfg.JWTSecret,courses.UpdateCourse(storage)))
		router.HandleFunc("GET /api/courses/search",courses.SearchCoursesByName(storage))

	// Enrollment
		router.HandleFunc("POST /api/enrollment",middlewares.StudentMiddleware(cfg.JWTSecret,enrollment.EnrollStudent(storage)))
		router.HandleFunc("POST /api/unenrollment", middlewares.StudentMiddleware(cfg.JWTSecret,enrollment.UnenrollStudent(storage)))
		router.HandleFunc("GET /api/enrolled/students/{id}",middlewares.StudentMiddleware(cfg.JWTSecret,enrollment.GetCoursesByStudentID(storage)))
		router.HandleFunc("GET /api/enrolled/courses/{id}",middlewares.AdminMiddleware(cfg.JWTSecret,enrollment.GetStudentsByCourseID(storage)))

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