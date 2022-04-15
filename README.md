# Simple Replicated Database

# About
This db is very bad for financial database and good for use cases with eventual consistency - maybe root comments (and not parent-child path), news updates.

# Design
We have a static leader and configurable follower count. All the writes go to leader and get queries a random follower.
Each follower periodically (configurable) queries leader for state and leader responds with few records. These records are used by follower to update it's local store.

For now we have just inc(key, value) and get(key) operation.

# Results
Read demo code and checkout the logs

# Further Steps

## Transactions
For transaction, we can have a list of these operations and these must be executed atomically. We can just lock a store, execute the transaction and unlock the store. Locking an entire store is a bad idea since we prevent other keys from either read / writes.

If we can have predefined transaction involving keys - we can store a hash of these keys and just lock this particular hashed key, otherwise
we need to lock list of keys sequentially (or use wait group) and then modify the store.

With multiple leaders , we need to have a separate co-ordinater to perform these operations on (n/2 + 1) leaders and read from (n/2 + 1) leaders.

This gets messy and is out of scope for me.