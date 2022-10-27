.PHONY: server_local
server_local:
	DOCKER_CONTEXT=local ENV=local zsh ./scripts/rebuild_server.sh