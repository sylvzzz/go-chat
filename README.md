# go-chat

A TCP chat server and client written in Go, using only the standard library and the `sync` package.

## Layout

- `server/` - the chat server. Listens on `:8080`, accepts clients, and broadcasts each message to everyone else in the room.
- `client/` - a command line client. Connects to the server, asks for a username, sends what you type, and prints incoming messages.

Both are separate programs with their own `go.mod`.

## Running

Open two terminals for clients, plus one for the server. Port `8080` must be free before you start; if an old server is still running, the new one fails to bind and prints an error.

Terminal 1 (server):

```bash
cd server
make run
```

Terminal 2 and 3 (clients):

```bash
cd client
go run .
```

Wait for the username prompt, type a name, then start chatting. Type `exit` to leave.

## Example output

Server terminal with two clients connected, one of them saying hello:

```text
Server started ...
alice joined ...
bob joined ...
alice: hey everyone
bob: hi alice
alice left the chat ...
```

Client terminal for bob:

```text
Enter your username: bob
hi alice
```
