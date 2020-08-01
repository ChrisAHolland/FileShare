# FileShare
FileShare is an effort to prototype a distributed, peer-to-peer file sharing (or torrenting) network. In the first phase of development for FileShare, the system was a client/server based architecture in which the Peers could send files to the SwarmMaster (server) and request files from the SwarmMaster that other Peers uploaded. In the second phase of development, and for development in the future, FileShare has moved to a peer-to-peer style architecture. 

## Contents
This project was started as a project for CSC-462 (Distributed Systems) at UVIC.  
* `src/`: Contains the code for FileShare, a peer-to-peer distributed file sharing system.
    * `fileshare/`: Contains the Peer and SwarmMaster code, all part of the `fileshare` package.  
    * `labgob/`: Code provided for course labs. Wrapper for the `encoding/gob` package.  
    * `labrpc/`: Code provided for course labs. Simple RPC framework.  
    * `main/`: Contains `main.go` and Peer directories, used for running a test case of the FileShare system.  
* `images/`: Contains architecture related images for the documentation.   

## Peer-to-peer System Architecture
FileShare is a distributed, client/server + peer-to-peer file sharing network (think BitTorrent, Gnutell, or Napster). A high level architecture diagram can be seen below:  

![](images/FileShare-P2P.png)  

## Data flow and Behaviour Diagrams
All components of FileShare (Peers and the SwarmMaster) communicate using Remote Procedure Calls (RPC) over a TCP connection. Below you will find behaviour diagrams for the main functionalities of the system.  

### Peer Connecting to SwarmMaster
![](images/peer-connect-server.png)

### Peer Connecting to Peer
![](images/peer-connect-peer.png)

### Registering a File
![](images/register-file.png)

### Searching for a File
![](images/search-for-file.png)
