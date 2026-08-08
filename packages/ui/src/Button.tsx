import type { ButtonHTMLAttributes, PropsWithChildren } from "react";

type Variant = "primary" | "secondary" | "danger";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
  loading?: boolean;
}

const variantClasses: Record<Variant, string> = {
  primary: "bg-accent hover:bg-accent-ink text-white border-transparent",
  secondary: "bg-surface hover:bg-surface-2 text-ink border-border",
  danger: "bg-warn hover:opacity-90 text-white border-transparent",
};

export function Button({
  children,
  type = "button",
  variant = "secondary",
  loading = false,
  disabled,
  className = "",
  ...props
}: PropsWithChildren<ButtonProps>) {
  const isDisabled = disabled || loading;
  return (
    <button
      type={type}
      disabled={isDisabled}
      className={`inline-flex items-center gap-2 px-3 py-2 text-sm font-medium rounded-control border cursor-pointer transition-colors duration-150 disabled:opacity-50 disabled:cursor-not-allowed ${variantClasses[variant]} ${className}`}
      {...props}
    >
      {loading && (
        <svg className="animate-spin h-4 w-4" viewBox="0 0 24 24" fill="none">
          <circle
            className="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            strokeWidth="4"
          />
          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
        </svg>
      )}
      {children}
    </button>
  );
}
