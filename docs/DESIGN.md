# go-karai General Design Spec

## Basic Functionality of the go-karai node

The go-karai node needs to at a minimum do the following:

-   generate/manage a TRTL wallet (custom prefix?)
-   send TRTL transactions containing tx_extra data
-   accept a user input payload, determine its size
-   communicate w/ pinning server
-   connect to ipfs
-   create or connect to a transaction channel
-   parse channel events
-   coordinate transaction linkages
-   assign transaction timestamps

## Vision

The go-karai node provides a means for a user to establish and interact with or coordinate transactions within a transaction channel. The go-karai node receives input intended as a script or payload from the user and is given a cost of pinning that script or payload from the desired pinning server.

The go-karai node sends the size of the intended script or payload in bytes to the pinning server, which returns a price for pinning it and a window of time for which the price of pinning is valid. If the user sends payment and payload within the time alotted, the payload is pinned and the content address is generated and released by the pinning server to the go-karai user.

When a transaction channel has been created, the node that created the channel is able to define a set of variable for transaction stream defaults. This go-karai node will also serve as a transaction coordinator for that channel. Transactions are encouraged to move chronologically, and transaction stream width is elastic according to current transaction volume, similar to Cryptonote dynamic block size.

A stream coordinator is able to set limits on stream width, transaction frequency, and snapshotting intervals. A 'snapshot' is a name for the user initiated action of storing a merkle hash on the TRTL chain as proof that a transaction ledger has not been modified between snapshots. It is not required, but is a mechanism to establish trust that the transaction stream has not been altered.

## Coordinating Transactions

Transaction volume is not always predictable, and the mesh must be able to assemble and process these interlinked transactions as fast as possible. To compensate for bursts of transaction volume, an N+1 stream width management system has been proposed to create expanding transaction wave capacity similar to TRTL's elastic block size.

#### Example Scenario

A stream is initiated with a default configuration of N+1 stream width, and pins its first script. The script is announced, and an influx of users arrive, generating 10 transaction events within the user defined receive-interval (in this example, we'll use 2 seconds as the receive interval). The stream mesh width will widen with interlinked transactions by 1 extra transaction per wave.

The process for handling these transactions is as follows:

0. Individual transactions are timestamped in the order they were received in the interval,the transactions are counted and assigned slots in a stream lattice. When the lattice order has been determined, the previous snapshot hash is serialized with the transaction array and hashed to form the current snapshot hash.
1. Block wave 1: `SH[0], tx (1A)` (First transaction enters slot A, stream width = 1) 9 tx remain
1. Block wave 2: `SH[0], tx 2A->(1A)`, `tx 2B->(1A)` (Next 2 tx enter slots A+B, stream width = 2) 7 tx remain
1. Block wave 3: `SH[0], tx 3A->(2A)`, `tx 3B->(2A & 2B)`, `tx 3C->(2B)` (Next 3 tx enter slots A+B+C, stream width = 3) 4 tx remain
1. Block wave 4: `SH[0], tx 4A->(3A)`, `tx 4B->(3A & 3B)`, `tx 4C->(3B & 3C)`, `tx 4D->(3C)` (Final 4 tx enter slots A+B+C+D, stream width = 4) No tx remain
1. Optional snapshot send to TRTL chain

#### Snapshot Diagram

A snapshot is not a backup. Periodic snapshots are wise but not necessary. A snapshot incurs a cost in the form of a fee paid to store the snapshot string in the tx_extra field of a transaction on the TRTL Network chain.

A snapshot string consists of the following serialized elements:

```
[Prefix][Snapshot ID][# of tx since snapshot N[self-1]][# of tx waves][Snapshot Hash]
```

Snapshot string example: `AA248032033ac35179663b5654d83e4f4bcf805e`

-   Prefix is `AA`
-   Snapshot ID is `24`, because 23 snapshots have happened before it
-   Number of tx in the lattice since last snapshot is `80`
-   Number of waves those transactions occupied was `32`
-   Snapshot Hash was `033ac35179663b5654d83e4f4bcf805e`

A snapshot hash consists of the hashed representation of the following elements:

```
[Hashed array of transactions] + [Previous Snapshot Hash]
```

Snapshot hash example: `6c7d48e42e0a0fde48064c4ec6ce7dc6`

-   `[Prefix]` is a static (0-9a-f) identifying characteristic to allow for easy scanning of a set of TRTL transactions containing a go-karai snapshot.
-   `[Snapshot ID]` is an incrementing number assigned to the snapshot group of waves and their constituent transactions. If you never snapshot, this number will stay 0. If you then create a snapshot, that number will increment to 1, and so on.
-   `[# of tx]` this is the sum total number of transactions marked with this Snapshot ID. By storing the number of transactions in a snapshot, it adds an added level of complexity to a possible attacker trying to generate a hash collision to alter the history undetected.
-   `[# of tx waves]` this is the sum total number of transaction waves containing transactions marked with this Snapshot ID. By storing the number of transactions in a snapshot, it adds an added level of complexity to a possible attacker trying to generate a hash collision to alter the history undetected.
-   `[Snapshot Hash]` This is a hashed representation of the following two serialized elements:
    -   Hashed array of transactions
    -   Previous snapshot hash
-   `[Hashed array of transactions]` is a hash made from the serialized transactions
