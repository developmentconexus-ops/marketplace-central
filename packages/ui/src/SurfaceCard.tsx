import type { PropsWithChildren } from "react";

interface SurfaceCardProps {
  className?: string;
}

export function SurfaceCard({ children, className = "" }: PropsWithChildren<SurfaceCardProps>) {
  return (
    <section className={`bg-surface border border-border rounded-card p-6 ${className}`}>
      {children}
    </section>
  );
}
