import type { JSX, ReactNode } from "react";
import { DetailPanel } from "./DetailPanel";

export interface DetailDrawerProps {
  open: boolean;
  onClose: () => void;
  title: string;
  subtitle?: string;
  children: ReactNode;
  actions?: ReactNode;
  closeLabel?: string;
  width?: number;
}

export function DetailDrawer({
  open,
  onClose,
  title,
  subtitle,
  children,
  actions,
  closeLabel,
  width,
}: DetailDrawerProps): JSX.Element | null {
  return (
    <DetailPanel
      open={open}
      onClose={onClose}
      title={title}
      subtitle={subtitle}
      closeLabel={closeLabel ?? "Fechar"}
      width={width ?? 360}
      footer={actions}
    >
      {children}
    </DetailPanel>
  );
}
