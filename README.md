# CoolCoin v 0.0.1
## Concept
### Blocks and Transactions
#### Blocks
a **block** consists of the following:
- block header
    - the producer's address
    - producer's signature
- block data
    - height (incrementing 64-bit integer)
    - timestamp (unix time, in milliseconds)
    - previous block hash
    - transactions

a **transaction** consists of the following:
- sender address
- receiver address
- data (transferred amount)
- nonce

### Consensus Algorithm
- PoA (Proof of Authority)
- if more than 50 % of pre-set validators agree for a given block, then the block is considered valid
- blocks are not produced when 50% of pre-set validators are against or offline.

- AccountManager
    - GetAddress
    - GetAccount

