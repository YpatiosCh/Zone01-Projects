package lobby

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"bomberman/server/internal/player"
	"bomberman/server/internal/protocol"
)

type testEngine struct {
	mu    sync.Mutex
	ticks int
}

func (e *testEngine) GetState() any {
	e.mu.Lock()
	defer e.mu.Unlock()
	return map[string]any{"tick": e.ticks}
}

func (e *testEngine) Advance(_ map[int]map[string]bool) (any, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.ticks++
	return map[string]any{"tick": e.ticks}, true
}

func newTestPlayer(id string) *player.Player {
	return &player.Player{
		ID:   id,
		Name: id,
		Send: make(chan []byte, 1024),
	}
}

func waitFor(t *testing.T, timeout time.Duration, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("condition not met before timeout")
}

func TestLobbyConcurrentGameTraffic_NoDeadlock(t *testing.T) {
	l := NewLobby()
	go l.Run()
	t.Cleanup(l.Close)

	players := make([]*player.Player, 0, MaxPlayers)
	for i := 0; i < MaxPlayers; i++ {
		p := newTestPlayer(fmt.Sprintf("p-%d", i))
		players = append(players, p)
		if ok := l.RegisterPlayer(p); !ok {
			t.Fatalf("failed to register player %d", i)
		}
	}

	waitFor(t, 2*time.Second, func() bool {
		return l.GetPlayerCount() == MaxPlayers
	})

	stopDrainers := make(chan struct{})
	var drainWG sync.WaitGroup
	for _, p := range players {
		drainWG.Add(1)
		go func(pl *player.Player) {
			defer drainWG.Done()
			for {
				select {
				case <-stopDrainers:
					return
				case _, ok := <-pl.Send:
					if !ok {
						return
					}
				}
			}
		}(p)
	}

	l.gameMu.Lock()
	l.game = &testEngine{}
	l.gameMu.Unlock()
	l.setPhase(PhaseRunning)
	go l.runGameLoop()

	end := time.Now().Add(500 * time.Millisecond)
	var workerWG sync.WaitGroup
	for i := 0; i < 8; i++ {
		workerWG.Add(1)
		go func(seed int64) {
			defer workerWG.Done()
			r := rand.New(rand.NewSource(seed))
			for time.Now().Before(end) {
				p := players[r.Intn(len(players))]
				msgType := r.Intn(3)

				switch msgType {
				case 0:
					select {
					case l.Input <- player.PlayerMessage{
						Player: p,
						Msg: protocol.ClientMessage{
							Type:    "INPUT",
							Key:     "UP",
							Pressed: true,
						},
					}:
					default:
					}
				case 1:
					select {
					case l.Input <- player.PlayerMessage{
						Player: p,
						Msg: protocol.ClientMessage{
							Type: "REQUEST_STATE",
						},
					}:
					default:
					}
				default:
					l.broadcastJSON(protocol.ServerMessage{
						Type: "NOOP",
						Data: map[string]any{"ts": time.Now().UnixNano()},
					})
				}
			}
		}(time.Now().UnixNano() + int64(i))
	}

	done := make(chan struct{})
	go func() {
		workerWG.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("workers did not finish; possible deadlock")
	}

	close(stopDrainers)
	drainWG.Wait()
}

func TestLobbyConcurrentRegisterUnregister_NoDeadlock(t *testing.T) {
	l := NewLobby()
	go l.Run()
	t.Cleanup(l.Close)

	const n = 24
	players := make([]*player.Player, 0, n)
	for i := 0; i < n; i++ {
		players = append(players, newTestPlayer(fmt.Sprintf("u-%d", i)))
	}

	var registerWG sync.WaitGroup
	for _, p := range players {
		registerWG.Add(1)
		go func(pl *player.Player) {
			defer registerWG.Done()
			_ = l.RegisterPlayer(pl)
		}(p)
	}

	registerDone := make(chan struct{})
	go func() {
		registerWG.Wait()
		close(registerDone)
	}()

	select {
	case <-registerDone:
	case <-time.After(3 * time.Second):
		t.Fatal("register phase timed out")
	}

	var unregisterWG sync.WaitGroup
	for _, p := range players {
		unregisterWG.Add(1)
		go func(pl *player.Player) {
			defer unregisterWG.Done()
			select {
			case l.GetUnregisterChan() <- pl:
			case <-time.After(2 * time.Second):
				t.Errorf("unregister send timed out for player %s", pl.ID)
			}
		}(p)
	}

	unregisterDone := make(chan struct{})
	go func() {
		unregisterWG.Wait()
		close(unregisterDone)
	}()

	select {
	case <-unregisterDone:
	case <-time.After(3 * time.Second):
		t.Fatal("unregister phase timed out")
	}
}

func TestLobbyRepeatedCloseWithCountdown_NoPanic(t *testing.T) {
	l := NewLobby()
	go l.Run()

	p1 := newTestPlayer("c1")
	p2 := newTestPlayer("c2")
	if !l.RegisterPlayer(p1) || !l.RegisterPlayer(p2) {
		t.Fatal("failed to register countdown players")
	}

	waitFor(t, 2*time.Second, func() bool {
		return l.GetPlayerCount() >= 2
	})

	var closeWG sync.WaitGroup
	for i := 0; i < 20; i++ {
		closeWG.Add(1)
		go func() {
			defer closeWG.Done()
			l.Close()
		}()
	}

	done := make(chan struct{})
	go func() {
		closeWG.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("concurrent Close() calls timed out")
	}
}
