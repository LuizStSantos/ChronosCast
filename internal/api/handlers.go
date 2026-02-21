package api

import (
	"ChronosCast/internal/scheduler"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

// Estrutura que define o que recebemos no JSON
type StreamRequest struct {
	ID         string   `json:"id" binding:"required"`
	Input      string   `json:"entrada" binding:"required"`
	Output     string   `json:"saida" binding:"required"`
	DaysOfWeek []string `json:"dias_semana" binding:"required"`
	StartTime  string   `json:"hora_inicio" binding:"required"`
	EndTime    string   `json:"hora_fim" binding:"required"`
}

// APIContext guarda as referências globais para os handlers usarem
type APIContext struct {
	DB     *sql.DB
	Engine *scheduler.Engine
	JobIDs map[string]cron.EntryID
	Mu     *sync.Mutex
}

// HandlePost - Cria ou Atualiza um agendamento
func (api *APIContext) HandlePost(c *gin.Context) {
	var req StreamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Campos obrigatórios ausentes"})
		return
	}

	api.Mu.Lock()
	defer api.Mu.Unlock()

	// 1. Limpa se já existir (para atualizar)
	api.removerLogica(req.ID)

	// 2. Salva no Banco de Dados
	diasStr := strings.Join(req.DaysOfWeek, ",")
	_, err := api.DB.Exec("INSERT INTO streams (id, entrada, saida, dias, inicio, fim) VALUES (?, ?, ?, ?, ?, ?)",
		req.ID, req.Input, req.Output, diasStr, req.StartTime, req.EndTime)
	if err != nil {
		c.JSON(500, gin.H{"erro": "Falha ao salvar no banco"})
		return
	}

	// 3. Agenda no Motor (Scheduler)
	entryID, err := api.Engine.AddStreamJob(req.ID, req.Input, req.Output, req.StartTime, req.EndTime, req.DaysOfWeek)
	if err != nil {
		c.JSON(500, gin.H{"erro": "Falha ao agendar no cron"})
		return
	}

	api.JobIDs[req.ID] = entryID
	c.JSON(http.StatusOK, gin.H{"status": "sucesso", "id": req.ID})
}

// HandleDelete - Remove um agendamento pelo ID na URL
func (api *APIContext) HandleDelete(c *gin.Context) {
	id := c.Param("id")

	api.Mu.Lock()
	defer api.Mu.Unlock()

	// Aqui verificamos no banco também para garantir que apague mesmo que o job não esteja em memória
	api.removerLogica(id)
	c.JSON(http.StatusOK, gin.H{"status": "removido", "id": id})
}

// HandleList - Retorna a lista de todos os agendamentos ativos no banco
func (api *APIContext) HandleList(c *gin.Context) {
	api.Mu.Lock()
	defer api.Mu.Unlock()

	rows, err := api.DB.Query("SELECT id, entrada, saida, dias, inicio, fim FROM streams")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao consultar o banco de dados"})
		return
	}
	defer rows.Close()

	var agendamentos []StreamRequest
	for rows.Next() {
		var s StreamRequest
		var dias string
		err := rows.Scan(&s.ID, &s.Input, &s.Output, &dias, &s.StartTime, &s.EndTime)
		if err != nil {
			continue
		}
		// Converte a string "MON,TUE" de volta para ["MON", "TUE"]
		s.DaysOfWeek = strings.Split(dias, ",")
		agendamentos = append(agendamentos, s)
	}

	// Retorna lista vazia [] em vez de null se não houver nada
	if agendamentos == nil {
		agendamentos = []StreamRequest{}
	}

	c.JSON(http.StatusOK, agendamentos)
}

// Função auxiliar interna para evitar repetição de código
func (api *APIContext) removerLogica(id string) {
	if entryID, ok := api.JobIDs[id]; ok {
		api.Engine.Cron.Remove(entryID)
		delete(api.JobIDs, id)
	}
	api.DB.Exec("DELETE FROM streams WHERE id = ?", id)
}

// LoadSchedules - Lê o banco e reativa todos os agendamentos no Cron
func (api *APIContext) LoadSchedules() {
	api.Mu.Lock()
	defer api.Mu.Unlock()

	rows, err := api.DB.Query("SELECT id, entrada, saida, dias, inicio, fim FROM streams")
	if err != nil {
		fmt.Println("❌ Erro ao carregar agendamentos do banco:", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var s StreamRequest
		var dias string
		rows.Scan(&s.ID, &s.Input, &s.Output, &dias, &s.StartTime, &s.EndTime)
		s.DaysOfWeek = strings.Split(dias, ",")

		// Reagenda no Motor
		entryID, err := api.Engine.AddStreamJob(s.ID, s.Input, s.Output, s.StartTime, s.EndTime, s.DaysOfWeek)
		if err == nil {
			api.JobIDs[s.ID] = entryID
			count++
		}
	}
	fmt.Printf("📦 Recuperação: %d agendamentos reativados com sucesso!\n", count)
}
