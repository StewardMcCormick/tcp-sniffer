# Zero-Allocation TCP Analyzer

[![Test](https://github.com/StewardMcCormick/tcp-sniffer/actions/workflows/TEST.yml/badge.svg)](https://github.com/StewardMcCormick/tcp-sniffer/actions/workflows/TEST.yml)

## Flow обработки пакетов

```mermaid
flowchart TD
    A[Sniffer] --> B[Parser]

    B -->|Invalid / Unsupported| X[Drop]
    B -->|Valid Packet| C[Engine]

    C --> D[Direction Detection]
    D --> E[Canonical Flow Key]
    E --> F[Storage Get / Create]

    F --> G{New Flow?}
    G -->|No| K[Retransmission Detection]
    G -->|Yes| H{SYN without ACK?}

    H -->|Yes| I[Track TCP State]
    H -->|No| J[Midstream Initialization]

    I --> K
    J --> K

    K --> L{Retransmission?}
    L -->|Yes| M[Skip RTT Calculation]
    L -->|No| N[Calculate RTT / SRTT]

    M --> O[Build Stats Update]
    N --> O

    O --> P(Apply)

```
