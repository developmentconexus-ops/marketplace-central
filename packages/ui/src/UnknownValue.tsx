interface UnknownValueProps {
  hint?: string;
}

export function UnknownValue({ hint }: UnknownValueProps) {
  return <span className="text-faint" title={hint || undefined}>—</span>;
}
