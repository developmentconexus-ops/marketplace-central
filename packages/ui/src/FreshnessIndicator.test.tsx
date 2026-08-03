import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { FreshnessIndicator } from "./FreshnessIndicator";

describe("FreshnessIndicator", () => {
  it("mostra a idade e guarda o instante absoluto no title", () => {
    const iso = new Date(Date.now() - 3 * 60 * 60 * 1000).toISOString();
    render(<FreshnessIndicator asOf={iso} />);

    const el = screen.getByLabelText("Atualização dos dados");
    expect(el).toHaveTextContent("há 3 h");
    // A idade relativa é o que o operador lê de relance; o instante exato tem
    // que continuar alcançável para ele conseguir cruzar com um log.
    expect(el).toHaveAttribute("title", expect.stringContaining(":"));
  });

  it("é honesto quando não há instante", () => {
    render(<FreshnessIndicator asOf={null} />);
    const el = screen.getByLabelText("Atualização dos dados");
    expect(el).toHaveTextContent("idade desconhecida");
    // Sem instante não há title: um tooltip vazio ou "Invalid Date" mentiria.
    expect(el).not.toHaveAttribute("title");
  });
});
