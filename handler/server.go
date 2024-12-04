package handler

import (
	"database/sql"
	"fmt"
	"roxy/config"
	"roxy/repository"
	"roxy/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	barangUc    usecase.MstBarangUseCase
	transaksiUc usecase.TransaksiUsecase

	engine *gin.Engine
	host   string
}

func (s *Server) initRoute() {
	rg := s.engine.Group(config.ApiGroup)

	NewBarangHandler(s.barangUc, rg).Route()
	NewTransaksiHandler(s.transaksiUc, rg).Route()
}

func (s *Server) Run() {
	s.initRoute()
	if err := s.engine.Run(s.host); err != nil {
		panic(fmt.Errorf("server not running on host %s, becauce error %v", s.host, err.Error()))
	}
}

func NewServer() *Server {
	cfg, _ := config.NewConfig()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		fmt.Println("connection error", err)
	}

	//inject dependencies repo layer
	barangRepo := repository.NewBarangRepository(db)
	transaksiRepo := repository.NewTransaksiRepository(db)
	//inject dependencies usecase layer
	barangUc := usecase.NewBarangUseCase(barangRepo)
	transaksiUc := usecase.NewTransaksiUsecase(transaksiRepo)

	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)
	return &Server{
		barangUc:    barangUc,
		transaksiUc: transaksiUc,

		engine: engine,
		host:   host,
	}
}
