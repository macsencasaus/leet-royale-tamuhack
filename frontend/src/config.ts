export const debug = import.meta.env.VITE_LR_DEBUG !== undefined;
export const ws_scheme = import.meta.env.VITE_LR_SECURE !== undefined ? "wss" : "ws";
