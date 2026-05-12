import { h } from "../../AuraJS/index.js";
import { unbindGameInput } from "./game.js";
import { app, ws, reInitSocket } from "../app.js";
import { closeWs } from "../../game/ws.js";

export let unbindGameOver = null;

// Same key bindings as game
export default function GameOver(state, dispatch, ctx) {
    if (state.chat.chatOpen) dispatch({ type: "TOGGLE_CHAT", open: false })
    if (!unbindGameOver) {
        const onKey = (e) => {
            // Quit
            if (e.key === "q" || e.key === "Q") {
                cleanup();
                reInitSocket();
                dispatch({ type: "START_GAME", payload: false })
                ctx.navigate("/");
                return;
            }

            // Continue to another game
            if (e.key === "c" || e.key === "C") {
                cleanup();
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
        }
        // Clean up listeners
        const cleanup = () => {
            window.removeEventListener("keydown", onKey);
            unbindGameOver = null;
        };

        window.addEventListener("keydown", onKey);

        unbindGameOver = cleanup;
    };

    return h("div", { class: "page" },
        h("div", { class: "banner" },
            h("div", { class: "bomber" }),
            h("div", { class: "man" }),
        ),
        h("div", { class: "lobby-player" },
            h("span", { class: `chat-color1` }, localStorage.getItem("wsName")),
            h("div", { class: `lobby-player-icon player p1 down` }),
        ),
        h("div", { class: "page", style: "margin-top: 100px;" },
            h("h2", {}, `Press:`),
            h("h2", {}, `(q) to restart`),
            h("h2", {}, `(c) to continue as: ${localStorage.getItem("wsName")}`),
        )
    );
}