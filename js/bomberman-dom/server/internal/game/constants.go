package game

var ()

const (
	P2ArenaWidth  = 19
	P2ArenaHeight = 13

	P3ArenaWidth  = 25
	P3ArenaHeight = 13

	P4ArenaWidth  = 31
	P4ArenaHeight = 13

	CellSize   = 1.0
	wallWidth  = 1.0
	wallHeight = 1.0

	maxLives = 3

	Circle    = 0
	Rectangle = 1

	PlayerSpeed     = 0.06
	speedPowerupInc = 0.008

	ExplosionStartSize = 1

	BombsMax = 1

	BombDelayFrames = 130

	BombCooldownMax = 60
	EndOfGameTimer  = 10 * 60

	BombPushPower = 0.20
	impulseBleed  = 0.90

	ExplodeDelay = 5

	InvincibleDuration = 30

	BombPlaceStopTimer = 10
)

var cardinalDirections = [4]pos{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
var diagonalDirections = [4]pos{{-1, -1}, {1, 1}, {-1, 1}, {1, -1}}
