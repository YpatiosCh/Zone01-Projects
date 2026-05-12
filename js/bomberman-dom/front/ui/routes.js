import Register from "./views/register.js";
import Lobby from "./views/lobby.js";
import Game from "./views/game.js";
import GameOver from "./views/gameover.js";

export const routes = [
  { path: "/", handler: Register },
  { path: "/lobby", handler: Lobby },
  { path: "/game", handler: Game },
  { path: "/over", handler: GameOver },
];
