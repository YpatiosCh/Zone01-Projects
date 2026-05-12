import { ws, initSocket } from "../ui/app.js"

export let shouldReconnect = true;
let reconnectTimeout = null;
let wsConn;

let prevGameState = "hello"

export function connectWS(state, dispatch, ctx, savedId) {
  let myPosition = null;
  let myLives = null;

  let connect = (name, id) => {
    let nameExt = `${encodeURIComponent(name || localStorage.getItem("wsName") || '')}`
    console.log("name", nameExt)
    // let url = `wss://bomberman-latest.onrender.com/ws?name=${nameExt}`;
    let url = `ws://localhost:8080/ws?name=${nameExt}`;
    if (id) url += `&id=${id}`;
    wsConn = new WebSocket(url);
    shouldReconnect = true;

    wsConn.onopen = () => {
      console.log("WS connected");
    };

    wsConn.onmessage = e => {
      const msg = JSON.parse(e.data);

      if (msg.type === "LOBBY_CANCEL") {
        console.log("HERE!")
      }


      if (msg.type !== "FRAME_UPDATE") {
        console.log(`Ws Message: ${msg.type}`, msg);
      }

      if (msg.type === "FRAME_UPDATE" || msg.type === "GAME_STATE" || msg.type === "INITIAL_STATE") {
        if (myPosition && msg.data.players?.[myPosition]) {
          const me = msg.data.players[myPosition];
          // Check if our player picked up a powerup
          if (me.item_picked_up) {
            ctx.events.emit("power:up");
          }
          // Check if our player lost a life
          if (me.lives != null && myLives != null && me.lives < myLives) {
            ctx.events.emit("life:lost");
          }
          if (me.lives != null) myLives = me.lives;
        }

        // Check if new bomb appears
        if (msg.data.bombs != null && Object.values(msg.data.bombs).some(b => b.destroy != null && b.destroy == false)) {
          ctx.events.emit("bomb-placed:start");
        }

        if (msg.type === "INITIAL_STATE") dispatch({ type: "GAME_STATE", payload: msg.data, started: false });
        if (msg.type === "FRAME_UPDATE" || msg.type === "GAME_STATE") {
          const started = !(msg.data.go_to_game != null && msg.data.go_to_game === false);
          dispatch({ type: "GAME_STATE", payload: msg.data, started });
        }
      }

      if (msg.type === "PRE_GAME_TICK" || msg.type === "PRE_GAME") {
        dispatch({ type: "PRE_GAME_TICK", payload: msg.data });
      }

      if (msg.type === "LOBBY_TICK") {
        dispatch({ type: "LOBBY_TICK", payload: msg.data });
      }

      if (msg.type === "PLAYER_JOINED") {
        dispatch({ type: "PLAYER_JOINED", payload: msg.data });
      }

      if (msg.type === "ASSIGN_ID") {
        myPosition = String(msg.data.position);
        localStorage.setItem("wsId", msg.data.id);
        if (!localStorage.getItem("wsName")) {
          localStorage.setItem("wsName", msg.data.name);
        }
        dispatch({ type: "ASSIGN_ID", payload: msg.data });
      }

      if (msg.type === "START_GAME") {
        dispatch({ type: "START_GAME", payload: true })
        ctx.navigate("/game");
      }

      if (msg.type === "GAME_CHAT") {
        dispatch({ type: "CHAT_MESSAGE", message: msg.message });
      }

      if (msg.type === "ANNOUNCE") {
        let timeout = setTimeout(() => {
          dispatch({ type: "REMOVE_ANNOUNCE" })
        }, 10000)
        dispatch({ type: "ANNOUNCE", payload: msg.data, timeout: timeout })
      }

      if (msg.type === "REMOVE_ANNOUNCE") {
        dispatch({ type: msg.type })
      }

      if (msg.type === "GAME_OVER") {
        dispatch({ type: "RESET" })
        closeWs();
        ctx.navigate("/over")
        return;
      }
    };

    wsConn.onclose = () => {
      console.log("WS disconnected");
      if (!shouldReconnect || !initSocket) return null;

      // Try reconnecting after delay
      if (!reconnectTimeout && shouldReconnect && initSocket) {
        reconnectTimeout = setTimeout(() => {
          reconnectTimeout = null;
          console.log("Reconnecting...", initSocket);
          connect(state.user, state.wsId);
        }, 2000);
      }
    };

    wsConn.onerror = (err) => {
      console.log("WS error", err);
      wsConn.close(); // triggers onclose -> reconnect
    };

    return wsConn;
  };

  return connect(state.user, savedId); // initial connection uses user, savedId for reconnect
}

export function closeWs() {
  if (reconnectTimeout) clearTimeout(reconnectTimeout);
  shouldReconnect = false;
  if (ws) ws.close();
  // console.log("closing socket", reconnectTimeout, shouldReconnect, ws)
}

