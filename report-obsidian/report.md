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
![[data-flow.png|600]]
<div style="page-break-after: always;"></div>

## Technology stack:
* Client side - C++ application using SFML for UI and ASIO/Beast for networking
* Server side - Golang using gorrila/websocket
* Version control - git + Github
## Github repo & Proof of Concept: 
[https://github.com/aringq10/TCP](https://github.com/aringq10/TCP/tree/main/proof_of_concept)
![[proof_of_concept.png|450]]

## UML class diagram after TP2 (with plans for TP3)
![[class_diagram.png|500]]
## Project costs

### TP1
|                    | Time spent by |             |           |         |        |          |           |
| -----------------  | ------------- | ----------- | --------- | ------- | ------ | -------- | --------- |
| Name               | Person        | Facilitator | Architect | Builder | Critic | Designer | Exhibitor |
| Aringas Civilka    | 4h40min       | 5min        | 4h        | 35min   | 0h     | 0h       | 0h        |
| Adomas Pimpė       | 4h30min       | 0h          | 4h        | 0h      | 30min  | 0h       | 0h        |
| Jonas Kirkilovskis | 4h            | 0h          | 4h        | 0h      | 0h     | 0h       | 0h        |
| Girius Frankonis   | 4h            | 0h          | 4h        | 0h      | 0h     | 0h       | 0h        |
| Total:             | 17h10min      | 5min        | 16h       | 35min   | 30min  | 0h       | 0h        |

### TP2
|                    | Time spent by |             |           |         |         |          |           |
| ------------------ | ------------- | ----------- | --------- | ------- | ------- | -------- | --------- |
| Name               | Person        | Facilitator | Architect | Builder | Critic  | Designer | Exhibitor |
| Aringas Civilka    | 12h35min      | 10min       | 3h50min   | 5h35min | 2h      | 0h       | 0h        |
| Adomas Pimpė       | 10h           | 0h          | 4h        | 4h      | 2h      | 0h       | 0h        |
| Jonas Kirkilovskis | 10h15min      | 0h          | 3h15min   | 5h30min | 1h30min | 0h       | 0h        |
| Girius Frankonis   | 8h            | 0h          | 30min     | 5h      | 1h      | 1h30min  | 0h        |
| Total:             | 40h50min      | 10min       | 11h35min  | 20h5min | 6h30min | 1h30min  | 0h        |

<div style="page-break-after: always;"></div>

### TP1 + TP2
|                    | Time spent by |             |           |          |         |          |           |
| ------------------ | ------------- | ----------- | --------- | -------- | ------- | -------- | --------- |
| Name               | Person        | Facilitator | Architect | Builder  | Critic  | Designer | Exhibitor |
| Aringas Civilka    | 17h15min      | 15min       | 7h50min   | 6h10min  | 2h      | 0h       | 0h        |
| Adomas Pimpė       | 14h30min      | 0h          | 8h        | 4h       | 2h30min | 0h       | 0h        |
| Jonas Kirkilovskis | 14h15min      | 0h          | 7h15min   | 5h30min  | 1h30min | 0h       | 0h        |
| Girius Frankonis   | 12h           | 0h          | 4h30min   | 5h       | 1h      | 1h30min  | 0h        |
| Total:             | 58h           | 15min       | 27h35min  | 20h40min | 7h      | 1h30min  | 0h        |
