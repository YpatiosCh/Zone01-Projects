import { h } from "../../AuraJS/index.js";
import { renderEntities, renderAnnouncement } from "../../game/renderer.js";
import { closeWs } from "../../game/ws.js";
import { app, ws, reInitSocket } from "../app.js";
import Chat from "../components/chat.js";

export let unbindGameInput = null;

const gameKeys = {
  ArrowUp: "up",
  ArrowDown: "down",
  ArrowLeft: "left",
  ArrowRight: "right",


  w: "up",
  W: "up",
  a: "left",
  A: "left",
  s: "down",
  S: "down",
  d: "right",
  D: "right",

  " ": "bomb",
  "t": "chat",
  "T": "chat",
  "P": "leave_game", // re-register
  "p": "leave_game",
  "O": "leave_game", // new game with same name
  "o": "leave_game"
};

export default function Game(state, dispatch, ctx) {
  // console.log("GAME UNBIND", unbindGameInput)
  if (!unbindGameInput) {

    const onKey = (e) => {
      if (app.store.getState().chat.chatOpen) return;

      if (!gameKeys[e.key]) return;

      // Quit
      if (e.key === "P" || e.key === "p") {
        cleanup();
        reInitSocket();
        ws.send(JSON.stringify({ type: "LEAVE_GAME" }));
        dispatch({ type: "START_GAME", payload: false })
        ctx.navigate("/");
        return;
      }

      // Continue to another game
      if (e.key === "O" || e.key === "o") {
        cleanup();
        ws.send(JSON.stringify({ type: "LEAVE_GAME" }));
        closeWs();
        localStorage.setItem("wsId", crypto.randomUUID());
        app.store.batch(() => {
          dispatch({ type: "RESET" });
          dispatch({ type: "START_GAME", payload: false })
          dispatch({ type: "REGISTER", user: localStorage.getItem("wsName") });
        })
        reInitSocket();
        ctx.navigate("/lobby");
        return;
      }

      // Chat
      if (e.key === "t") {
        e.preventDefault();
        dispatch({ type: "TOGGLE_CHAT", open: true });
        setTimeout(() => {
          document.querySelector(".game-chat-input")?.focus();
        }, 0);
        return;
      }

      e.preventDefault();
      ws.send(JSON.stringify({
        type: "INPUT",
        key: gameKeys[e.key],
        pressed: true
      }));
    };

    const offKey = (e) => {
      if (app.store.getState().chat.chatOpen) return;
      if (!gameKeys[e.key]) return;

      e.preventDefault();
      ws.send(JSON.stringify({
        type: "INPUT",
        key: gameKeys[e.key],
        pressed: false
      }));
    };

    // Clean up listeners
    const cleanup = () => {
      window.removeEventListener("keydown", onKey);
      window.removeEventListener("keyup", offKey);
      unbindGameInput = null;
    };

    window.addEventListener("keydown", onKey);
    window.addEventListener("keyup", offKey);

    unbindGameInput = cleanup;
  }


  const players = Object.entries(state.game.players).filter(([, p]) => p);

  return h("div", { class: "page" },

    h("div", { class: "banner" },
      h("div", { class: "bomber" }),
      h("div", { class: "man" }),
    ),
    h("div", {},
      h("h2", {class: "game-chat-hint"}, `Press (p) to restart (o) for new game as: ${localStorage.getItem("wsName")}`),
    ),
    h("div", { class: "game-and-players" },
      ...players.map(([id, p]) =>
        h("div", { class: `player-info info-${id}` },
          h("div", { class: "lobby-player" },
            h("span", { class: `chat-color${id}` }, p.name),
            h("div", { class: "icon-wrapper" },
              h("div", { class: `lobby-player-icon player p${id} down` }),
            ),
            h("div", { class: "lives-count" }, `Lives ${state.game.players[id].lives}`)
          ),
        ),
      ),

      h("div", { class: "game", style: `width: ${state.game.width}px; height: ${state.game.height}px;` },
        ...renderEntities(state.game, dispatch, ctx.events),
      ),
    ),
    h("div", { class: "announcement" },
      renderAnnouncement(state.game.announcement, dispatch),
    ),
    Chat(state, dispatch, ctx),
  );
}
