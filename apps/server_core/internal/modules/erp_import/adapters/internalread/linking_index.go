package internalread

import (
	"context"
	"strconv"
	"strings"
	"time"

	erpdomain "marketplace-central/apps/server_core/internal/modules/erp_import/domain"
	readdomain "marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	readports "marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

// defaultLinkingIndexTTL bounds how long one loaded mirror snapshot serves
// linking lookups. Candidate generation asks the matcher about every listing in
// an account — thousands of calls in one run — and each call used to re-read the
// whole mirror, which made a full-catalog run impossible: the request died on
// the server's write timeout before the matcher had seen a tenth of the
// listings. The mirror only changes when a sync writes it, so re-reading it per
// listing bought nothing.
const defaultLinkingIndexTTL = time.Minute

// linkingIndex is one loaded mirror snapshot plus the lookup maps the linking
// anchors need. It is a read-through cache of rows the sync produced, never a
// separate source of truth: every candidate is still built from the row itself.
type linkingIndex struct {
	source     erpdomain.ImportSource
	builtAt    time.Time
	rows       []erpdomain.MirrorProduct
	lowerDesc  []string
	collisions map[string]int
	byCode     map[string][]int
	byEAN      map[string][]int
}

func (r *Reader) linkingIndex(ctx context.Context, source erpdomain.ImportSource) (*linkingIndex, error) {
	ttl := r.linkingTTL
	if ttl <= 0 {
		ttl = defaultLinkingIndexTTL
	}

	r.linkingMu.Lock()
	defer r.linkingMu.Unlock()

	if cached := r.linkingCache; cached != nil && cached.source == source && r.now().Sub(cached.builtAt) < ttl {
		return cached, nil
	}

	rows, err := r.repo.MirrorRows(ctx, r.tenantID, source)
	if err != nil {
		return nil, readdomain.NewReadError(readdomain.ReadErrorSourceUnavailable, "product mirror is unavailable", err)
	}

	index := &linkingIndex{
		source:     source,
		builtAt:    r.now(),
		rows:       rows,
		lowerDesc:  make([]string, len(rows)),
		collisions: validEANCounts(rows),
		byCode:     make(map[string][]int, len(rows)),
		byEAN:      make(map[string][]int, len(rows)),
	}
	for i, row := range rows {
		if code := strings.TrimSpace(row.CodigoProduto); code != "" {
			index.byCode[code] = append(index.byCode[code], i)
		}
		if row.EAN != nil {
			if ean := strings.TrimSpace(*row.EAN); ean != "" {
				index.byEAN[ean] = append(index.byEAN[ean], i)
			}
		}
		if row.Descricao != nil {
			index.lowerDesc[i] = strings.ToLower(*row.Descricao)
		}
	}
	r.linkingCache = index
	return index, nil
}

// candidateRows narrows the mirror to the rows an anchor could possibly match.
// It never decides a match — matches() still does, on every row returned — so
// the answer is identical to scanning the whole mirror, only cheaper. An input
// with no anchor at all keeps the historical meaning: every row.
func (i *linkingIndex) candidateRows(input readports.FindProductsInput) []erpdomain.MirrorProduct {
	productID, sku, ean, title := input.ProductID, matchableSellerSKU(input), trimmed(input.EAN), trimmed(input.Title)
	if !hasSuppliedAnchor(input) {
		return i.rows
	}

	selected := make(map[int]struct{})
	add := func(indexes []int) {
		for _, index := range indexes {
			selected[index] = struct{}{}
		}
	}
	if productID != nil {
		add(i.byCode[strconv.Itoa(*productID)])
	}
	if sku != nil {
		add(i.byCode[*sku])
	}
	if ean != nil {
		add(i.byEAN[*ean])
	}
	if title != nil {
		lowered := strings.ToLower(*title)
		for index, description := range i.lowerDesc {
			if description != "" && strings.Contains(description, lowered) {
				selected[index] = struct{}{}
			}
		}
	}

	rows := make([]erpdomain.MirrorProduct, 0, len(selected))
	for index := range selected {
		rows = append(rows, i.rows[index])
	}
	return rows
}
