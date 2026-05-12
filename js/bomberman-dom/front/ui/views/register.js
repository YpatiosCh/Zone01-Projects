import { h } from "../../AuraJS/index.js";
import { closeWs } from "../../game/ws.js";

// Close Ws, clear state, clear local storage
export default function Register(state, dispatch, { navigate, events }) {
  // console.log("STATE: ", state)
  closeWs();
  dispatch({ type: "RESET" })
  localStorage.removeItem("wsId");
  localStorage.removeItem("wsName");

  // const id = Object.keys(state.game.players).length +1 
  return h("div", { class: "page" },
    h("div", { class: "banner" }),
    h("h1", {}, "Welcome to Bomberman!"),
    h("br", {}),
    h("h1", {}, "Enter Username to play:"),
    h("input", {
      class: "register_input",
      onKeydown: e => {
        if (e.key === "Enter" && e.target.value.trim()) {
          events.emit("music:play", "waiting");
          dispatch({ type: "REGISTER", user: e.target.value.trim() });
          navigate("/lobby");
        }
      },
    }),
    h("br", {}),
    h("br", {}),

    h("h6", {}, "press Enter when ready")
  );
}
