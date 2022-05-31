# Websocket learning and simple echo server demo

# Overview of websockets

WebSocket is a technology that enables bidirectional, full-duplex
communication between client and server over a persistent, single-socket connection.

The base WebSocket protocol consists of an opening handshake (upgrading
the connection from HTTP to WebSockets), followed by data transfer. After the client and server successfully negotiate the opening handshake, the WebSocket connection acts as a persistent full-duplex communication channel. Each side can, independently, send data at will. Clients and servers transfer data back and forth in conceptual units referred to as messages, which can consist of one or more frames. Once the WebSocket connection has served its purpose, it can be terminated via a closing handshake.

The WebSocket protocol defines two URI schemes for traffic between server and client:
• ws, used for unencrypted connections.
• wss, used for secure, encrypted connections over Transport Layer Security (TLS). 

The rest of the WebSocket URI follows a generic syntax, similar to HTTP. It consists of several components. First the scheme (ws or wss), then host followed by a port, path, and query.  

`wss://example.com:443/websocket/demo?foo=bar`

The port component is optional with the default port for ws being 80 and 443 for wss.

## Opening handshake

The process of establishing a WebSocket connection is known as the opening handshake and consists of an HTTP/1.1 request/response exchange between the client and the server. 

The client always initiates the handshake. It sends a GET request to the server, indicating that it wants to upgrade the connection from HTTP to WebSockets. The server must return an HTTP `101 Switching Protocols` response code for the WebSocket connection to be established. 

Once that happens, the WebSocket connection can be used for ongoing, bidirectional, full-duplex communications between server and client.

The request must contain the following headers:
• Host
• Connection
• Upgrade
• Sec-WebSocket-Version
• Sec-WebSocket-Key

In addition to the required headers, the request may also contain optional ones. See the Opening handshake headers section later in this chapter for more information on headers.

The server must return an HTTP `101 Switching Protocols` response code for the
WebSocket connection to be successfully established. If the status code returned by the server is anything but HTTP `101 Switching Protocol`, the handshake will fail, and the WebSocket connection will not be
established.

Two of the required headers used during the opening handshake are
`Sec-WebSocket-Key`, and `Sec-WebSocket-Accept`. Together, these headers are essential in guaranteeing that both the server and the client are capable of communicating over WebSockets. 

`Sec-WebSocket-Key`, which is passed by the client to the server contains
a 16-byte, base64-encoded one-time random value (nonce). Its purpose is to help ensure that the server does not accept connections from non-WebSocket clients (e.g., HTTP clients) that are being abused (or misconfigured) to send data to unsuspecting WebSocket servers. An example of a `Sec-WebSocket-Key` would be `dGhlIHNhbXBsZSBub25jZQ`.

In direct relation to `Sec-WebSocket-Key`, the server response includes a `Sec-WebSocketAccept` header. This header contains a base64-encoded SHA-1 hashed value generatedmby concatenating the `Sec-WebSocket-Key` nonce sent by the client, and the static value `258EAFA5-E914-47DA-95CA-C5AB0DC85B11`.

This value `258EAFA5-E914-47DA-95CA-C5AB0DC85B11` is a magic string used by all websocket servers. See https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API/Writing_WebSocket_servers

If the `Sec-WebSocket-Key` header is missing from the client-initiated handshake, the server will stop processing the request and return an HTTP response with an appropriate error code (400 Bad Request, for example). If there’s something wrong with the value of `Sec-WebSocket-Accept`, or if the header is missing from the server response, the WebSocket connection will not be established (the client fails the connection)

## Message frames

After a successful opening handshake, the client and the server can use the WebSocket connection to exchange messages in full-duplex mode. A WebSocket message consists of one or more frames. The WebSocket frame has a binary syntax and contains several pieces of information.

• FIN bit - indicates whether the frame is the final fragment in a WebSocket message.
• RSV 1, 2, 3 - reserved for WebSocket extensions.
• Opcode - determines how to interpret the payload data.
• Mask - indicates whether the payload is masked or not.
• Masking key - key used to unmask the payload data.
• (Extended) payload length - the length of the payload.
• Payload data - consists of application and extension data.

## Fragmentation

There are numerous scenarios where fragmenting a WebSocket message into multiple frames is required (or at least desirable). For example, fragmentation is often used to improve performance. Without fragmentation, an endpoint would have to buffer the entire message before sending it. With fragmentation, the endpoint can choose a reasonably sized buffer, and when that is full, send subsequent frames as a continuation. The receiving endpoint then assembles the frames to recreate the WebSocket message. 

All data frames that comprise a WebSocket message must be of the same type
(text or binary); you can’t have a fragmented message that consists of both text and binary frames. However, a fragmented WebSocket message may include
control frames. See the Opcodes section later in this chapter for more details
about frame types.

The WebSocket protocol makes fragmentation possible via the first bit of the WebSocket frame — the FIN bit, which indicates whether the frame is the final fragment in a message. If it is, the FIN bit must be set to 1. Any other frame must have the FIN bit clear. 

## RSV

