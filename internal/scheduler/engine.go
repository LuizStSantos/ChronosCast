package scheduler

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

type Engine struct {
	Cron *cron.Cron
}

func NewEngine() *Engine {
	e := cron.New(cron.WithSeconds())
	e.Start()
	return &Engine{Cron: e}
}

func (e *Engine) AddStreamJob(id, input, output, startTime, endTime string, days []string) (cron.EntryID, error) {
	t := strings.Split(startTime, ":")
	cronExp := fmt.Sprintf("%s %s %s * * %s", t[2], t[1], t[0], strings.Join(days, ","))

	layout := "15:04:05"
	startT, _ := time.Parse(layout, startTime)
	endT, _ := time.Parse(layout, endTime)
	duration := endT.Sub(startT)
	if duration < 0 {
		duration += 24 * time.Hour
	}

	return e.Cron.AddFunc(cronExp, func() {
		horaInicioJob := time.Now()
		fmt.Printf("🎬 [%s] Iniciando rotina de agendamento às %s\n", id, horaInicioJob.Format("15:04:05"))

		// O contexto garante que o FFmpeg pare no horário de fim, não importa as tentativas
		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		for {
			// Verifica se o tempo total do agendamento já acabou
			if ctx.Err() != nil {
				break
			}

			cmd := exec.CommandContext(ctx, "ffmpeg",
				"-hide_banner",
				"-loglevel", "error",
				"-rw_timeout", "10000000",
				"-analyzeduration", "5000000",
				"-probesize", "5000000",
				"-fflags", "+genpts+igndts",
				"-i", input,
				"-c", "copy",
				"-f", "mpegts",
				"-reconnect", "1",
				"-reconnect_at_eof", "1",
				"-reconnect_streamed", "1",
				"-reconnect_delay_max", "5",
				output,
			)

			// Executa e aguarda a saída
			out, err := cmd.CombinedOutput()

			if err != nil {
				// Se o erro foi por causa do Timeout (fim do horário), encerramos o loop silenciosamente
				if ctx.Err() != nil {
					break
				}

				// Caso contrário, foi um erro de conexão (o I/O Error )
				fmt.Printf("⚠️  [%s] Conexão falhou. Tentando reconectar em 5 segundos...\n", id)

				// Opcional: Logar o erro técnico para debug
				lines := strings.Split(strings.TrimSpace(string(out)), "\n")
				if len(lines) > 0 {
					fmt.Printf("📝 [%s] Erro técnico: %s\n", id, lines[len(lines)-1])
				}

				// Espera antes de tentar novamente para não travar o processador
				time.Sleep(5 * time.Second)
				continue
			}

			// Se o FFmpeg encerrar sem erro (raro em streams infinitas), sai do loop
			break
		}

		horaFim := time.Now()
		duracaoReal := horaFim.Sub(horaInicioJob)
		fmt.Printf("🏁 [%s] Transmissão finalizada às %s\n", id, horaFim.Format("15:04:05"))
		fmt.Printf("⏱️  [%s] Duração total do processo: %v\n", id, duracaoReal.Round(time.Second))
	})
}
