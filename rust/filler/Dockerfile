FROM rust:1.86-slim-bookworm

COPY ./game/maps			        /filler/maps
COPY ./game/linux_robots		    /filler/linux_robots
COPY ./game/linux_game_engine	/filler/linux_game_engine

WORKDIR /filler/

ENTRYPOINT ["/bin/bash"]