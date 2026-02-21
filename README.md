
# 🛰️ CHRONOSCAST
## Sistema de Agendamento Profissional

**Autor:** Luiz Stormorwski dos Santos  
**Versão:** 1.0.0  
**Licença:** MIT  
**Ambiente:** Fedora Linux / Go (Golang)

---

## ⚠️ Dependência Obrigatória

O **ChronosCast utiliza o FFmpeg como motor principal de transmissão**.

A instalação do **FFmpeg é obrigatória** para o funcionamento do sistema.  
Sem ele, as transmissões não serão executadas.

Verifique se está instalado:

```bash
ffmpeg -version
```

Instalação no Fedora:

```bash
sudo dnf install ffmpeg
```

---

## 📦 Requisitos do Sistema

- Go 1.21+
- SQLite
- FFmpeg (obrigatório)
- Linux (testado em Fedora)

---

## 📌 1. Sobre o Projeto

O **ChronosCast** é um motor de automação para transmissões de vídeo.

Ele utiliza:

- **SQLite** → Persistência de dados  
- **Cron** → Agendamento com precisão  
- **FFmpeg** → Processamento de streamings SRT/RTMP  

---

## 🚀 2. Como Executar

No terminal, no projeto, execute:

```bash
./chronoscast-linux-amd64
```

---

## 🌐 3. API – Endpoints Disponíveis

### 🔹 A) Agendar / Atualizar Transmissão

**Método:** POST  
**URL:** http://localhost:8080/agendar  

Exemplo de Corpo JSON:

```json
{
  "id": "transmissao_01",
  "entrada": "srt://URL_DE_ENTRADA",
  "saida": "rtmp://URL_DE_SAIDA",
  "dias_semana": ["MON", "WED", "FRI"],
  "hora_inicio": "14:00:00",
  "hora_fim": "16:00:00"
}
```

---

### 🔹 B) Listar Agendamentos

**Método:** GET  
**URL:** http://localhost:8080/agendar  

Retorna todos os agendamentos salvos no banco de dados.

---

### 🔹 C) Deletar Agendamento

**Método:** DELETE  
**URL:** http://localhost:8080/agendar/{id_da_transmissao}

Remove o agendamento do banco e cancela o temporizador imediatamente.

---

## 📅 4. Formato dos Dias da Semana

SUN – Domingo  
MON – Segunda  
TUE – Terça  
WED – Quarta  
THU – Quinta  
FRI – Sexta  
SAT – Sábado  

---

## ⚙️ 5. Recursos Automáticos

### Auto-Load
Ao iniciar o servidor, o sistema lê o banco \`chronos.db\` e reativa automaticamente todos os agendamentos.

### Logs de Tempo
O terminal exibe:
- Hora exata de início  
- Hora exata de término  
- Duração real da transmissão  

---

## 📂 6. Estrutura de Arquivos

/cmd/server/main.go  
/internal/api/  
/internal/scheduler/  
/internal/storage/  
/chronos.db  
/LICENSE  

---

> "Código é poesia, automação é liberdade."