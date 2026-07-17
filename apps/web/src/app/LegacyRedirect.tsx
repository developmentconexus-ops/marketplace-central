import { Navigate, useLocation } from "react-router-dom";

interface LegacyRedirectProps {
  to: string;
}

export function LegacyRedirect({ to }: LegacyRedirectProps) {
  const location = useLocation();

  return <Navigate replace to={{ pathname: to, search: location.search }} />;
}
