# TCP - The Chess Project
A traditional chess game that can be played online, either P2P or on a centralized server.

<div style="page-break-after: always;"></div>

## Name of team: 3 Idiotai
### Member list (Role):
* Aringas Civilka (Facilitator, Architect, Builder, Exhibitor)
* Adomas Pimpė (Architect, Builder, Critic)
* Jonas Kirkilovskis (Architect, Builder, Critic)
* Girius Frankonis (Designer, Builder, Architect)
## Use case diagram
![[use_case_diagram.png|500]]


<div style="page-break-after: always;"></div>

## Activity diagram
![[activity_diagram.png|400]]


<div style="page-break-after: always;"></div>

## High level architecture of the system
### Centralized Server
![[high_level_arch_centralized.png|500]]
### Peer to Peer
![[high_level_arch_p2p.png|500]]

<div style="page-break-after: always;"></div>

## Data flow diagram
![[data_flow_diagram.png|600]]

<div style="page-break-after: always;"></div>

## Technology stack:
* Client side - C++ application using SFML for UI and ASIO/Beast for networking
* Server side - Golang using gorrila/websocket
* Version control - git + Github
## Github repo & Proof of Concept: 
[https://github.com/aringq10/TCP](https://github.com/aringq10/TCP/tree/main/proof_of_concept)
![[proof_of_concept.png|450]]
## Project costs

|                    | Time spent by |             |           |          |         |          |           |
| ------------------ | ------------- | ----------- | --------- | -------- | ------- | -------- | --------- |
| Name               | Person        | Facilitator | Architect | Builder  | Critic  | Designer | Exhibitor |
| Aringas Civilka    | 17h15min      | 15min       | 7h50min   | 6h10min  | 2h      | 0h       | 0h        |
| Adomas Pimpė       | 14h30min      | 0h          | 8h        | 4h       | 2h30min | 0h       | 0h        |
| Jonas Kirkilovskis | 13h15min      | 0h          | 7h15min   | 4h30min  | 1h30min | 0h       | 0h        |
| Girius Frankonis   | 9h30min       | 0h          | 4h30min   | 3h       | 1h      | 1h       | 0h        |
| Total:             | 54h30min      | 15min       | 27h35min  | 17h40min | 7h      | 1h       | 0h        |
