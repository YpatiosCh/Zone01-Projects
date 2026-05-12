SOLUTION_DIR := solution
GAME_DIR     := game
ENGINE       := $(GAME_DIR)/linux_game_engine
ROBOTS_DIR   := $(GAME_DIR)/linux_robots
MAPS_DIR     := $(GAME_DIR)/maps

export CARGO_TARGET_DIR := $(CURDIR)/target
PLAYER := $(CARGO_TARGET_DIR)/release/solution

DEFAULT_MAP := map00
DEFAULT_BOT := bender

CARGO := cargo
CARGO_FLAGS := --manifest-path $(SOLUTION_DIR)/Cargo.toml

.PHONY: all build run test clean build-docker run-docker

all: build

build:
	$(CARGO) build --release $(CARGO_FLAGS)

test:
	$(CARGO) test $(CARGO_FLAGS)

run: build
	$(ENGINE) -f $(MAPS_DIR)/$(DEFAULT_MAP) -p1 $(PLAYER) -p2 $(ROBOTS_DIR)/$(DEFAULT_BOT) -q

build-docker:
	docker build -t filler .

run-docker: build-docker
	docker run -v "$(CURDIR)/solution":/filler/solution -it filler

clean:
	$(CARGO) clean $(CARGO_FLAGS)
	rm -rf $(CARGO_TARGET_DIR)



## p1 me p2 bender map00
# ./linux_game_engine -f ./maps/map00 -p1 ./solution/target/release/solution -p2 ./linux_robots/bender

## p1 me p2 wall_e map00
# ./linux_game_engine -f ./maps/map00 -p1 ./solution/target/release/solution -p2 ./linux_robots/wall_e

## p1 me p2 h2_d2 map00
# ./linux_game_engine -f ./maps/map00 -p1 ./solution/target/release/solution -p2 ./linux_robots/h2_d2

## p1 me p2 terminator map00
# ./linux_game_engine -f ./maps/map00 -p1 ./solution/target/release/solution -p2 ./linux_robots/terminator
