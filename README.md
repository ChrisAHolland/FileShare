# FileShare
Prototyping a distributed, peer-to-peer file sharing network.

## Contents
This project was started as a project for CSC-462 (Distributed Systems) at UVIC.  
* `cs/`: Contains the initial, client/server style prototype for the system.  
* `p2p/`: Contains the peer-to-peer architecture system.

## Client/Server System Architecture
The initial prototype of the system is a client/server architecture in which Peers can send and receive files to/from the Server (SwarmMaster) as demonstrated below:  

![](images/FileShare-CS.png)  

A typical use case scenario of the client/server FileShare system is demonstrated in the behaviour diagram below:  

![](images/FileShare-Behaviour-Diagram.png)  

## Peer-to-peer System Architecture
I am currently working on making FileShare a distributed, peer-to-peer system. Initially, I plan to have it's system architecture similar to what the below diagram represents:  

![](images/FileShare-P2P.png)  

There is more to come for this part of the project!