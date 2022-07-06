# Aplicação de Exemplo do Livro [Let's Go! Learn to Build Professional Web Applications With Golang](https://lets-go.alexedwards.net/)

## Testes

**Parando a execução do teste na primeira falha:**

```bash
go test -failfast -v ./cmd/web
```

**Detalhando para saber qual teste está sendo executado, use `-v`:**

```bash
go test -v ./cmd/web/
go test -v ./pkg/models/mysql
```

**Use `-run` para especificar uma expressão regex para excutar apenas os testes dessa expressão:**

```bash
go test -v ="TestSignupUser" ./cmd/web/
```

**Limpando os caches dos resultados dos testes:**

```bash
go clean -testcache
```

**Excluíndo testes utilizando a opção `-short`. No código deve ser feito `if testing.Short()`:**

```bash
go test -v -short ./...
```

### Profiling

**Para saber a cobertura de testes no código, use `-cover`:**

```bash
go test -v -short -cover ./...
```

**Podemos obter uma análise mais detalhada da cobertura de teste por método e função usando a opção de `-coverprofile` da seguinte forma:**

```bash
go test -coverprofile=/tmp/profile.out ./...
```

Você pode então visualizar o perfil de cobertura usando o comando `go tool cover` como da seguinte forma:


```bash
go tool cover -func=/tmp/profile.out
```

**Uma maneira alternativa e mais visual de visualizar o perfil de cobertura é usar a opção de `-html` em vez de `-func`:**

```bash
go tool cover -html=/tmp/profile.out
```

### Adicionando modo de debug

**Ao executar no modo de depuração, quaisquer erros detalhados e rastreamentos da *stack* devem ser exibidos no navegador:**

```
go run ./cmd/web -debug
```
