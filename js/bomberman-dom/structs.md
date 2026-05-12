## golang struct 

#### In Go, using map[string] with the entity IDs as keys:                                                               
```go                                                                                                                  
  type Entity struct {                           
      X  int `json:"x",omitempty`                                                                                               
      Y  int `json:"y",omitempty`
      Destroy bool `json:"destroy",omitempty`                                                                                               
  }

  type GameState struct {
      Stones     map[int]stone `json:"stones,omitempty"` // indestructible walls
      Players    map[int]player `json:"players,omitempty"`
      Bombs      map[int]bomb `json:"bombs,omitempty"`
      Explosions map[int]explosion `json:"explosions,omitempty"`
      Walls      map[int]wall `json:"walls,omitempty"`
  }
  ```

####  Building the payload:
```go
  //   Map keys are the entity ids
  state := GameState{
      Players: map[string]Entity{
          "1": {X: 34, Y: 96},
          "2": {X: 927, Y: 160},
          "3": {X: 500, Y: 200},
      },
      Bombs: map[string]Entity{
          "4": {X: 32, Y: 32},
      },
      // Explosions omitted — won't appear in JSON thanks to omitempty
      Walls: map[string]Entity{
          "5": {X: 64, Y: 32, destroy: true },
      },
  }

  json.Marshal(state) produces:

  {
    "players": {
      "1": {"x": 34, "y": 96},
      "2": { "x": 927, "y": 160},
      "3": {"x": 500, "y": 200}
    },
    "bombs": {
      "4": { "x": 32, "y": 32}
    },
    "walls": {
      "5": {"x": 64, "y": 32, "destroy": true}
    }
  }
```

#### Partial update (only player 1 moved):
```go
  update := GameState{
      Players: map[string]Entity{
          "1": {X: 36, Y: 96},
      },
  }
```
#### Destroy
```go
  Destroy:

  update := GameState{
      Players: map[string]Entity{
          "1": {dead: true},
      },
  }

```
Produces just {"players":{"1":{"id":1,"x":36,"y":96}}} — bombs, walls, explosions absent thanks to omitempty, so the
client's reducer keeps them unchanged.






BATCH PAYLOAD (BATCH DELTAS)

GAME_STATE = payload
REMOVE_PLAYER = player.id
DETONATE_BOMB = bomb.id
CHAT_MESSAGE = message

example (just an idea):
server => PLAYER_DIED = player.id -> reducer marks player as dead: true and setTimeOut(2secs) -> renderer sees player.dead and applies fadeout class => timer goes to 0 and event emits dispatch("remove_player") reducer removes player from state. 