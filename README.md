# chesster-server
(WIP) Back end for Chess server

This is meant to work with the Android app [here](https://github.com/cactorium/chesster-app)

Maybe eventually include a web version that'll run from the same server

## How to build

Install [Go](https://golang.org/doc/install) and [protocol buffers for Go](https://developers.google.com/protocol-buffers/docs/gotutorial) (see section for "Compiling your protocol buffers").

 TODO finish and test instructions

```
make # Or go build ./cmd/chessterd
```

## Basic architecture

- [ ] SQLite database for backend
  - [ ] Stores user login info (username + salted password hash) [TODO: explore using additional authentication using a key stored in the Android app that can be revoked by the user]
  - [ ] Stores board history
  - [ ] Stores latest board configuration
  - [ ] (Extra feature): friend's list
  - [ ] Allow spectating on (public) matches
- [ ] Chess engine
  - [x] Validates moves
  - [ ] Validates board history
  - [ ] Determines checkmate/stalemate/draws
  - [x] Flags mates
- [ ] TCP-based server
  - [ ] Authenticates users and provides secure tokens
  - [ ] Validates packets
  - [ ] Allows account management
  - [ ] Provides push notifications
  - [ ] Password recovery
  - [ ] Specify packet format
  - [ ] Specify payloads using Protobuf
- [ ] CLI client (used for testing server)
  - [ ] Can communication with server
  - [ ] Can read in test patterns from stdin for integration testing
  - [ ] Can log in
  - [ ] Can list matches
  - [ ] Can start match
  - [ ] Can make a move
  - [ ] Can display board
  - [ ] Lets you play chess from the command line
- [ ] Maybe eventually web app
  - [ ] Lets you play chess online

## Protocol
We'll use a gRPC server for communication, with all messages being serviced by the gRPC service on TCP port 8888.
Authentication will be done using an authentication service on 8888, with all other service types containing a header with authentication details.
Push notifications will be done by a long-term connection to a particular endpoint.

### Message header

|Name            | Protobuf type    |Description                    |
|----------------|-----------|-------------------------------|
| Version        | varint    | Version number, currently 0 |
| Type           | enum      | Packet type |
| Session token  | bytes     | Authentication token received in handshake |
| Server nonce   | bytes     | Nonce generated by server |
| Client none    | bytes     | Nonce chosen by client |
| HMAC           | bytes     | Hash((Key1 xor padding1) + Hash((Key2 xor padding2) + snonce + message + cnonce + token)) |

Key1 and Key2 will be different depending on whether the packet originated from the client, or server, and all four keys will be generated during the handshake

### Handshake
The client will send an authentication request with a randomly generated nonce and username.

The server will reply with its own nonce, the salt for the given user, and a session token.
The server will generate the session keys using its salted password, the client nonce, and its nonce, and use it to generate the HMAC in the packet header.

The client will send an acknowledgement with a valid HMAC generated from the same components, and the handshake is completed.
The acknowledgement may include some additional randomly generated bytestrings to use to generate encryption keys, which would then be used to encrypt packets within this session.

The session token is valid for two hours, and at any point the client can send a request to generate a new session token.
The old token will be revoked after the new session token has been successfully used.
There will be a packet type to allow logging out and elimination of the session information.

A guest session uses a Diffie-Hellman key exchange since there's a lack of a shared secret between the server and the user.

TODO: flesh out with hash function and encryption choices

### Packet types
TODO
