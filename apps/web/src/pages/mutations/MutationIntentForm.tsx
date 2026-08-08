import type { FormEvent } from "react";
import type { MutationType } from "@marketplace-central/sdk-runtime";

export interface MutationIntentFormProps {
  type: MutationType;
  disabled?: boolean;
  onSubmit: (intent: Record<string, unknown>) => void;
}

const inputClassName =
  "rounded-lg border border-slate-300 bg-white px-3 py-2 text-slate-900 outline-none focus:border-blue-600 focus:ring-2 focus:ring-blue-100";
const labelClassName = "flex flex-col gap-1 text-sm font-medium text-slate-700";

function field(form: FormData, name: string): string {
  return String(form.get(name) ?? "").trim();
}

function intentFromForm(type: MutationType, form: FormData): Record<string, unknown> {
  switch (type) {
    case "price_update":
      return { new_price: { amount: field(form, "new_price"), currency: "BRL" } };
    case "stock_correct":
      return { publish_quantity: Number(field(form, "publish_quantity")) };
    case "link_apply": {
      const action = field(form, "link_action");
      if (action === "approve_candidate")
        return { action, candidate_id: field(form, "candidate_id") };
      if (action === "manual_resolve") return { action, product_id: field(form, "product_id") };
      return { action };
    }
    case "listing_edit":
      return {
        attributes: [
          { id: field(form, "attribute_id"), value_name: field(form, "attribute_value") },
        ],
      };
    case "listing_pause":
    case "listing_resync":
      return {};
    default:
      return {};
  }
}

export function MutationIntentForm({ type, disabled = false, onSubmit }: MutationIntentFormProps) {
  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    onSubmit(intentFromForm(type, new FormData(event.currentTarget)));
  };

  return (
    <form id="mutation-intent-form" className="space-y-4" onSubmit={handleSubmit}>
      {type === "price_update" ? (
        <label className={labelClassName}>
          Novo preço
          <input
            className={inputClassName}
            name="new_price"
            inputMode="decimal"
            pattern="\d+(?:[.,]\d{1,2})?"
            required
            disabled={disabled}
            placeholder="49.90"
          />
        </label>
      ) : null}
      {type === "stock_correct" ? (
        <label className={labelClassName}>
          Quantidade a publicar
          <input
            className={inputClassName}
            name="publish_quantity"
            type="number"
            min="0"
            step="1"
            required
            disabled={disabled}
          />
        </label>
      ) : null}
      {type === "link_apply" ? <LinkIntentFields disabled={disabled} /> : null}
      {type === "listing_edit" ? (
        <div className="grid gap-3 sm:grid-cols-2">
          <label className={labelClassName}>
            ID do atributo
            <input className={inputClassName} name="attribute_id" required disabled={disabled} />
          </label>
          <label className={labelClassName}>
            Valor do atributo
            <input className={inputClassName} name="attribute_value" required disabled={disabled} />
          </label>
        </div>
      ) : null}
      {type === "listing_pause" ? (
        <p className="text-sm text-slate-600">
          Os anúncios selecionados serão pausados após sua confirmação.
        </p>
      ) : null}
      {type === "listing_resync" ? (
        <p className="text-sm text-slate-600">
          Os dados dos anúncios selecionados serão buscados novamente.
        </p>
      ) : null}
    </form>
  );
}

function LinkIntentFields({ disabled }: { disabled: boolean }) {
  return (
    <div className="grid gap-3 sm:grid-cols-2">
      <label className={labelClassName}>
        Ação de vínculo
        <select className={inputClassName} name="link_action" required disabled={disabled}>
          <option value="approve_candidate">Aprovar candidato</option>
          <option value="manual_resolve">Resolver manualmente</option>
          <option value="reject_listing">Rejeitar anúncio</option>
        </select>
      </label>
      <label className={labelClassName}>
        ID do candidato
        <input className={inputClassName} name="candidate_id" disabled={disabled} />
      </label>
      <label className={labelClassName}>
        ID do produto
        <input className={inputClassName} name="product_id" disabled={disabled} />
      </label>
    </div>
  );
}
