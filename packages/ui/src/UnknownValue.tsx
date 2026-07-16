interface UnknownValueProps {
  hint?: string;
}

export function UnknownValue({ hint }: UnknownValueProps) {
  return <span title={hint || undefined}>—</span>;
}
