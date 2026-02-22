
# 🛰️ CHRONOSCAST

Motor de Automação para Transmissões Profissionais

O ChronosCast é um motor de agendamento e processamento de fluxos de áudio e vídeo. Ele atua como um orquestrador para o FFmpeg, permitindo programar entradas e saídas de streaming com persistência em banco de dados e precisão de cronograma.

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
./chronoscast-linux-amd64-v1.1.0
```

---

## 🌐 3. API – Endpoints Disponíveis

### 🔹 A) Agendar / Atualizar Transmissão

**Método:** POST  
**URL:** http://localhost:8080/agendar 

Cria ou atualiza uma transmissão

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

Lista todos os agendamentos ativos

---

### 🔹 C) Deletar Agendamento

**Método:** DELETE  
**URL:** http://localhost:8080/agendar/{id_da_transmissao}

Remove e cancela uma transmissão

---

### 🔹 d) Status

**Método:** GET  
**URL:** http://localhost:8080/status

Verifica o status do motor e versão

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

## ⚙️ 5. Diferenciais Técnicos

### Auto-Recovery:
Ao iniciar, o sistema lê o banco chronos.db e reativa todos os agendamentos automaticamente.

### Precisão Cron:
Gerenciamento nativo de dias da semana (SUN, MON, TUE, WED, THU, FRI, SAT).

### Observabilidade:
Logs detalhados informando duração real e status do processo FFmpeg no terminal.

### Graceful Handling:
O cancelamento de um agendamento interrompe o processo do FFmpeg de forma limpa.

---

## 📂 6. Estrutura de Arquivos

```
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   └── handlers.go
│   ├── models/
│   │   └── models.go
│   ├── scheduler/
│   │   └── engine.go
│   └── storage/
│       └── sqlite.go
├── .gitignore
├── build.sh
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```
---

> "Código é poesia, automação é liberdade."