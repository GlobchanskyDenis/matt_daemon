DAEMON_BIN					=	daemon_bin
DAEMON_DIR					=	cmd/daemon
DAEMON_FILES				=	main.go      socket.go    config.go
DAEMON_FILENAMES			=	$(addprefix $(DAEMON_DIR)/,$(DAEMON_FILES))

MATT_BIN					=	Matt_daemon
MATT_DIR					=	cmd/matt
MATT_FILES					=	main.go   config.go
MATT_FILENAMES				=	$(addprefix $(MATT_DIR)/,$(MATT_FILES))

CLIENT_APP_BIN				=	client_bin
CLIENT_APP_DIR				=	cmd/clientApp
CLIENT_APP_FILES			=	main.go      socket.go
CLIENT_APP_FILENAMES		=	$(addprefix $(CLIENT_APP_DIR)/,$(CLIENT_APP_FILES))

all : $(MATT_BIN) $(DAEMON_BIN) $(CLIENT_APP_BIN)

$(MATT_BIN) : $(MATT_FILENAMES)
	@echo "компилирую бинарник запускающий демон как процесс"
	@go build -o $(MATT_BIN) $(MATT_FILENAMES)

$(DAEMON_BIN) : $(DAEMON_FILENAMES)
	@echo "компилирую бинарник демона"
	@go build -o $(DAEMON_BIN) $(DAEMON_FILENAMES)

$(CLIENT_APP_BIN) : $(CLIENT_APP_FILENAMES)
	@echo "компилирую бинарник приложения клиента"
	@go build -o $(CLIENT_APP_BIN) $(CLIENT_APP_FILENAMES)

fclean:
	@echo "удаляю бинарники"
	@rm -rf $(MATT_BIN)
	@rm -rf $(DAEMON_BIN)
	@rm -rf $(CLIENT_APP_BIN)

re: fclean all