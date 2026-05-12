import { h } from "../AuraJS/index.js";

const tile = (cls, entity, id, dispatch, entityKey) => {
  // console.log("TILE:", cls, entity, id, dispatch, entityKey)

  if (entity?.remove) {
    console.log("removing entity", entity, entityKey)
    dispatch({ type: "REMOVE_ENTITIES", entityKey: entityKey, ids: [id] })
  }

  let empty = false;
  const { x, y } = entity;

  if (Object.entries(entity).length === 0) {
    empty = true
    console.log("no coords", entity, id) // debug
  }

  return h("div", {
    key: id,
    class: empty ? "hidden" : `${cls}`,
    style: `left:${convertCoordsToPix(x)}px; top:${convertCoordsToPix(y)}px`,
  });
};


// think of it as a lock
// we will use it below to lock an action that is already in process
const breakingWalls = new Set();
const dyingPlayers = new Set();
const hitPlayers = new Set();
const playerLives = {};
const explosions = new Set();
const announcements = new Set();

export function renderAnnouncement(announcement, dispatch) {
  if (!announcement.text) return;

  if (!announcements.has(announcement.id)) {
    announcements.add(announcement.id)
    let timeout = setTimeout(() => {
      announcements.delete(announcement.id)
      dispatch({ type: "CLEAR_ANNOUNCEMENT" })
    }, 10000)
    dispatch({ type: "STORE_ANNOUNCE_TIMEOUT", timeout: timeout })
  }

  return h("div", {class: "announcement"}, announcement.text);
}

export function renderEntities(game, dispatch, events) {
  const entries = obj => Object.entries(obj).filter(([, v]) => v && Object.keys(v).length);

  return [
    ...entries(game.bombs).map(([id, b]) => {
      return tile("bomb", b, id, dispatch, "bombs")
    }),

    ...entries(game.walls).map(([id, w]) => {
      if (w.destroy === true && !breakingWalls.has(id)) {
        breakingWalls.add(id);
        events.emit("wall:break");
        setTimeout(() => {
          dispatch({ type: "REMOVE_ENTITIES", entityKey: "walls", ids: [id] });
          breakingWalls.delete(id);
        }, 500);
      }

      if (w.permanent) {
        return tile("stone", w, id, dispatch, "stones")
      }
      return tile(w.destroy ? "wall_destroy" : "wall", w, id, dispatch, "walls")
    }),

    ...entries(game.players).map(([id, p]) => {
      if (p.destroy === true && !dyingPlayers.has(id)) {
        dyingPlayers.add(id);
        events.emit("player:destroy");
        setTimeout(() => {
          dispatch({ type: "REMOVE_ENTITIES", entityKey: "players", ids: [id] });
          dyingPlayers.delete(id);
        }, 1500);
      }

      // Detect life loss — flash red for any player
      if (p.lives != null && playerLives[id] != null && p.lives < playerLives[id] && !hitPlayers.has(id)) {
        hitPlayers.add(id);
        setTimeout(() => hitPlayers.delete(id), 500);
      }
      if (p.lives != null) playerLives[id] = p.lives;

      const hit = hitPlayers.has(id) ? "hit" : "";

      if (p.destroy) return tile(`player_destroy p${id}`, p, id, dispatch, "players");

      p.direction = p.direction ? p.direction : "";
      p.animate = p.animate ? p.animate : "";
      return tile(`player p${id} ${p.direction} ${p.animate} ${hit}`, p, id, dispatch, "players")
    }),

    ...entries(game.explosions).map(([id, e]) => {
      if (e.remove === true) {
        events.emit(`stop:explosion:${id}`);
        explosions.delete(id)
      } else {
        animateExplosion(id, dispatch, events)
      }
      return explosion(e, id, dispatch, "explosions")
    }),

    ...entries(game.powerups).map(([id, w]) => {
      return tile(`powerups ${w.name}`, w, id, dispatch, "powerups")
    }),
  ];
}


