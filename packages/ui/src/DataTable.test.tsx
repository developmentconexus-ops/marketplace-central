import { describe, expect, it, vi } from "vitest";
import { fireEvent, render, screen, within } from "@testing-library/react";
import { DataTable, type DataTableColumn } from "./DataTable";

type Item = {
  id: string;
  name: string;
  amount: string;
};

const rows: Item[] = [
  { id: "item-a", name: "Alpha", amount: "—" },
  { id: "item-b", name: "Beta", amount: "R$ 12" },
];

const columns: DataTableColumn<Item>[] = [
  { key: "name", header: "Name", render: (row) => row.name, sortable: true },
  { key: "amount", header: "Amount", render: (row) => row.amount, align: "right" },
];

describe("DataTable", () => {
  it("renders column headers and cells from the supplied renderers", () => {
    render(<DataTable columns={columns} rows={rows} rowKey={(row) => row.id} />);

    expect(screen.getByRole("columnheader", { name: "Name" })).toBeInTheDocument();
    expect(screen.getByRole("columnheader", { name: "Amount" })).toBeInTheDocument();
    expect(screen.getByRole("cell", { name: "Alpha" })).toBeInTheDocument();
    expect(screen.getByRole("cell", { name: "Beta" })).toBeInTheDocument();
    expect(screen.getByRole("cell", { name: "—" })).toBeInTheDocument();
    expect(screen.getByRole("cell", { name: "R$ 12" })).toBeInTheDocument();
  });

  it("calls onRowClick for a row but not when its checkbox is clicked", () => {
    const onRowClick = vi.fn();
    const onSelectionChange = vi.fn();

    render(
      <DataTable
        columns={columns}
        rows={rows}
        rowKey={(row) => row.id}
        onRowClick={onRowClick}
        onSelectionChange={onSelectionChange}
      />,
    );

    const row = screen.getByRole("row", { name: /Alpha/ });
    fireEvent.click(row);
    expect(onRowClick).toHaveBeenCalledWith(rows[0]);

    fireEvent.click(within(row).getByRole("checkbox"));
    expect(onRowClick).toHaveBeenCalledTimes(1);
    expect(onSelectionChange).toHaveBeenCalledWith(new Set(["item-a"]));
  });

  it("supports controlled select-all and row selection", () => {
    const onSelectionChange = vi.fn();
    const { rerender } = render(
      <DataTable
        columns={columns}
        rows={rows}
        rowKey={(row) => row.id}
        selectedKeys={new Set()}
        onSelectionChange={onSelectionChange}
      />,
    );

    const checkboxes = screen.getAllByRole("checkbox");
    expect(checkboxes).toHaveLength(3);
    fireEvent.click(checkboxes[0]);
    expect(onSelectionChange).toHaveBeenCalledWith(new Set(["item-a", "item-b"]));

    rerender(
      <DataTable
        columns={columns}
        rows={rows}
        rowKey={(row) => row.id}
        selectedKeys={new Set(["item-a"])}
        onSelectionChange={onSelectionChange}
      />,
    );
    expect(checkboxes[0].indeterminate).toBe(true);
    expect(screen.getByRole("row", { name: /Alpha/ })).toHaveClass("bg-accent-soft");

    fireEvent.click(screen.getAllByRole("checkbox")[1]);
    expect(onSelectionChange).toHaveBeenLastCalledWith(new Set());

    rerender(
      <DataTable
        columns={columns}
        rows={rows}
        rowKey={(row) => row.id}
        selectedKeys={new Set(["item-a"])}
        onSelectionChange={onSelectionChange}
      />,
    );
    fireEvent.click(screen.getAllByRole("checkbox")[2]);
    expect(onSelectionChange).toHaveBeenLastCalledWith(new Set(["item-a", "item-b"]));
  });

  it("does not render selection controls without onSelectionChange", () => {
    render(<DataTable columns={columns} rows={rows} rowKey={(row) => row.id} />);

    expect(screen.queryByRole("checkbox")).not.toBeInTheDocument();
  });

  it("calls the controlled sort callback and shows the active direction", () => {
    const onSortChange = vi.fn();
    const { rerender } = render(
      <DataTable
        columns={columns}
        rows={rows}
        rowKey={(row) => row.id}
        sortKey="name"
        sortDir="asc"
        onSortChange={onSortChange}
      />,
    );

    const sortButton = screen.getByRole("button", { name: /Name/ });
    expect(sortButton).toHaveTextContent("▲");
    fireEvent.click(sortButton);
    expect(onSortChange).toHaveBeenCalledWith("name");

    rerender(
      <DataTable
        columns={columns}
        rows={rows}
        rowKey={(row) => row.id}
        sortKey="name"
        sortDir="desc"
        onSortChange={onSortChange}
      />,
    );
    expect(screen.getByRole("button", { name: /Name/ })).toHaveTextContent("▼");
    expect(screen.getByRole("columnheader", { name: "Amount" }).querySelector("button")).toBeNull();
  });

  it("shows loading and empty states", () => {
    const { rerender } = render(
      <DataTable columns={columns} rows={rows} rowKey={(row) => row.id} loading />,
    );
    expect(screen.getByRole("status")).toBeInTheDocument();
    expect(screen.queryByRole("cell", { name: "Alpha" })).not.toBeInTheDocument();

    rerender(
      <DataTable columns={columns} rows={[]} rowKey={(row) => row.id} emptyState="Nothing here" />,
    );
    expect(screen.getByText("Nothing here")).toBeInTheDocument();

    rerender(<DataTable columns={columns} rows={[]} rowKey={(row) => row.id} />);
    expect(screen.getByText("Nenhum item.")).toBeInTheDocument();
  });

  it("uses the semantic shell tokens and stays value-agnostic", () => {
    render(
      <DataTable
        columns={columns}
        rows={rows.slice(0, 1)}
        rowKey={(row) => row.id}
        selectedKeys={new Set(["item-a"])}
        onSelectionChange={vi.fn()}
      />,
    );

    const table = screen.getByRole("table");
    expect(table.parentElement).toHaveClass("border-border", "rounded-card");
    expect(table.querySelector("thead")).toHaveClass("bg-surface-2");
    expect(screen.getByRole("row", { name: /Alpha/ })).toHaveClass("bg-accent-soft");
    expect(screen.queryByRole("cell", { name: "R$ 12" })).not.toBeInTheDocument();
    expect(screen.getByRole("cell", { name: "—" })).toHaveClass("text-right");
    expect(screen.queryByText("0")).not.toBeInTheDocument();
    expect(screen.queryByText("R$ 0")).not.toBeInTheDocument();
  });

  it("sticks the header when requested", () => {
    render(<DataTable columns={columns} rows={rows} rowKey={(row) => row.id} stickyHeader />);

    expect(screen.getAllByRole("row")[0].parentElement).toHaveClass("sticky", "top-0");
  });
});