RSV1, RSV2, and RSV3 are reserved bits. They must be 0 unless an extension was negotiated during the opening handshake that defines non-zero values.

## Opcodes

Every frame has an opcode that determines how to interpret that frame’s payload data.

* 0 Continuation frame; continues the payload from the previous frame.
* 1 Indicates a text frame (UTF-8 text data).
* 2 Indicates a binary frame.
* 3-7 Reserved for custom data frames.
* 8 Connection close frame; leads to the connection being terminated.
* 9 A ping frame. Serves as a heartbeat mechanism ensuring the connection is still alive. The receiver must respond with a pong frame.
* 10 A pong frame. Serves as a heartbeat mechanism ensuring the connection is stillalive. Sent as a response after receiving a ping frame.
* 11-15 Reserved for custom control frames.


## Masking 

Each WebSocket frame sent by the client to the server needs to be masked with the help of a random masking-key (32-bit value). This key is contained within the frame, and it’s used to obfuscate the payload data. However, when data flows the other way around, the server must not mask any frames it sends to the client. 

A masking bit set to 1 indicates that the respective frame is masked (and
therefore contains a masking-key). The server will close the WebSocket
connection if it receives an unmasked frame.

On the server-side, frames received from the client must be unmasked before further processing. Masking is used as a security mechanism that helps prevent cache poisoning.

## Payload Length

The WebSocket protocol encodes the length of the payload data using a variable number of bytes:
• For payloads <126 bytes, the length is packed into the first two frame header bytes.
• For payloads of 126 bytes, two extra header bytes are used to indicate length.
• If the payload is 127 bytes, eight additional header bytes are used to indicate its length. 

## Payload Data

The WebSocket protocol supports two types of payload data: text (UTF-8 Unicode
text) and binary. Each frame’s payload type is indicated via a 4-bit opcode (1 for text or 2 for binary). 

## Closing Handshake

Compared to the opening handshake, the closing handshake is a much simpler process. It is initiated by sending a close frame with an opcode of 8. In addition to the opcode, the close frame may contain a body that indicates the reason for closing. This body consists of a status code (integer) and a UTF-8 encoded string (the reason).

The standard status codes that can be used during the closing handshake are defined by RFC 6455; additional, custom close codes can be registered with IANA21. 

* 0-999 N/A Codes below 1000 are invalid and cannot be used.
* 1000 Normal closure Indicates a normal closure, meaning that the purpose
for which the WebSocket connection was established has been fulfilled.
* 1001 Going away Should be used when closing the connection and there
is no expectation that a follow-up connection will be attempted (e.g., server shutting down, or browser navigating away from the page).
* 1002 Protocol error The endpoint is terminating the connection due to a
protocol error.
* 1003 Unsupported data The connection is being terminated because the
endpoint received data of a type it cannot handle (e.g. a text-only endpoint receiving binary data).
* 1004 Reserved Reserved. A meaning might be defined in the future.
* 1005 No status received Used by apps and the WebSocket API to indicate
that no status code was received, although one was
expected.
* 1006 Abnormal closure Used by apps and the WebSocket API to indicate that
a connection was closed abnormally (e.g., without
sending or receiving a close frame).
* 1007 Invalid payload data The endpoint is terminating the connection because it received a message containing inconsistent data (e.g., non-UTF-8 data within a text message).
* 1008 Policy violation The endpoint is terminating the connection because
it received a message that violates its policy. This is a generic status code; it should be used when other status codes are not suitable, or if there is a need to hide specific details about the policy.
* 1009 Message too big The endpoint is terminating the connection due to
receiving a data frame that is too large to process.
* 1010 Mandatory extension The client is terminating the connection because the
server failed to negotiate an extension during the opening handshake.
* 1011 Internal error The server is terminating the connection because it
encountered an unexpected condition that prevented it from fulfilling the request.
* 1012 Service restart The server is terminating the connection because it is
restarting.
* 1013 Try again later The server is terminating the connection due to a
temporary condition, e.g., it is overloaded.
* 1014 Bad gateway The server was acting as a gateway or proxy and
received an invalid response from the upstream server. Similar to 502 Bad Gateway HTTP status code.
* 1015 TLS handshake Reserved. Indicates that the connection was closed due
to a failure to perform a TLS handshake (e.g., the server certificate can’t be verified).
* 1016-1999 N/A Reserved for future use by the WebSocket standard.
* 2000-2999 N/A Reserved for future use by WebSocket extensions.
* 3000-3999 N/A Reserved for use by libraries, frameworks, and applications. Available for registration at IANA via firstcome, first-serve.
* 4000-4999 N/A Range reserved for private use in applications. 

Both the client and the server can initiate the closing handshake. Upon receiving a close frame, an endpoint (client or server) has to send a close frame as a response (echoing the status code received).

Once a close frame has been sent, no more data frames can pass over the
WebSocket connection. 

After an endpoint has both sent and received a close frame, the closing handshake is complete, and the WebSocket connection is considered closed. 

# Create a websocket demo

I want to understand the API of the `golang.org/x/net/websocket` package so I have used this package to create an echo server and also written some client code to connect to this server, send a message and receive a message.
