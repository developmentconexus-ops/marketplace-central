package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	syncapp "marketplace-central/apps/server_core/internal/modules/sync/application"
)

// icmsMatrixResolver é a metade de leitura da sincronização da matriz: lê o
// TGFICM da origem fixa e devolve uma célula resolvida por (uf_destino,
// grupo). Interface declarada aqui, no consumidor, para esta camada nunca
// importar adapters/oracle — quem satisfaz hoje é
// internal_read/adapters/oracle.ICMSMatrixReader.
type icmsMatrixResolver interface {
	ResolveCells(ctx context.Context) ([]domain.ICMSMatrixCell, error)
}

// icmsMatrixApplier é a metade de escrita: aplica as células ao espelho
// versionado icms_matrix_mirror e devolve quantas linhas escreveu. Quem
// satisfaz hoje é internal_read/adapters/mirror.ICMSMatrixWriter.
type icmsMatrixApplier interface {
	ApplyCells(ctx context.Context, tenantID string, cells []domain.ICMSMatrixCell) (int, error)
}

// ICMSMatrixCursor é o cursor persistido em sync_state.cursor. Registra o que
// o ciclo fez, não que ele ocorreu: um scheduler que tica com Cells=0 é
// distinguível de um que sincronizou de fato. Ambiguos é contado à parte
// porque uma célula ambígua é gravada sem alíquota (ADR-17) e portanto é uma
// pendência fiscal silenciosa no cálculo do pedido — o operador precisa ver o
// número subir.
type ICMSMatrixCursor struct {
	Cells       int       `json:"cells"`
	Ambiguos    int       `json:"ambiguos"`
	Written     int       `json:"written"`
	CompletedAt time.Time `json:"completed_at"`
}

// NewICMSMatrixJob monta o corpo do ciclo de sincronização da matriz de ICMS.
//
// Existe porque o P2.b entregou o leitor e o escritor sem nenhum relógio: os
// dois só tinham chamador em _test.go, icms_matrix_mirror ficou com 0 linhas, e
// o consumidor inteiro (root.go -> orders/adapters/pricingtax ->
// pricing/adapters/postgres.MatrixReader) leu a tabela vazia e devolveu
// Found:false, deixando margem NULL em 39 de 39 pedidos. O ADR-17 funcionou
// perfeitamente e escondeu que o dado nunca existiu.
//
// Falha fechada em ambas as pontas: erro de leitura não aplica nada (aplicar
// uma lista vazia fecharia toda célula aberta do tenant), e erro de escrita
// sobe para o scheduler, que grava sync_state.last_error e faz a falha
// aparecer na tela /integracoes. Nenhum dos dois inventa um cursor de sucesso.
func NewICMSMatrixJob(
	tenantID string,
	resolver icmsMatrixResolver,
	applier icmsMatrixApplier,
	now func() time.Time,
) syncapp.JobFunc {
	if now == nil {
		now = time.Now
	}
	return func(ctx context.Context, cursor json.RawMessage) (json.RawMessage, error) {
		cells, err := resolver.ResolveCells(ctx)
		if err != nil {
			return cursor, fmt.Errorf("icms matrix sync: resolver células: %w", err)
		}
		written, err := applier.ApplyCells(ctx, tenantID, cells)
		if err != nil {
			return cursor, fmt.Errorf("icms matrix sync: aplicar células: %w", err)
		}
		ambiguos := 0
		for _, c := range cells {
			if c.Ambiguo {
				ambiguos++
			}
		}
		next, err := json.Marshal(ICMSMatrixCursor{
			Cells:       len(cells),
			Ambiguos:    ambiguos,
			Written:     written,
			CompletedAt: now().UTC(),
		})
		if err != nil {
			// A escrita já aconteceu. Reportar o ciclo como falho diria que a
			// matriz não chegou, o que é falso — mantém o cursor anterior.
			return cursor, nil
		}
		return next, nil
	}
}
