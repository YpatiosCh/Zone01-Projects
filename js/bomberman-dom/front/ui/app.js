import { createApp } from "../AuraJS/index.js";
import * as store from "./store.js";
import { routes } from "./routes.js";
import { connectWS, shouldReconnect } from "../game/ws.js";
import { soundPlugin } from "../game/plugins/audio/audioPlayer.js";
import { unbindLobbyInput } from "./views/lobby.js";
import { unbindGameInput } from "./views/game.js";
import { unbindGameOver } from "./views/gameover.js";
import { startCountdown } from "../ui/views/lobby.js";

export let ws;

export let initSocket = true;
export function reInitSocket() { initSocket = true; }

function view(state, dispatch, ctx) {
  if (ctx.path !== "/game" && unbindGameInput) unbindGameInput();
  if (ctx.path !== "/lobby" && unbindLobbyInput) unbindLobbyInput();
  if (ctx.path !== "/over" && unbindGameOver) {unbindGameOver()};
  if (state.game.started && ctx.path !== "/game") ctx.navigate("/game");

  if (initSocket && (ctx.path === "/game" || ctx.path === "/lobby")) {
    initSocket = false;
    ws = connectWS(state, dispatch, ctx, localStorage.getItem("wsId"));
    console.log("initial connection to socket", localStorage.getItem("wsId"))
  }
  
  return ctx.handler(state, dispatch, ctx);
}

export const app = createApp({
  root: "#app",
  state: store.initialState,
  reducer: store.reducer,
  routes,
  view,
  plugins: [soundPlugin],
});
