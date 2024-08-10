# User file access checks via unix sockets

This is a proof of concept of verifying a requested user's access to a given file, assuming the file is on a given mounted file path. 

The idea is to simply give a username and a file path, and receive back "true" if the user can access the file and "false" if the user cannot. The intent is to avoid parsing group memberships or ACLs to infer the access level and rather let the system do that for you, for which the easiest method is to assume that user's identity.

Given that assuming a user's identity requires elevated privileges (i.e. root), this approach uses a client/server relationship and communicates via JSON over unix sockets. 

## Usage

The server process (vrfy.d) must first be running as root:

```
sudo go run server.go 
[sudo] password for mstack: 
Server is listening on /tmp/unix_socket
```

The client can then send its request

```
go run client.go 
Response received: {testuser /home/mstack/go/src/github.com/mstackhouse/verify/testfile.txt false}
```

## Where to go from here
- Client needs a simple CLI to take parameters of `user` and `path`
- Add the access level (i.e. rwx) for further clarity
