nats-randomdata
---

Random Data Generator for NATS

The goal is to have a simple way of generating a constant stream of random data to at custom subject. Useful for debugging, exploring features, etc.

Requires the `nsc` command to grab the appropriate credentials for connecting to the cluster.

WIP - currently you can send a constant flow of "order data" to a single subject at a specified interval. Maybe more to come, maybe not - this does what I want for now.


### Example

```shell
nats-randomdata push orders.new -c nsc://mainbus/website/pushtest -i 1
```

This will push random order data to the `orders.new` subject ever 1 second using the `pushtest` user in account `website` from operator `mainbus`.
