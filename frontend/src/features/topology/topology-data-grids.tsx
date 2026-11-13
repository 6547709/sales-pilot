"use client";

import { useTable, type ColumnDef } from "@tanstack/react-table";
import {
  DataGrid,
  DataGridContainer,
  dataGridFeatures,
  type DataGridFeatures,
} from "@/components/reui/data-grid/data-grid";
import { DataGridTable } from "@/components/reui/data-grid/data-grid-table";
import { DataGridPagination } from "@/components/reui/data-grid/data-grid-pagination";
import { DataGridColumnHeader } from "@/components/reui/data-grid/data-grid-column-header";
import { Badge } from "@/components/reui/badge";
import { Button } from "@/components/ui/button";

export type LayerRow = {
  id: number;
  level: number;
  title: string;
  subtitle: string;
  sort_order: number;
};

export type CategoryRow = {
  id: number;
  slug: string;
  label: string;
  column_type: string;
  layer_id?: number | null;
  is_active: boolean;
};

export type VendorRow = {
  id: number;
  category_id: number;
  category_label: string;
  market: string;
  name: string;
};

type GridColumns<T extends object> = ColumnDef<DataGridFeatures, T>[];

export function LayersGrid({
  layers,
  onDelete,
}: {
  layers: LayerRow[];
  onDelete: (id: number) => void;
}) {
  const columns: GridColumns<LayerRow> = [
    {
      accessorKey: "id",
      header: ({ column }) => <DataGridColumnHeader column={column} title="ID" />,
      size: 64,
    },
    {
      accessorKey: "level",
      header: ({ column }) => <DataGridColumnHeader column={column} title="层级" />,
      size: 72,
    },
    {
      accessorKey: "title",
      header: ({ column }) => <DataGridColumnHeader column={column} title="标题" />,
      cell: ({ row }) => <span className="font-medium">{row.original.title}</span>,
    },
    {
      accessorKey: "subtitle",
      header: ({ column }) => <DataGridColumnHeader column={column} title="副标题" />,
      cell: ({ row }) => (
        <span className="line-clamp-1 text-muted-foreground">
          {row.original.subtitle || "—"}
        </span>
      ),
    },
    {
      accessorKey: "sort_order",
      header: ({ column }) => <DataGridColumnHeader column={column} title="排序" />,
      size: 72,
    },
    {
      id: "actions",
      header: () => <span className="text-right">操作</span>,
      cell: ({ row }) => (
        <div className="flex justify-end">
          <Button
            variant="destructive"
            size="sm"
            onClick={() => onDelete(row.original.id)}
          >
            删除
          </Button>
        </div>
      ),
    },
  ];

  const table = useTable({
    features: dataGridFeatures,
    data: layers,
    columns,
    getRowId: (row) => String(row.id),
    initialState: { pagination: { pageIndex: 0, pageSize: 10 } },
  });

  return (
    <DataGrid
      table={table}
      recordCount={layers.length}
      emptyMessage="暂无分层"
      tableLayout={{ dense: true, headerBackground: true, stripped: true, rowRounded: true }}
    >
      <DataGridContainer>
        <DataGridTable />
      </DataGridContainer>
      <DataGridPagination sizes={[10, 20, 50]} rowsPerPageLabel="每页" />
    </DataGrid>
  );
}

