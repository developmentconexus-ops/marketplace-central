// Package composition liga o job periódico de pedidos (orders/application's
// NewOrdersJob) ao scheduler de sync — uma instância de
// sync/application.Scheduler por instalação Mercado Livre ativa, seguindo o
// mesmo padrão que listings/composition.Group já estabeleceu.
//
// O fan-out por instalação (e não um scheduler compartilhado) é imposto pelo
// próprio sync: syncapp.Scheduler nasce amarrado a um installationID e
// RegisterJob aceita uma registração por entidade por instância. É a mesma
// razão que listings/composition documenta.
//
// Reusar o scheduler em vez de um ticker próprio é o que faz a falha aparecer:
// RecordFailure grava consecutive_failures em sync_state, /sync/health devolve
// a entidade sem allowlist, e SyncHealthCard pinta a linha de vermelho. Um
// ticker paralelo custaria a mesma quantidade de código e deixaria a tela
// verde com o token expirado.
package composition

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	integrationsdomain "marketplace-central/apps/server_core/internal/modules/integrations/domain"
	ordersapp "marketplace-central/apps/server_core/internal/modules/orders/application"
	syncapp "marketplace-central/apps/server_core/internal/modules/sync/application"
	synccomposition "marketplace-central/apps/server_core/internal/modules/sync/composition"
	"marketplace-central/apps/server_core/internal/modules/sync/domain"
)

// mercadoLivreProviderCode é o único provider para o qual este scheduler roda
// o job de pedidos hoje — mesma literal que listings/composition usa.
const mercadoLivreProviderCode = "mercado_livre"

// installationLister enumera as instalações do tenant — mesmo seam que
// listings/composition.installationLister usa, implementado por
// *integrationsapp.InstallationService.
type installationLister interface {
	List(ctx context.Context) ([]integrationsdomain.Installation, error)
}

// Schedulers é um *syncapp.Scheduler por instalação Mercado Livre ativa.
type Schedulers []*syncapp.Scheduler

// StartAll sobe o ticker de cada scheduler em sua própria goroutine e volta
// imediatamente.
func (s Schedulers) StartAll(ctx context.Context) {
	for _, scheduler := range s {
		scheduler := scheduler
		go scheduler.Start(ctx)
	}
}

// Group captura os inputs de composição do scheduler de pedidos. Construir
// via NewOrdersSchedulers NÃO toca o banco — instalações são resolvidas
// preguiçosamente, dentro da goroutine de StartAll, nunca durante a
// composição síncrona. Mesmo contrato que todo outro composto em root.go
// respeita (provado por
// TestNewRootRouterBuildsWithoutLegacyConnectorCredentials).
type Group struct {
	pool          *pgxpool.Pool
	interval      time.Duration
	pageSize      int
	overlap       time.Duration
	installations installationLister
	importer      ordersapp.OrdersImporter
}

// NewOrdersSchedulers captura os inputs necessários para construir um
// scheduler de pedidos (Task 5's NewOrdersJob, registrado em
// domain.EntityOrders) por instalação Mercado Livre ativa. Nada é resolvido
// ou consultado aqui; chame StartAll para de fato listar instalações e subir
// os schedulers por instalação.
func NewOrdersSchedulers(
	pool *pgxpool.Pool,
	interval time.Duration,
	pageSize int,
	overlap time.Duration,
	installations installationLister,
	importer ordersapp.OrdersImporter,
) Group {
	return Group{
		pool:          pool,
		interval:      interval,
		pageSize:      pageSize,
		overlap:       overlap,
		installations: installations,
		importer:      importer,
	}
}

// StartAll resolve as instalações (via o seam installationLister) e sobe um
// scheduler de pedidos por instalação Mercado Livre ativa encontrada, tudo
// dentro de uma goroutine de fundo para nunca bloquear ou falhar o boot do
// servidor (a chamada de composição em root.go é
// NewOrdersSchedulers(...).StartAll(ctx)).
//
// Uma falha em List, ou qualquer panic ao resolvê-la (ex.: pool não
// configurado), é logada e não produz nenhum scheduler em vez de derrubar o
// processo (mesma postura "log, não crashe" que
// sync/composition.NewProductsScheduler já tinha para seu próprio
// RegisterJob). É um fail-open deliberado: um restart do operador tenta de
// novo; nada fabrica um agendamento silenciosamente nesse meio-tempo.
//
// Instalações resolvidas aqui são um snapshot ÚNICO tirado nesta chamada —
// uma instalação ML conectada DEPOIS do boot não ganha seu próprio scheduler
// até o próximo restart.
func (g Group) StartAll(ctx context.Context) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("orders scheduler: panic resolvendo instalacoes", "panic", r)
			}
		}()
		for _, scheduler := range resolveOrdersSchedulers(ctx, g) {
			scheduler := scheduler
			go scheduler.Start(ctx)
		}
	}()
}

// resolveOrdersSchedulers é a lógica síncrona de fan-out que a goroutine de
// StartAll roda: listar instalações, filtrar para Mercado Livre, construir um
// Scheduler + job EntityOrders registrado por instalação encontrada. Separado
// de StartAll para que os testes exerçam o fan-out diretamente com um
// installationLister falso (que nunca toca um banco real, portanto seguro de
// chamar sincronamente em teste).
func resolveOrdersSchedulers(ctx context.Context, g Group) Schedulers {
	all, err := g.installations.List(ctx)
	if err != nil {
		slog.Error("orders scheduler: listar instalacoes falhou", "err", err)
		return nil
	}

	var schedulers Schedulers
	for _, installation := range all {
		if installation.ProviderCode != mercadoLivreProviderCode {
			continue
		}
		scheduler := synccomposition.NewInstallationScheduler(g.pool, installation.TenantID, installation.InstallationID, g.interval)
		// EntityOrders é válida e é a única (primeira) registração nesta
		// instância, então RegisterJob não pode falhar; os modos de falha dele
		// são cobertos genericamente em sync/application/scheduler_test.go.
		_ = scheduler.RegisterJob(domain.EntityOrders,
			ordersapp.NewOrdersJob(g.importer, installation.InstallationID, g.pageSize, g.overlap, time.Now))
		schedulers = append(schedulers, scheduler)
	}
	return schedulers
}
