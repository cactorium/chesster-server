# chesster-server
(WIP) Back end for Chess server

This is meant to work with the Android app [here](https://github.com/cactorium/chesster-app)

Maybe eventually include a web version that'll run from the same server

## Basic architecture

- [ ] SQLite database for backend
  - [ ] Stores user login info (username + salted password hash) [TODO: explore using additional authentication using a key stored in the Android app that can be revoked by the user]
  - [ ] Stores board history
  - [ ] Stores latest board configuration
  - [ ] (Extra feature): friend's list
  - [ ] Allow spectating on (public) matches
- [ ] Chess engine
  - [ ] Validates moves
  - [ ] Validates board history
  - [ ] Determines checkmate/stalemate
  - [ ] Flags mates
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
We'll use a TCP socket for communication with the Android app, using Protobuf-encoded messages with a simple header for communication.
These can be plaintext with some MAC code to ensure security (I don't think anyone cares about their match's privacy __that__ much).
Authentication will use an elliptic curve Diffie-Hellman exchange to generate a set of secrets for authenticating all the following requests.

### Header format

|Name            | Length    |Description                    |
|----------------|-----------|-------------------------------|
| Length         | 2         | Message length in bytes excluding header |
| Version        | 2         | Version number, currently 0 |
| Type           | 4         | Packet type |
| Session token  | 8         | Authentication token received in handshake |
| Server nonce   | 8         | Nonce generated by server |
| Client none    | 8         | Nonce chosen by client |
| HMAC           | 4         | Hash((Key1 xor padding1) + Hash((Key2 xor padding2) + snonce + message + cnonce + token)) |

All components will be in big-endian format

Key1 and Key2 will be different depending on whether the packet originated from the client, or server, and all four keys will be generated during the handshake

### Handshake
TODO: needs to generate a session token and key for each session

Handshake will be initiated a packet of type zero, where the HMAC is ignored.

### Packet types
TODO
