# gored
gored is a in-memory database similar to Redis. Lately I took interest in database development and decided to implement a simple Database.
I chose to base it on Redis due to its power and simplicity. Compared to relational databases there is minimal parsing, the persistence is implemented with AOF(append only file)
and there is no need to write a VM. Also it touches on distributed systems and concurrent programming.
gored is able to parse RESP and send messages using this protocol to redis client.
## To Do
- Create server type that allows to handle multiple connections
- Implement connection abstraction
- Pub/Sub
- Pipelining
- Add remaining data structures (set, sorted set etc)

## References
Build your own Redis challenge and  Build Your Own Redis with C/C++ book.
