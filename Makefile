DAEMON_BIN					=	daemon_bin
DAEMON_DIR					=	cmd/daemon
DAEMON_FILES				=	main.go      socket.go    config.go
DAEMON_FILENAMES			=	$(addprefix $(DAEMON_DIR)/,$(DAEMON_FILES))

MATT_BIN					=	Matt_daemon
MATT_DIR					=	cmd/matt
MATT_FILES					=	main.go   config.go
MATT_FILENAMES				=	$(addprefix $(MATT_DIR)/,$(MATT_FILES))

TEST_SERVER_BIN				=	server
TEST_SERVER_DIR				=	cmd/socketServer
TEST_SERVER_FILES			=	main.go
TEST_SERVER_FILENAMES		=	$(addprefix $(TEST_SERVER_DIR)/,$(TEST_SERVER_FILES))

all : $(MATT_BIN) $(DAEMON_BIN) $(TEST_SERVER_BIN)

$(MATT_BIN) : $(MATT_FILENAMES)
	@echo "компилирую бинарник запускающий демон как процесс"
	@go build -o $(MATT_BIN) $(MATT_FILENAMES)

$(DAEMON_BIN) : $(DAEMON_FILENAMES)
	@echo "компилирую бинарник демона"
	@go build -o $(DAEMON_BIN) $(DAEMON_FILENAMES)

$(TEST_SERVER_BIN) : $(TEST_SERVER_FILENAMES)
	@echo "компилирую бинарник тестового сервера"
	@go build -o $(TEST_SERVER_BIN) $(TEST_SERVER_FILENAMES)

test :
	rm -rf client_test
	rm -rf tcp_server_test
	rm -rf udp_server_test
	go build -o client_test pkg/netSocket/deprecated/client.go
	go build -o tcp_server_test pkg/netSocket/deprecated/tcpServer.go
	go build -o udp_server_test pkg/netSocket/deprecated/udpServer.go

fclean:
	@echo "удаляю бинарники"
	@rm -rf $(MATT_BIN)
	@rm -rf $(DAEMON_BIN)
	@rm -rf $(TEST_SERVER_BIN)

re: fclean all