function explosion(e, id, dispatch, entityKey) {
  if (e?.remove) {
    dispatch({ type: "REMOVE_ENTITIES", entityKey: entityKey, ids: [id] });
    explosions.delete(id);
  }

  const parts = [];

  parts.push(h("div", {
    key: id + "-center",
    class: `exp level-${e.level} exp-center`,
    style: `left:${convertCoordsToPix(e.x)}px; top:${convertCoordsToPix(e.y)}px`
  }));

  parts.push(...expandExp(e, id));
  return parts;
}

function animateExplosion(id, dispatch, events) {
  if (explosions.has(id)) return;
  explosions.add(id)

  let level = 1
  let op = 1;
  events.emit("explosion:start");

  const interval = setInterval(() => {
    level += op
    if (level < 1) {
      clearInterval(interval)
      dispatch({ type: "REMOVE_ENTITIES", entityKey: "explosions", ids: [id] });
      explosions.delete(id)
    } else if (level === 4) {
      op = -1
      dispatch({ type: "SET_EXPLOSION_LEVEL", id, level })
    } else {
      dispatch({ type: "SET_EXPLOSION_LEVEL", id, level })
    }
  }, 60)

  events.on(`stop:explosion:${id}`, () => clearInterval(interval))
}


function expandExp(e, id) {
  const children = [];

  let offset = convertCoordsToPix(e.y) - 32;
  if (!e.level) e.level = 1;

  if (e.upLen > 1) {
    for (let i = 1; i < e.upLen; i++) {
      children.push({
        key: id + "-mid-up" + i,
        class: `exp level-${e.level} exp-mid-vertical`,
        style: `left:${convertCoordsToPix(e.x)}px; top:${offset}px`
      })
      offset -= 32
    }
  }

  if (e.upLen >= 1) {
    children.push({
      key: id + "-top-edge",
      class: `exp level-${e.level} exp-top`,
      style: `left:${convertCoordsToPix(e.x)}px; top:${offset}px`,
    })
  }

  offset = convertCoordsToPix(e.y) + 32;
  if (e.downLen > 1) {
    for (let i = 1; i < e.downLen; i++) {
      children.push({
        key: id + "-mid-down" + i,
        class: `exp level-${e.level} exp-mid-vertical`,
        style: `left:${convertCoordsToPix(e.x)}px; top:${offset}px`
      })
      offset += 32
    }
  }

  if (e.downLen >= 1) {
    children.push({
      key: id + "-bottom-edge",
      class: `exp level-${e.level} exp-bottom`,
      style: `left:${convertCoordsToPix(e.x)}px; top:${offset}px`,
    })
  }

  offset = convertCoordsToPix(e.x) - 32;

  if (e.leftLen > 1) {
    for (let i = 1; i < e.leftLen; i++) {
      children.push({
        key: id + "-mid-left" + i,
        class: `exp level-${e.level} exp-mid-horizontal`,
        style: `left:${offset}px; top:${convertCoordsToPix(e.y)}px`
      })
      offset -= 32
    }
  }

  if (e.leftLen >= 1) {
    children.push({
      key: id + "-left-edge",
      class: `exp level-${e.level} exp-left`,
      style: `left:${offset}px; top:${convertCoordsToPix(e.y)}px`,
    })
  }

  offset = convertCoordsToPix(e.x) + 32;
  if (e.rightLen > 1) {
    for (let i = 1; i < e.rightLen; i++) {
      children.push({
        key: id + "-mid-right" + i,
        class: `exp level-${e.level} exp-mid-horizontal`,
        style: `left:${offset}px; top:${convertCoordsToPix(e.y)}px`,
      })
      offset += 32
    }
  }

  if (e.rightLen >= 1) {
    children.push({
      key: id + "-right-edge",
      class: `exp level-${e.level} exp-right`,
      style: `left:${offset}px; top:${convertCoordsToPix(e.y)}px`,
    })
  }

  const divs = [];
  children.map((c) => {
    divs.push(h("div", c))
  })
  return divs
}

function convertCoordsToPix(coord) {
  return (coord * 32)
}


