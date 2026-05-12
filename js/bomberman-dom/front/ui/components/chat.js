import { h } from "../../AuraJS/index.js";
import { ws } from "../app.js"

// let listener = null;

// export const chatPages = ["/lobby", "/game"];

// export function bindChat(dispatch, getState) {
//     if (listener) return;
//     listener = e => {
//         if (e.key === "t" && !getState().chat.chatOpen) {
//             e.preventDefault();
//             dispatch({ type: "TOGGLE_CHAT", open: true });
//             setTimeout(() => {
//                 const input = document.querySelector(".game-chat-input");
//                 if (input) input.focus();
//             }, 0);
//         }
//     };
//     window.addEventListener("keydown", listener);
// }

// export function unbindChat() {
//     if (listener) {
//         window.removeEventListener("keydown", listener);
//         listener = null;
//     }
// }

export default function Chat(state, dispatch) {
    const { chat, chatOpen } = state.chat;

    function onKeydown(e) {
        e.stopPropagation();
        if (e.key === "Enter" && e.target.value.trim()) {
            const text = e.target.value.trim();
            e.target.value = "";
            ws.send(JSON.stringify({ type: "SEND_MSG", msg: text }))
            dispatch({ type: "TOGGLE_CHAT", open: false });
        }
        if (e.key === "Escape") {
            dispatch({ type: "TOGGLE_CHAT", open: false });
        }
    }


    const messages = chat.slice(-12);

    const players = state.game.players;
    function getColorClass(name) {
        const id = Object.keys(players).find(id => players[id] && players[id].name === name);
        return id ? `chat-color${id}` : "";
    }

    return h("div", {
        class: "game-chat",
        onFocusout: () => dispatch({ type: "TOGGLE_CHAT", open: false }),
    },
        h("div", { class: "game-chat-messages"},
            ...messages.map(msg =>
                h("div", { class: "game-chat-msg" },
                    h("span", { class: `${getColorClass(msg.user)}` }, msg.user + ": "),
                    msg.text,
                )
            )
        ),
        chatOpen
            ? h("input", {

                class: "game-chat-input",
                placeholder: "Press Enter to send, Esc to close",
                onKeydown: onKeydown,
                autofocus: true,
            })
            : h("div", { class: "game-chat-hint" }, "Press T to chat"),
    );
}
