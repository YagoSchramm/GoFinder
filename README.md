# GoFinder

GoFinder e uma API em Go para buscar produtos em lojas brasileiras, normalizar os resultados, ordenar por preco e manter historico para alertas futuros.

## Status atual

O projeto ainda esta em fase inicial. O que ja existe no repositorio:

- Base da API em Go com `gorilla/mux`
- Conexao com PostgreSQL via `pgx`
- Autenticacao com cadastro, login e JWT
- Middleware de sessao autenticada
- Migracao inicial da tabela `users`
- Leitura simples de variaveis a partir de `.env`

O restante da proposta abaixo esta marcado como TODO enquanto ainda nao estiver implementado.

## Entidades

As tabelas do projeto ficam nas migrations em `infrastructure/script/migrate`.

Migration atual: `infrastructure/script/migrate/001-create-tables.up.sql`

| Entidade | Descricao | Status |
| --- | --- | --- |
| `users` | Usuarios cadastrados para autenticacao | Implementado |
| `searches` | Buscas feitas pelos usuarios | Implementada |
| `products` | Produtos encontrados em cada busca | Implementada |
| `price_history` | Historico de preco por produto e loja | Implementada |
| `alerts` | Alertas de preco por usuario | Implementada |

## Endpoints

### Implementados

| Metodo | Rota | Descricao |
| --- | --- | --- |
| POST | `/auth/register` | Cadastra usuario e retorna JWT |
| POST | `/auth/login` | Autentica usuario e retorna JWT |

### TODO

| Metodo | Rota | Descricao |
| --- | --- | --- |
| POST | `/search` | Dispara busca nos scrapers |
| GET | `/search/:id/results` | Retorna resultados ordenados por preco |
| GET | `/searches` | Lista historico de buscas do usuario |
| POST | `/alerts` | Cria alerta de preco: avisa se baixar de R$ X |
| GET | `/alerts` | Lista alertas ativos |

## Diferenciais planejados

| Recurso | Descricao | Status |
| --- | --- | --- |
| Goroutines paralelas | Cada scraper roda ao mesmo tempo para reduzir o tempo de resposta | TODO |
| Cache por query | A mesma busca feita nos ultimos 30 minutos retorna do banco | TODO |
| Normalizacao de titulo | Ex.: "Notebook Dell Inspiron 15" e "Dell Inspiron 15 Notebook" viram o mesmo produto | TODO |
| Alerta de preco | Job periodico re-raspa produtos e notifica por email se o preco baixou | TODO |
| Score de relevancia | Ranking por relevancia do titulo + preco, nao apenas por preco bruto | TODO |

## Bibliotecas Go previstas

| Lib | Uso | Status |
| --- | --- | --- |
| `colly` | Scraping principal | TODO |
| `chromedp` | Sites com JS pesado como Shopee e Americanas | TODO |
| `robfig/cron` | Job periodico de alertas | TODO |
| `redis` | Cache de buscas | TODO |

Bibliotecas ja presentes:

| Lib | Uso |
| --- | --- |
| `gorilla/mux` | Roteamento HTTP |
| `pgx` | PostgreSQL |
| `golang-jwt/jwt` | JWT |
| `golang.org/x/crypto` | Hash de senha |
| `google/uuid` | UUIDs |

## Configuracao

Use o arquivo `.env-example` como base para criar um `.env` na raiz do projeto.

## Como rodar

```bash
go mod download
go run .
```

Antes de subir a API, aplique as migracoes em `infrastructure/script/migrate`.
