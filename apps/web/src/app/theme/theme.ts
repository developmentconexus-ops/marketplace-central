export const STORAGE_KEY = "marketplace-central-theme";

export type Theme = "light" | "dark";

const isTheme = (value: string | null): value is Theme =>
  value === "light" || value === "dark";

export function readStoredTheme(): Theme {
  try {
    const storedTheme = localStorage.getItem(STORAGE_KEY);
    return isTheme(storedTheme) ? storedTheme : "light";
  } catch {
    return "light";
  }
}

export function applyTheme(theme: Theme): void {
  document.documentElement.setAttribute("data-theme", theme);
}

export function persistTheme(theme: Theme): void {
  try {
    localStorage.setItem(STORAGE_KEY, theme);
  } catch {
    return;
  }
}

export function getInitialTheme(): Theme {
  const rootTheme = document.documentElement.getAttribute("data-theme");
  return isTheme(rootTheme) ? rootTheme : readStoredTheme();
}
