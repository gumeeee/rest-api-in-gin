# 🚀 REST API em Go com Gin

Uma API REST robusta desenvolvida em Go utilizando o framework Gin, com sistema de autenticação JWT, banco de dados SQLite e documentação automática com Swagger.

## 📋 Índice

- [Sobre o Projeto](#sobre-o-projeto)
- [Tecnologias Utilizadas](#tecnologias-utilizadas)
- [Arquitetura do Projeto](#arquitetura-do-projeto)
- [Funcionalidades](#funcionalidades)
- [Estrutura do Banco de Dados](#estrutura-do-banco-de-dados)
- [Endpoints da API](#endpoints-da-api)
- [Instalação e Configuração](#instalação-e-configuração)
- [Como Executar](#como-executar)
- [Documentação da API](#documentação-da-api)
- [Variáveis de Ambiente](#variáveis-de-ambiente)
- [Desenvolvimento](#desenvolvimento)
- [Contribuição](#contribuição)
- [Licença](#licença)

## 🎯 Sobre o Projeto

Esta API REST foi desenvolvida para gerenciar eventos e participantes, oferecendo funcionalidades completas de autenticação, criação e gerenciamento de eventos, e controle de participantes. O projeto segue as melhores práticas de desenvolvimento em Go, incluindo uma arquitetura limpa, validação de dados, tratamento de erros e documentação automática.

### Principais Características

- ✅ **Autenticação JWT**: Sistema seguro de autenticação com tokens JWT
- ✅ **CRUD Completo**: Operações completas de criação, leitura, atualização e exclusão
- ✅ **Validação de Dados**: Validação robusta de entrada usando Gin Validator
- ✅ **Documentação Automática**: Swagger/OpenAPI integrado
- ✅ **Migrações de Banco**: Sistema de migrações para controle de versão do banco
- ✅ **Hot Reload**: Desenvolvimento com recarregamento automático usando Air
- ✅ **Arquitetura Limpa**: Separação clara de responsabilidades

## 🛠 Tecnologias Utilizadas

### Backend
- **[Go 1.24.2](https://golang.org/)** - Linguagem principal
- **[Gin](https://github.com/gin-gonic/gin)** - Framework web HTTP
- **[SQLite](https://www.sqlite.org/)** - Banco de dados
- **[JWT](https://github.com/golang-jwt/jwt)** - Autenticação com tokens
- **[bcrypt](https://golang.org/x/crypto/bcrypt)** - Hash de senhas
- **[Swagger](https://swaggo.github.io/swaggo/)** - Documentação da API
- **[Air](https://github.com/cosmtrek/air)** - Hot reload para desenvolvimento
- **[Golang Migrate](https://github.com/golang-migrate/migrate)** - Migrações de banco

### Ferramentas de Desenvolvimento
- **[Go Modules](https://go.dev/blog/using-go-modules)** - Gerenciamento de dependências
- **[Godotenv](https://github.com/joho/godotenv)** - Gerenciamento de variáveis de ambiente

## 🏗 Arquitetura do Projeto

O projeto segue uma arquitetura limpa e modular, organizada da seguinte forma:

```
rest-api-go/
├── cmd/                    # Pontos de entrada da aplicação
│   ├── api/               # Servidor principal da API
│   │   ├── main.go        # Ponto de entrada principal
│   │   ├── server.go      # Configuração do servidor
│   │   ├── routes.go      # Definição das rotas
│   │   ├── auth.go        # Handlers de autenticação
│   │   ├── events.go      # Handlers de eventos
│   │   ├── middleware.go  # Middlewares personalizados
│   │   └── context.go     # Contextos personalizados
│   └── migrate/           # Ferramenta de migração
│       ├── main.go        # Ponto de entrada das migrações
│       └── migrations/    # Arquivos de migração SQL
├── internal/              # Código interno da aplicação
│   ├── database/          # Camada de acesso a dados
│   │   ├── models.go      # Estrutura dos modelos
│   │   ├── users.go       # Operações de usuários
│   │   ├── events.go      # Operações de eventos
│   │   └── attendees.go   # Operações de participantes
│   └── env/               # Gerenciamento de variáveis de ambiente
│       └── env.go         # Funções de configuração
├── docs/                  # Documentação gerada pelo Swagger
├── go.mod                 # Dependências do Go
├── go.sum                 # Checksums das dependências
├── .air.toml             # Configuração do Air (hot reload)
└── README.md             # Este arquivo
```

### Padrões Arquiteturais

- **Separação de Responsabilidades**: Cada camada tem uma responsabilidade específica
- **Injeção de Dependência**: Dependências são injetadas através de construtores
- **Repository Pattern**: Acesso a dados abstraído através de interfaces
- **Middleware Pattern**: Funcionalidades cross-cutting através de middlewares
- **RESTful Design**: Endpoints seguindo padrões REST

## ⚡ Funcionalidades

### 🔐 Autenticação e Autorização
- Registro de usuários com validação de dados
- Login com JWT tokens
- Middleware de autenticação para rotas protegidas
- Hash seguro de senhas com bcrypt

### 📅 Gerenciamento de Eventos
- Criação, leitura, atualização e exclusão de eventos
- Validação de dados de entrada
- Controle de propriedade (apenas o criador pode editar/excluir)
- Busca de eventos por ID

### 👥 Gerenciamento de Participantes
- Adicionar participantes a eventos
- Remover participantes de eventos
- Listar participantes de um evento específico
- Listar eventos de um participante específico

### 📊 Relatórios e Consultas
- Listagem de todos os eventos
- Consultas relacionais entre eventos e participantes
- Dados estruturados em JSON

## 🗄 Estrutura do Banco de Dados

### Tabela `users`
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    password TEXT NOT NULL
);
```

### Tabela `events`
```sql
CREATE TABLE events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    date DATETIME NOT NULL,
    location TEXT NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);
```

### Tabela `attendees`
```sql
CREATE TABLE attendees (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    event_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES events (id) ON DELETE CASCADE
);
```

## 🔌 Endpoints da API

### Autenticação

| Método | Endpoint | Descrição | Autenticação |
|--------|----------|-----------|--------------|
| `POST` | `/api/v1/auth/register` | Registrar novo usuário | ❌ |
| `POST` | `/api/v1/auth/login` | Fazer login | ❌ |

### Eventos

| Método | Endpoint | Descrição | Autenticação |
|--------|----------|-----------|--------------|
| `GET` | `/api/v1/events` | Listar todos os eventos | ❌ |
| `GET` | `/api/v1/events/:id` | Buscar evento por ID | ❌ |
| `POST` | `/api/v1/events` | Criar novo evento | ✅ |
| `PUT` | `/api/v1/events/:id` | Atualizar evento | ✅ |
| `DELETE` | `/api/v1/events/:id` | Excluir evento | ✅ |

### Participantes

| Método | Endpoint | Descrição | Autenticação |
|--------|----------|-----------|--------------|
| `GET` | `/api/v1/events/:id/attendees` | Listar participantes de um evento | ❌ |
| `GET` | `/api/v1/attendees/:id/events` | Listar eventos de um participante | ❌ |
| `POST` | `/api/v1/events/:id/attendees/:userId` | Adicionar participante ao evento | ✅ |
| `DELETE` | `/api/v1/events/:id/attendees/:userId` | Remover participante do evento | ✅ |

### Documentação
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `GET` | `/swagger/*` | Documentação Swagger |

## 🚀 Instalação e Configuração

### Pré-requisitos

- Go 1.24.2 ou superior
- Git

### 1. Clone o repositório

```bash
git clone https://github.com/gumeeee/rest-api-go.git
cd rest-api-go
```

### 2. Instale as dependências

```bash
go mod download
```

### 3. Configure as variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto:

```env
PORT=8080
JWT_SECRET=sua-chave-secreta-aqui
```

### 4. Execute as migrações do banco

```bash
go run cmd/migrate/main.go
```

## ▶ Como Executar

### Desenvolvimento (com Hot Reload)

```bash
# Instale o Air globalmente (se ainda não tiver)
go install github.com/cosmtrek/air@latest

# Execute com hot reload
air
```

### Produção

```bash
# Compile o projeto
go build -o bin/api cmd/api/main.go

# Execute o binário
./bin/api
```

### Executar diretamente

```bash
go run cmd/api/main.go
```

## 📚 Documentação da API

A API possui documentação automática gerada pelo Swagger. Após iniciar o servidor, acesse:

- **Swagger UI**: http://localhost:8080/swagger/
- **JSON da API**: http://localhost:8080/swagger/doc.json
- **YAML da API**: http://localhost:8080/swagger/doc.yaml

### Exemplos de Uso

#### 1. Registrar um usuário

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "email": "joao@example.com",
    "password": "senha123456"
  }'
```

#### 2. Fazer login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "joao@example.com",
    "password": "senha123456"
  }'
```

#### 3. Criar um evento (com autenticação)

```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN_JWT" \
  -d '{
    "name": "Meetup Go",
    "description": "Encontro da comunidade Go",
    "date": "2024-01-15T19:00:00Z",
    "location": "São Paulo, SP"
  }'
```

#### 4. Listar todos os eventos

```bash
curl -X GET http://localhost:8080/api/v1/events
```

## 🔧 Variáveis de Ambiente

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `PORT` | Porta do servidor | `8080` |
| `JWT_SECRET` | Chave secreta para JWT | `secret-jwt-key-123456` |

## 💻 Desenvolvimento

### Estrutura de Desenvolvimento

O projeto utiliza o **Air** para hot reload durante o desenvolvimento. A configuração está no arquivo `.air.toml`.

### Comandos Úteis

```bash
# Executar testes
go test ./...

# Verificar cobertura de testes
go test -cover ./...

# Formatar código
go fmt ./...

# Verificar problemas de código
go vet ./...

# Gerar documentação Swagger
swag init -g cmd/api/main.go

# Executar migrações
go run cmd/migrate/main.go
```

### Padrões de Código

- **Nomenclatura**: Seguir convenções Go (camelCase para variáveis, PascalCase para exportados)
- **Tratamento de Erros**: Sempre verificar e tratar erros adequadamente
- **Validação**: Usar tags de binding do Gin para validação
- **Documentação**: Comentários Swagger para documentação da API
- **Logs**: Usar o logger padrão do Go para logs estruturados

## 🤝 Contribuição

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

### Diretrizes de Contribuição

- Siga os padrões de código estabelecidos
- Adicione testes para novas funcionalidades
- Atualize a documentação quando necessário
- Mantenha commits pequenos e descritivos

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 👨‍💻 Autor

**Gumeeee**

- GitHub: [@gumeeee](https://github.com/gumeeee)

## 🙏 Agradecimentos

- [Gin Framework](https://github.com/gin-gonic/gin) - Framework web HTTP
- [Swaggo](https://github.com/swaggo/swaggo) - Documentação automática
- [Air](https://github.com/cosmtrek/air) - Hot reload para Go

---

⭐ Se este projeto te ajudou, considere dar uma estrela no repositório!
