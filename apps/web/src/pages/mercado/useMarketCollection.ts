import { useCallback, useRef, useState } from "react";
import type {
  MarketPriceIntelCollectionContagens,
  MarketPriceIntelCollectionResponse,
} from "@marketplace-central/sdk-runtime";
import type { MercadoMarketClient } from "./marketClient";

/**
 * The backend collection endpoint is single-product and synchronous by ratified
 * decision (D-F4-p: no batch POST), and every run hits the live ML search. The
 * radar still needs a page-level "Atualizar agora", so the fan-out lives here:
 * bounded concurrency over the codprods on screen, with honest progress and an
 * honest per-cause summary at the end.
 */
const CONCURRENCY = 4;

/** Ceiling on one run, so a full catalog cannot turn the button into an hour-long job. */
export const COLLECTION_CEILING = 200;

export interface MarketCollectionSummary {
  /** Codprods actually collected in this run (≤ COLLECTION_CEILING). */
  attempted: number;
  /** Codprods left out because the run hit the ceiling — surfaced, never silent. */
  skipped: number;
  /** Per-cause totals summed across every single-product response. */
  contagens: MarketPriceIntelCollectionContagens;
  /** Codprods whose collection call itself failed (transport / 5xx / 409). */
  failed: number;
}

export interface MarketCollectionState {
  running: boolean;
  /** Codprods finished so far in the current run. */
  done: number;
  /** Codprods this run will attempt. */
  total: number;
  summary: MarketCollectionSummary | null;
  /** True when the last run finished with no attempt at all (nothing linked). */
  nothingToCollect: boolean;
}

const EMPTY_CONTAGENS: MarketPriceIntelCollectionContagens = {
  ok: 0,
  no_price_evidence: 0,
  insufficient_market: 0,
  no_candidate: 0,
  sem_custo: 0,
};

function addContagens(
  into: MarketPriceIntelCollectionContagens,
  from: MarketPriceIntelCollectionResponse["contagens"] | undefined,
): MarketPriceIntelCollectionContagens {
  if (!from) return into;
  return {
    ok: into.ok + (from.ok ?? 0),
    no_price_evidence: into.no_price_evidence + (from.no_price_evidence ?? 0),
    insufficient_market: into.insufficient_market + (from.insufficient_market ?? 0),
    no_candidate: into.no_candidate + (from.no_candidate ?? 0),
    sem_custo: into.sem_custo + (from.sem_custo ?? 0),
  };
}

export interface UseMarketCollectionOptions {
  client: MercadoMarketClient;
  /** Called once the run finishes, so the caller can invalidate its queries. */
  onFinished?: () => void;
}

export function useMarketCollection({ client, onFinished }: UseMarketCollectionOptions): {
  state: MarketCollectionState;
  collect: (codprods: string[]) => Promise<void>;
} {
  const [state, setState] = useState<MarketCollectionState>({
    running: false,
    done: 0,
    total: 0,
    summary: null,
    nothingToCollect: false,
  });
  const runningRef = useRef(false);

  const collect = useCallback(
    async (codprods: string[]) => {
      if (runningRef.current) return;
      const unique = Array.from(new Set(codprods.filter((code) => code.trim() !== "")));
      if (unique.length === 0) {
        setState({ running: false, done: 0, total: 0, summary: null, nothingToCollect: true });
        return;
      }
      const targets = unique.slice(0, COLLECTION_CEILING);
      const skipped = unique.length - targets.length;

      runningRef.current = true;
      setState({
        running: true,
        done: 0,
        total: targets.length,
        summary: null,
        nothingToCollect: false,
      });

      let contagens = EMPTY_CONTAGENS;
      let failed = 0;
      let next = 0;
      let done = 0;

      async function worker(): Promise<void> {
        for (;;) {
          const index = next++;
          if (index >= targets.length) return;
          try {
            const response = await client.collectMarketPriceIntel(targets[index]);
            contagens = addContagens(contagens, response.contagens);
          } catch {
            // A single product failing (409 already collecting, provider 5xx,
            // transport) must not abort the sweep — it is counted and reported.
            failed++;
          }
          done++;
          setState((previous) => (previous.running ? { ...previous, done } : previous));
        }
      }

      await Promise.all(Array.from({ length: Math.min(CONCURRENCY, targets.length) }, worker));

      runningRef.current = false;
      setState({
        running: false,
        done: targets.length,
        total: targets.length,
        summary: { attempted: targets.length, skipped, contagens, failed },
        nothingToCollect: false,
      });
      onFinished?.();
    },
    [client, onFinished],
  );

  return { state, collect };
}