export function CategoriesGrid({
  categories,
  onEdit,
  onDelete,
}: {
  categories: CategoryRow[];
  onEdit: (row: CategoryRow) => void;
  onDelete: (id: number) => void;
}) {
  const columns: GridColumns<CategoryRow> = [
    {
      accessorKey: "id",
      header: ({ column }) => <DataGridColumnHeader column={column} title="ID" />,
      size: 64,
    },
    {
      accessorKey: "label",
      header: ({ column }) => <DataGridColumnHeader column={column} title="名称" />,
      cell: ({ row }) => <span className="font-medium">{row.original.label}</span>,
    },
    {
      accessorKey: "slug",
      header: ({ column }) => <DataGridColumnHeader column={column} title="slug" />,
      cell: ({ row }) => (
        <span className="text-muted-foreground">{row.original.slug}</span>
      ),
    },
    {
      accessorKey: "column_type",
      header: ({ column }) => <DataGridColumnHeader column={column} title="栏目" />,
      size: 96,
    },
    {
      accessorKey: "layer_id",
      header: ({ column }) => <DataGridColumnHeader column={column} title="分层" />,
      size: 80,
      cell: ({ row }) => row.original.layer_id ?? "—",
    },
    {
      id: "status",
      header: "状态",
      size: 88,
      cell: ({ row }) => (
        <Badge variant={row.original.is_active ? "success-light" : "outline"}>
          {row.original.is_active ? "启用" : "停用"}
        </Badge>
      ),
    },
    {
      id: "actions",
      header: () => <span className="text-right">操作</span>,
      cell: ({ row }) => (
        <div className="flex justify-end gap-1">
          <Button variant="secondary" size="sm" onClick={() => onEdit(row.original)}>
            编辑
          </Button>
          <Button variant="destructive" size="sm" onClick={() => onDelete(row.original.id)}>
            删除
          </Button>
        </div>
      ),
    },
  ];

  const table = useTable({
    features: dataGridFeatures,
    data: categories,
    columns,
    getRowId: (row) => String(row.id),
    initialState: { pagination: { pageIndex: 0, pageSize: 10 } },
  });

  return (
    <DataGrid
      table={table}
      recordCount={categories.length}
      emptyMessage="暂无解决方案节点"
      tableLayout={{ dense: true, headerBackground: true, stripped: true, rowRounded: true }}
    >
      <DataGridContainer>
        <DataGridTable />
      </DataGridContainer>
      <DataGridPagination sizes={[10, 20, 50]} rowsPerPageLabel="每页" />
    </DataGrid>
  );
}

export function VendorsGrid({
  rows,
  onDelete,
}: {
  rows: VendorRow[];
  onDelete: (id: number) => void;
}) {
  const columns: GridColumns<VendorRow> = [
    {
      accessorKey: "id",
      header: ({ column }) => <DataGridColumnHeader column={column} title="ID" />,
      size: 64,
    },
    {
      id: "category",
      accessorFn: (row) => row.category_label,
      header: ({ column }) => <DataGridColumnHeader column={column} title="方案" />,
      cell: ({ row }) => (
        <span>
          {row.original.category_label}{" "}
          <span className="text-muted-foreground">(#{row.original.category_id})</span>
        </span>
      ),
    },
    {
      id: "market",
      header: ({ column }) => <DataGridColumnHeader column={column} title="市场" />,
      size: 88,
      cell: ({ row }) => (
        <Badge variant={row.original.market === "foreign" ? "secondary" : "info-light"}>
          {row.original.market === "foreign" ? "国外" : "国内"}
        </Badge>
      ),
    },
    {
      accessorKey: "name",
      header: ({ column }) => <DataGridColumnHeader column={column} title="名称" />,
      cell: ({ row }) => <span className="font-medium">{row.original.name}</span>,
    },
    {
      id: "actions",
      header: () => <span className="text-right">操作</span>,
      cell: ({ row }) => (
        <div className="flex justify-end">
          <Button variant="destructive" size="sm" onClick={() => onDelete(row.original.id)}>
            删除
          </Button>
        </div>
      ),
    },
  ];

  const table = useTable({
    features: dataGridFeatures,
    data: rows,
    columns,
    getRowId: (row) => String(row.id),
    initialState: { pagination: { pageIndex: 0, pageSize: 10 } },
  });

  return (
    <DataGrid
      table={table}
      recordCount={rows.length}
      emptyMessage="暂无数据或拓扑加载失败"
      tableLayout={{ dense: true, headerBackground: true, stripped: true, rowRounded: true }}
    >
      <DataGridContainer>
        <DataGridTable />
      </DataGridContainer>
      <DataGridPagination sizes={[10, 20, 50]} rowsPerPageLabel="每页" />
    </DataGrid>
  );
}
