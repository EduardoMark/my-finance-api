# 💰 My Finance API

Uma API REST completa para gerenciamento de finanças pessoais desenvolvida em Go, com autenticação JWT e banco de dados PostgreSQL.

## 📋 Sobre o Projeto

A **My Finance API** é uma solução robusta para controle financeiro pessoal que permite aos usuários gerenciar suas finanças de forma organizada e segura. A API oferece funcionalidades completas para:

- 👤 **Gestão de usuários** com autenticação segura
- 🏦 **Contas bancárias** múltiplas por usuário
- 📊 **Categorias** personalizadas para organização
- 💳 **Transações** detalhadas com filtros avançados

## 🚀 Tecnologias Utilizadas

- **Go 1.24+** - Linguagem de programação
- **Chi Router** - Framework web minimalista e performático
- **PostgreSQL** - Banco de dados relacional
- **JWT** - Autenticação stateless
- **Docker & Docker Compose** - Containerização
- **SQLC** - Geração de código SQL type-safe
- **Tern** - Migrações de banco de dados
- **bcrypt** - Hash de senhas

## 🏗️ Arquitetura

O projeto segue os princípios de **Clean Architecture** com separação clara de responsabilidades:

```
├── cmd/                    # Pontos de entrada da aplicação
│   ├── api/               # Servidor principal da API
│   └── terndotenv/        # Utilitário para migrações
├── internal/              # Código interno da aplicação
│   ├── api/               # Configuração de rotas e API
│   ├── middlewares/       # Middlewares HTTP
│   ├── store/pgstore/     # Camada de dados PostgreSQL
│   ├── user/              # Módulo de usuários
│   ├── account/           # Módulo de contas
│   ├── category/          # Módulo de categorias
│   ├── transaction/       # Módulo de transações
│   └── validator/         # Validadores customizados
└── pkg/                   # Pacotes reutilizáveis
    ├── config/            # Configurações
    ├── database/          # Conexão com banco
    ├── token/             # Gestão de JWT
    └── httputils/         # Utilitários HTTP
```

## 🔧 Configuração

### Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:

```env
# Servidor
PORT=3000

# Banco de Dados
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=root
DB_NAME=my_finance_api
DB_TIMEZONE=UTC

# JWT
JWT_SECRET=your_super_secret_jwt_key_here
```

### Pré-requisitos

**Para Docker (Recomendado):**
- Docker
- Docker Compose

**Para instalação manual:**
- Go 1.24+
- PostgreSQL 14+

## 🚀 Como Executar

### Método Recomendado: Docker Compose

A forma mais simples de executar a aplicação é usando Docker Compose, que configura automaticamente tanto a API quanto o banco de dados:

```bash
# Clone o repositório
git clone https://github.com/EduardoMark/my-finance-api.git
cd my-finance-api

# Execute a aplicação completa
docker-compose up -d
```

✅ **Pronto!** A API estará disponível em `http://localhost:3000`

> O Docker Compose irá:
> - Configurar o banco PostgreSQL automaticamente
> - Executar as migrações necessárias
> - Iniciar a API na porta 3000
> - Configurar as variáveis de ambiente

## 🛠️ Instalação Manual (Opcional)

### 1. Clone e configure o projeto

```bash
git clone https://github.com/EduardoMark/my-finance-api.git
cd my-finance-api

# Instale as dependências
go mod tidy
```

### 2. Configure o banco de dados

```bash
# Execute as migrações
go run cmd/terndotenv/main.go
```

### 3. Execute a aplicação

```bash
# Desenvolvimento
go run cmd/api/main.go

# Ou compile e execute
go build -o bin/api cmd/api/main.go
./bin/api
```

## 📚 Documentação da API

### Autenticação

Todos os endpoints (exceto registro e login) requerem autenticação via JWT:

```http
Authorization: Bearer <jwt_token>
```

### Endpoints Principais

#### 👤 Usuários
```http
POST /api/v1/users/register    # Registro de usuário
POST /api/v1/users/login       # Login
GET  /api/v1/users/profile     # Perfil do usuário
```

#### 🏦 Contas
```http
POST   /api/v1/accounts        # Criar conta
GET    /api/v1/accounts        # Listar contas
GET    /api/v1/accounts/:id    # Obter conta específica
PUT    /api/v1/accounts/:id    # Atualizar conta
DELETE /api/v1/accounts/:id    # Deletar conta
```

#### 📊 Categorias
```http
POST   /api/v1/categories      # Criar categoria
GET    /api/v1/categories      # Listar categorias
GET    /api/v1/categories/:id  # Obter categoria específica
PUT    /api/v1/categories/:id  # Atualizar categoria
DELETE /api/v1/categories/:id  # Deletar categoria
```

#### 💳 Transações
```http
POST   /api/v1/transactions    # Criar transação
GET    /api/v1/transactions    # Listar transações (com filtros)
GET    /api/v1/transactions/:id # Obter transação específica
PUT    /api/v1/transactions/:id # Atualizar transação
DELETE /api/v1/transactions/:id # Deletar transação
```

### Exemplos de Uso

#### Criar uma transação
```bash
curl -X POST http://localhost:3000/api/v1/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -d '{
    "description": "Compra no supermercado",
    "amount": 150.50,
    "date": "2024-01-15",
    "type": "expense",
    "account_id": "123e4567-e89b-12d3-a456-426614174000",
    "category_id": "123e4567-e89b-12d3-a456-426614174001"
  }'
```

#### Listar transações com filtros
```bash
# Todas as transações
GET /api/v1/transactions

# Filtrar por tipo
GET /api/v1/transactions?type=expense

# Filtrar por conta
GET /api/v1/transactions?account_id=uuid

# Filtrar por período
GET /api/v1/transactions?start_date=2024-01-01&end_date=2024-01-31
```

## 🔒 Segurança

- **Autenticação JWT** com tokens que expiram em 1 hora
- **Hash bcrypt** para senhas com salt automático
- **Validação rigorosa** de dados de entrada
- **Middleware de autenticação** em todas as rotas protegidas
- **Isolamento por usuário** - cada usuário acessa apenas seus próprios dados

## 🧪 Estrutura de Dados

### Tipos de Transação
- `income` - Receita
- `expense` - Despesa

### Validações
- **UUIDs válidos** para todos os IDs
- **Valores monetários** devem ser positivos
- **Datas** no formato YYYY-MM-DD
- **Campos obrigatórios** validados automaticamente

## 🚦 Status de Resposta

| Código | Descrição |
|--------|-----------|
| 200 | Operação bem-sucedida |
| 201 | Recurso criado com sucesso |
| 204 | Operação bem-sucedida sem conteúdo |
| 400 | Dados inválidos |
| 401 | Não autorizado |
| 404 | Recurso não encontrado |
| 422 | Erro de validação |
| 500 | Erro interno do servidor |

## 🤝 Contribuindo

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -m 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## 👨‍💻 Autor

**Eduardo Mark**
- GitHub: [@EduardoMark](https://github.com/EduardoMark)

---

⭐ Se este projeto te ajudou, considere dar uma estrela no repositório!