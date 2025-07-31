# 🐳 DOCUMENTAÇÃO DO DOCKER

## 📚 Sumário

- [📦 Conceitos do Docker e Container](#-conceitos-do-docker-e-container)
- [⚛️ Container do Frontend](#️-container-do-frontend)
- [🐹 Container do Backend](#-container-do-backend)
- [🐬 Container do MySQL](#-container-do-mysql)
- [🚀 Como rodar o projeto](#-como-rodar-o-projeto)
- [💡 Dicas úteis](#-dicas-úteis)

---

## 📦 Conceitos do Docker e Container

O Docker é uma plataforma que permite empacotar, distribuir e executar aplicações de forma isolada em containers. Com ela, podemos rodar o mesmo ambiente de desenvolvimento em qualquer máquina, como um PC, servidor ou nuvem.

Esse ambiente de desenvolvimento encapsulado pelo Docker se chama **Container**. Ele inclui tudo que o projeto precisa para funcionar: sistema operacional, dependências, variáveis, portas, etc.

Dessa forma, ao executar o projeto em um novo dispositivo, **não é necessário configurar dependências manualmente**. Basta clonar o projeto e executar o comando Docker correto.

---

## ⚛️ Container do Frontend

### Arquivo (dockerfile) do frontend

```dockerfile
FROM node:18

WORKDIR /app
COPY . .

RUN npm install
EXPOSE 3000
CMD ["npm", "run", "dev", "--", "--host"]
```

### Estrutura

- Diretório: `./OficinaMecanica`
- Ferramentas: Vite, React
- Porta padrão: `5173` (ou `3000`, se alterado)

---

## 🐹 Container do Backend

### Arquivo (dockerfile) do backend

```dockerfile
FROM golang:1.21

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o main .

EXPOSE 8080
CMD ["./main"]
```

### Estrutura

- Diretório: `./oficina_backend`
- Ferramenta: Go
- Porta padrão: `8080`

---

## 🐬 Container do MySQL

O banco de dados MySQL é executado em um container próprio, utilizando a imagem oficial do MySQL. Não é necessário criar um Dockerfile para ele, apenas configurar no `docker-compose.yaml`.

### Configuração no docker-compose.yaml

```yaml
mysql:
  image: mysql:8.0
  container_name: oficina-mysql
  restart: unless-stopped
  environment:
    MYSQL_ROOT_PASSWORD: senha123
    MYSQL_DATABASE: oficina
  ports:
    - "3306:3306"
  networks:
    - oficina-net
  volumes:
    - mysql_data:/var/lib/mysql
```

### Estrutura

- Imagem: `mysql:8.0`
- Porta padrão: `3306`
- Variáveis de ambiente:
  - `MYSQL_ROOT_PASSWORD`: senha do usuário root
  - `MYSQL_DATABASE`: nome do banco de dados inicial
- Volume persistente: `mysql_data` (mantém os dados mesmo após reiniciar/remover o container)
- Rede: `oficina-net` (permite comunicação com backend e frontend)

> **Dica:**  
> O backend deve usar `mysql` como host do banco de dados, pois é o nome do serviço na rede Docker.

---

## 🚀 Como rodar o projeto

1. Certifique-se de ter o [Docker](https://www.docker.com/) e o [Docker Compose](https://docs.docker.com/compose/) instalados. Se estiver utilizando Scoop, use os comandos "scoop install docker" e "scoop install docker-compose".
2. Clone o repositório do projeto.
3. Navegue até o diretório onde está o arquivo `docker-compose.yaml` (normalmente `OficinaMecanica`).
4. Execute o comando abaixo para construir e iniciar todos os containers:

```sh
docker-compose up --build
```

5. Para rodar em segundo plano (background):

```sh
docker-compose up --build -d
```

6. Para parar todos os containers:

```sh
docker-compose down
```

---

## 💡 Dicas úteis

- Para acessar o banco de dados MySQL, use um cliente como DBeaver, MySQL Workbench ou o próprio terminal, conectando em `localhost:3306` com usuário `root` e a senha definida em `MYSQL_ROOT_PASSWORD`.
- Se precisar alterar variáveis de ambiente, faça isso diretamente no `docker-compose.yaml`.
- O frontend estará disponível em `http://localhost:3000` e o backend em `http://localhost:8080` após o Docker subir os containers.
