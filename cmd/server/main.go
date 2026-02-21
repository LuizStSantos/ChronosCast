package main

import (
	"ChronosCast/internal/api"
	"ChronosCast/internal/scheduler"
	"ChronosCast/internal/storage"
	"fmt"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	fmt.Println(`

 ██████╗██╗  ██╗██████╗  ██████╗ ███╗   ██╗ ██████╗ ███████╗ ██████╗ █████╗ ███████╗████████╗
██╔════╝██║  ██║██╔══██╗██╔═══██╗████╗  ██║██╔═══██╗██╔════╝ ██╔════╝██╔══██╗██╔════╝╚══██╔══╝
██║     ███████║██████╔╝██║   ██║██╔██╗ ██║██║   ██║███████╗ ██║     ███████║███████╗   ██║   
██║     ██╔══██║██╔══██╗██║   ██║██║╚██╗██║██║   ██║╚════██║ ██║     ██╔══██║╚════██║   ██║   
╚██████╗██║  ██║██║  ██║╚██████╔╝██║ ╚████║╚██████╔╝███████║ ╚██████╗██║  ██║███████║   ██║   
 ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝ ╚═════╝ ╚══════╝ ╚═════╝╚═╝  ╚═╝╚══════╝   ╚═╝   

API de Agendamento Profissional | v1.0.0
Desenvolvido por: Luiz Stormorwski dos Santos

	`)

	// Inicializa Componentes
	db, err := storage.InitDB("./chronos.db")
	if err != nil {
		log.Fatal("❌ Erro fatal ao abrir o banco de dados: ", err)
	}
	engine := scheduler.NewEngine()

	// 2. Cria o Contexto da API
	apiCtx := &api.APIContext{
		DB:     db,
		Engine: engine,
		JobIDs: make(map[string]cron.EntryID),
		Mu:     &sync.Mutex{},
	}

	apiCtx.LoadSchedules()

	r := gin.Default()

	// 3. Rotas conectadas aos Handlers
	r.GET("/agendar", apiCtx.HandleList)
	r.POST("/agendar", apiCtx.HandlePost)
	r.DELETE("/agendar/:id", apiCtx.HandleDelete)

	r.Run(":8080")
}
