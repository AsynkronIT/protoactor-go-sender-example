This repository contains supplemental examples for my blog article, [[Golang] Protoactor-go 101: How actors communicate with each other](https://blog.oklahome.net/2018/09/protoactor-go-messaging-protocol.html) , to cover all message passing methods for all kinds of actors provided by protoactor-go.

![](https://raw.githubusercontent.com/oklahomer/protoactor-go-sender-example/master/docs/components.png)

# Local
For local message passing, see below directory:
- local-tell ... Use Tell() for local message passing. The recipient actor cannot refer to the sender actor.
- local-request ... Use Request() for local message passing. The recipient actor can refer to the sender actor.
- local-future ... Use RequestFuture() for local message passing. Context.Sender() does not return the PID of sender actor but that of actor.Future.

# Remote
- remote/messages ... Contain Protobuf serializable message structures.
- rmeote/remote-pong ... A process that returns pong message to sender.
- remote/remote-ping-tell ... A process that sends message to pong actor by Tell(). The recipient cannot refer to the sender actor.
- remote/remote-ping-request ... A process that sends message to pong actor by Request(). The recipient actor can refer to the sender actor.
- remote/remote-poing-future ... A process that sends message to pong actor by RequestFuture(). Context.Sender() does not return the PID of sender actor but that of actor.Future.

# Cluster Grian
- cluster/messages ... Contain Protobuf serializable message structures and generated actor.Actor implementation for gRPC based communication.

## Cluster Grain usage with remote communication
- cluster/cluster-pong-remote ... A process that returns pong message to the sender based on remote actor implementation.
- cluster/cluster-ping-tell ... A process that sends message to pong actor by Tell(). The recipient cannot refer to the sender actor.
- cluster/cluster-ping-request ... A process that sends message to pong actor by Request(). The recipient actor can refer to the sender actor.
- cluster/cluster-ping-future ... A process that sends message to pong actor by RequestFuture(). Context.Sender() does not return the PID of sender actor but that of actor.Future.

## Cluster Grain usage with gRPC based communication
- cluster/cluster-pong-grpc ... A process that returns pong message to the sender via gRPC service.
- cluster/cluster-ping-grpc ... A process that sends message to pong actor over gRPC based service.

# References
- [[Golang] Protoactor-go 101: Introduction to golang's actor model implementation](https://blog.oklahome.net/2018/07/protoactor-go-introduction.html)
- [[Golang] Protoactor-go 101: How actors communicate with each other](https://blog.oklahome.net/2018/09/protoactor-go-messaging-protocol.html)