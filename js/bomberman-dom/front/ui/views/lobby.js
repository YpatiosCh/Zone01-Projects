import { h } from "../../AuraJS/index.js";
import Chat from "../components/chat.js";
import { app } from "../app.js";

let countdownInterval = null;

export function startCountdown(events) {
  if (countdownInterval) return;
  
  console.log("Starting lobby count down", )
  
  countdownInterval = setInterval(() => {
    const cd = app.store.getState().lobby.countdown;
    if (cd === 4) {
      // events.emit("music:stop", "waiting");
      events.emit("game:countdown");
    }
    if (cd != null && cd === 1) {
      console.log("Stopping lobby count down")
      stopCountdown();
    }
  }, 1000);
}

export function stopCountdown() {
  if (countdownInterval) {
    clearInterval(countdownInterval);
    countdownInterval = null;
  }
}


export let unbindLobbyInput = null;


export default function Lobby(state, dispatch, ctx) {
  // console.log("State:", state)
  
  if (!unbindLobbyInput) {
    startCountdown(ctx.events);

    const onKey = (e) => {
      if (app.store.getState().chat.chatOpen) return;
      if (e.key === "t" || e.key === "T") {
        e.preventDefault();
        dispatch({ type: "TOGGLE_CHAT", open: true });
        setTimeout(() => {
          document.querySelector(".game-chat-input")?.focus();
        }, 0);
        return;
      }
    }
    window.addEventListener("keydown", onKey);
    const cleanup = () => {
      window.removeEventListener("keydown", onKey);
      unbindLobbyInput = null;
    }
    unbindLobbyInput = cleanup;
  }

  const players = Object.entries(state.game.players).filter(([, p]) => p);
  return h("div", { class: "page lobby" },
    h("h2", {}, `Welcome ${state.user}`),
    h("p", {}, `Players online: ${state.lobby.players}`),
    state.lobby.countdown !== null
      ? h("p", { class: "countdown" }, `Game starts in: ${state.lobby.countdown}s`)
      : h("p", {}, `Waiting for other players to join${state.lobby.waitTime != null ? ` (${state.lobby.waitTime})` : ""}`),
    h("div", { class: "lobby-players" },
      ...players.map(([id, p]) =>
        h("div", { class: "lobby-player" },
          h("span", { class: `chat-color${id}` }, p.name),
          h("div", { class: `lobby-player-icon player p${id} down` }),
        )
      ),
    ),
    Chat(state, dispatch, ctx),
  );
}
