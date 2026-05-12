const sounds = {
  "wall:break": new Audio("game/plugins/audio/sounds/break_wall.mp3"),
  "player:destroy": new Audio("game/plugins/audio/sounds/player_dies.mp3"),
  "game:countdown": new Audio("game/plugins/audio/sounds/countdown.mp3"),
  "explosion:start": new Audio("game/plugins/audio/sounds/explosion.mp3"),
  "bomb-placed:start": new Audio("game/plugins/audio/sounds/place-bomb.mp3"),
  "power:up": new Audio("game/plugins/audio/sounds/powerup.mp3"),
  "life:lost": new Audio("game/plugins/audio/sounds/lost_life.mp3"),
};
Object.values(sounds).forEach(a => a.volume = 0.06);

const music = {
  // waiting: new Audio("game/plugins/audio/sounds/waiting.mp3"),
  waiting: new Audio("game/plugins/audio/sounds/bomberman_theme.mp3"),
};
music.waiting.loop = true;
music.waiting.volume = 0.06;

export function soundPlugin({ events }) {
  const cleanups = [];
  for (const [event, audio] of Object.entries(sounds)) {
    const unsub = events.on(event, () => {
      audio.currentTime = 0;
      audio.play();
    });
    cleanups.push(unsub);
  }
  cleanups.push(events.on("music:play", (name) => {
    const track = music[name];
    if (track) {
      track.currentTime = 0;
      track.play();
    }
  }));
  cleanups.push(events.on("music:stop", (name) => {
    const track = music[name];
    if (track) {
      track.pause();
      track.currentTime = 0;
    }
  }));
  return () => {
    Object.values(music).forEach(t => { t.pause(); t.currentTime = 0; });
    cleanups.forEach((fn) => fn());
  };
}