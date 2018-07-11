#include <sys/types.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <stdio.h>

#define SOCK_FILE "socket"

int main() {
	int sock, msgsock, rval;
	struct sockaddr_un server;
	char buf[256];

	sock = socket(AF_UNIX, SOCK_STREAM, 0);
	if(sock < 0) {
		perror("failed to open new unix stream socket");
		exit(1);
	}

	server.sun_family = AF_UNIX;
	strcpy(server.sun_path, SOCK_FILE);

	if(bind(sock, (struct sockaddr *) &server, sizeof(struct sockaddr_un))) {
		perror("failed to bind unix stream socket");
		exit(1);
	}

	printf("opened new socket: %s\n", server.sun_path);
	listen(sock, 5);

	for(;;) {
		msgsock = accept(sock, 0, 0);
		if(msgsock == -1)
			perror("failed to accept new connection");
		else do {
			bzero(buf, sizeof(buf));
			if((rval = read(msgsock, buf, 256)) < 0) {
				perror("failed to read stream message");
			} else if(rval == 0) {
				printf("closing connection\n");
			} else {
				printf("-->%s\n", buf);
				if(write(sock, buf, sizeof(buf)) < 0) {
					perror("failed to echo back stream message");
				} else {
					printf("<--%s", buf);
				}
			}
		} while(rval > 0);
		close(msgsock);
	}

	close(sock);
	unlink(SOCK_FILE);
}
