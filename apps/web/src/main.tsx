import "./index.css";
import React from "react";
import ReactDOM from "react-dom/client";
import { App } from "./App";
import { ClientProvider } from "./app/ClientContext";
import { QueryClientProvider } from "@tanstack/react-query";
import { createWebQueryClient } from "@marketplace-central/web-query";

const queryClient = createWebQueryClient();

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <ClientProvider>
      <QueryClientProvider client={queryClient}>
        <App />
      </QueryClientProvider>
    </ClientProvider>
  </React.StrictMode>
);
