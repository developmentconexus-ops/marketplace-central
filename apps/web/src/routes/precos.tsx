import { PricingSimulatorPage } from "@marketplace-central/feature-simulator";
import { useClient } from "../app/ClientContext";

export function PrecosRoute() {
  const client = useClient();
  return <PricingSimulatorPage client={client} />;
}
