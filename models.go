package main

type StreamSchedule struct {
	ID         string   `json:"id"`
	Input      string   `json:"entrada"`
	Output     string   `json:"saida"`
	DaysOfWeek []string `json:"dias_semana"`
	StartTime  string   `json:"hora_inicio"`
	EndTime    string   `json:"hora_fim"`
}