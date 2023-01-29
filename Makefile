.PHONY: server_local
server_local:
	DOCKER_CONTEXT=local ENV=local zsh ./scripts/rebuild_server.sh

server_remote:
	DOCKER_CONTEXT=remote ENV=prod zsh ./scripts/rebuild_server.sh