import { createStore } from "../AuraJS/index.js";
import { closeWs } from "../game/ws.js";

export const initialState = {
  user: null,
  wsId: null,
  chat: {
    chat: [],
    chatOpen: false,
  },
  lobby: {
    players: 0,
    countdown: null,
    waitTime: null,
  },
  game: {
    started: false,
    width: "992",
    height: "415",
    announcement: {
      id: null,
      text: null,
      timeout: null,
    },
    players: {},
    bombs: {},
    explosions: {},
    walls: {},
    powerups: {},
  },
};

export function reducer(state = initialState, action) {
  switch (action.type) {
    case "REGISTER":
      return {
        ...state,
        user: action.user,
      }

    case "ASSIGN_ID":
      const mappedPlayers = Object.entries(action.payload.players).reduce(
        (acc, [id, name]) => {
          acc[id] = {
            name,
            lives: 3,
          };
          return acc;
        },
        {}
      );

      if (action.payload.position) {
        mappedPlayers[action.payload.position] = {
          name: action.payload.name,
          lives: 3,
        };
      }

      return {
        ...state,
        user: state.user || action.payload.name,
        wsId: action.payload.id,
        lobby: {
          ...state.lobby,
          players: action.payload.playersCount,
        },
        game: {
          ...state.game,
          players: {
            ...state.game.players,
            ...mappedPlayers,
          },
        },
      };

    case "PLAYER_JOINED":
      return {
        ...state,
        lobby: {
          ...state.lobby,
          players: action.payload.playersCount,
        },
        game: {
          ...state.game,
          players: {
            ...state.game.players,
            [action.payload.position]: { name: action.payload.name, lives: 3 },
          }
        }
      }

    case "CHAT_MESSAGE":
      return {
        ...state,
        chat: {
          ...state.chat,
          chat: [...state.chat.chat, action.message],
        },
      };

    case "PRE_GAME_TICK":
      return {
        ...state,
        lobby: {
          ...state.lobby,
          countdown: action.payload.countdown,
        },
      };

    case "LOBBY_TICK":
      return {
        ...state,
        lobby: {
          ...state.lobby,
          waitTime: action.payload.countdown,
        },
      };

    case "START_GAME":
      return {
        ...state,
        game: {
          ...state.game,
          started: action.payload,
        }
      }

    case "TOGGLE_CHAT":
      return {
        ...state,
        chat: { ...state.chat, chatOpen: action.open },
      };

    case "REMOVE_ENTITIES": {
      const { entityKey, ids } = action; // entityKey: "walls" | "players"

      const remaining = { ...state.game[entityKey] };

      for (const id of ids) {
        delete remaining[id];
      }

      return {
        ...state,
        game: {
          ...state.game,
          [entityKey]: remaining,
        },
      };
    }

    case "SET_EXPLOSION_LEVEL":
      return {
        ...state,
        game: {
          ...state.game,
          explosions: {
            ...state.game.explosions,
            [action.id]: {
              ...state.game.explosions[action.id],
              level: action.level
            }
          }
        }
      }

    case "CLEAR_ANNOUNCEMENT":
      if (state.game.announcement.timeout) {
        clearTimeout(state.game.announcement.timeout);
      }

      return {
        ...state,
        game: {
          ...state.game,
          announcement: {
            ...state.game.announcement,
            id: null,
            text: null,
            timeout: null,
          }
        }
      };

    case "STORE_ANNOUNCE_TIMEOUT":
      if (state.game.announcement.timeout) {
        clearTimeout(state.game.announcement.timeout);
      }

      return {
        ...state,
        game: {
          ...state.game,
          announcement: {
            ...state.game.announcement,
            timeout: action?.timeout ? action.timeout : null,
          }
        }
      };

    case "GAME_STATE": {
      const g = action.payload;
      // console.log(action)
      return {
        ...state,
        game: {
          ...state.game,
          started: action.started,
          width: g.match_state?.arena_width ? g.match_state.arena_width * 32 : state.game.width,
          height: g.match_state?.arena_height ? g.match_state.arena_height * 32 : state.game.height,
          announcement: g.match_state ? {
            ...state.game.announcement,
            id: crypto.randomUUID(),
            text: g.match_state.announcement
          } : state.game.announcement,
          bombs: g.bombs ? { ...state.game.bombs, ...g.bombs } : state.game.bombs,
          explosions: g.explosions ? { ...state.game.explosions, ...g.explosions } : state.game.explosions,
          walls: g.walls
            ? Object.entries(g.walls).reduce((acc, [id, val]) => {
              acc[id] = val === null ? null : { ...state.game.walls[id], ...val };
              return acc;
            }, { ...state.game.walls })
            : state.game.walls,
          players: g.players
            ? Object.entries(g.players).reduce((acc, [id, val]) => {
              acc[id] = val === null ? null : { ...state.game.players[id], ...val };
              return acc;
            }, { ...state.game.players })
            : state.game.players,
          powerups: g.powerups ? { ...state.game.powerups, ...g.powerups } : state.game.powerups,
        },
      };
    }

    case "RESET":
      return {
        user: null,
        wsId: null,
        chat: {
          chat: [],
          chatOpen: false,
        },
        lobby: {
          players: 0,
          countdown: null,
        },
        game: {
          started: false,
          width: "992",
          height: "415",
          announcement: {
            id: null,
            text: null,
            timeout: null,
          },
          players: {},
          bombs: {},
          explosions: {},
          walls: {},
          powerups: {},
        },
      };

    default:
      return state;
  }
}

export const store = createStore(reducer, initialState);
