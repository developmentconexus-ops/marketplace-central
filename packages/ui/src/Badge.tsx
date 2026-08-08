type Status = "pending" | "in_progress" | "succeeded" | "failed" | "completed";

const config: Record<Status, { label: string; classes: string }> = {
  pending: { label: "Pending", classes: "bg-surface-2 text-muted" },
  in_progress: { label: "In Progress", classes: "bg-info-soft text-info" },
  succeeded: { label: "Succeeded", classes: "bg-accent-soft text-accent-ink" },
  failed: { label: "Failed", classes: "bg-warn-soft text-warn" },
  completed: { label: "Completed", classes: "bg-accent-soft text-accent-ink" },
};

interface BadgeProps {
  status: Status;
  className?: string;
}

export function Badge({ status, className = "" }: BadgeProps) {
  const { label, classes } = config[status];
  return (
    <span
      className={`inline-flex items-center px-2 py-0.5 rounded-pill text-xs font-medium ${classes} ${className}`}
    >
      {label}
    </span>
  );
}
