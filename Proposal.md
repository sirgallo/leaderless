# Athn

## A Proposal for a Leaderless Distributed Consensus Algorithm with Global Versioning and Reputation Based Conflict Resolution


`Author: Nicholas Gallo - Nov, 2023`


## Proposal 

This proposal demonstrates a system that is meant to mimic the transfer of information between people, where information is spread out similar to an epidemic model, and information from more reputable sources, or the same information coming from many unique individuals, can determine the speed or scale at which data is accepted. Consensus on the global state in the system is managed through a combination of data propagation through a gossip protocol and proposal redunancy, global versioning inspired by `multi-version concurrency control`, and a reputation based conflict resolution system.


## Overview

A request to update state will create a proposal, which is a mutated copy of the current state with the update from the request on the node that received it in the system. A `vdf` is then computed with the `current version of the state + 1` to symbolize the next version with the updated state change proposal. The proposal is then propagated through the nodes of the network using a gossip protocol, with each node updating its current copy of the data with the proposed version. A reputation based model will be used to resolve conflicts where two nodes submit a data copy with the same version at the same time, favoring nodes with a long track record of accurate state change updates.


## Reputation Based Conflict Resolution Model

When two nodes attempt to modify the same version of the state concurrently and both submit new state proposals, the proposal is then propagated through the network of nodes using the `proposal redundancy` model, where a singular copy of the version will eventually be copied up to the `redundancy threshold` times or greater on a majority of nodes in the network. 

However, when there is conlict while quorum is being reached, concurrent writes to the same proposal version are handled through a reputation based conlict resolution model where the state change proposal coming from a node with a longer track record of historical successful writes to the state machine will be weighted heavier over other state change proposals with the same version tag in conjunction with a weight applied for total redundant proposal copies on a node attempting to broadcast its view of the state.

`Variables`

$$
\begin{align}
  &ps = \text{the preference score, where higher score indicates a higher chance of being selected} \\
  &v_{curr} = \text{the current version of the state, which is also is the total number of successful writes} \\
  &w_{recr} = \text{the weight applied to redundant proposal copies} \\ 
  &w_{RS} = \text{the weight applied to the reputation score}
\end{align}
$$

`Results`

$$
\begin{align}
  &RS = \text{the reputation score calculated for the node}
\end{align}
$$

Proposals will be considered successful based on which proposal has the highest preference score.

$$
\begin{align}
  &ps = w_{recr}\times{\text{number of redundant proposals}} +  w_{RS}\times{RS} \\
  &(w_{recr} + w_{RS})\le{1}
\end{align}
$$

The rejected state proposals will be retried with the next version of the state once an individual proposal has been accepted by the quorum and written to the state machine, where each failed attempt uses an exponential backoff strategy. 

On each successful write to the state machine from a particular node in the cluster, the node's reputation will increase logarithmically, with a base of $m$, the `propagation factor`. The node's score should increase logarithmically with the total number of direct interactions it has. Before applying the logarithm, the square root of the historical successes is taken.

$$
\begin{align}
  &RS = log_{m}(\sqrt{\text{total historical successes}})
\end{align}
$$

The logarithmic nature of reputation growth is meant to create a more democratic and fair system, where issues like reputation inflation can be mitigated since reputation will grow quickly as a node makes more successful writes once it joins the network, but will flatten out over time the more writes are made. This helps to ensure that nodes that are new to the network or that have low activity are not punished and are actually encouraged to make successful writes. This logarithmic reputation growth is meant to again mimic how reputation works in social settings, where unknown individuals can quickly gain reputation, but once more people learn about that individual, the reputation begins to flatten out until attrition is reached.

Since a node is essentially linked to a unique identifier, the unique id would be persisted even if a node goes offline, along with its historical successful writes to the state machine.


## Assumed Quorum through a Gossip Protocol and State Change Proposal Redundancy

### Model

A node can only make one proposal to state change per version and cannot then propose to change the global version tag until its proposal has either been successfully acknowledged by a quorum, or a different request updates the state + version and a the node receives a copy of the new data + version. The following algorithm outlines the state change proposal redundancy solution for updating global state in a leaderless distributed system. The algorithm itself does not handle conflict resolution, which is outlined in the `Reputation Based Conflict Resolution Model` section, but is meant to instead outline how quorum is achieved in a leaderless network through redundancy of state change proposals on a quorum of nodes. This is meant to mimic how information is spread amongst people in social environments. When an individual hears the same information from multiple unique sources, they are more likely to accept the information as fact the more sources they hear the information from. Data propagation is as follows:
```
  on received state change proposals:
    
    on first copy of the new state change proposal:
      verify the version of the incoming proposal:
        if the version tag is verified as being 1 + the current global version tag and there are no conflicts with other proposal copies for that version of the data from other nodes in the network:
          1.) the copied proposal is written to a tempory write ahead log or RAM
          2.) a subset of machines, determined by a propogation factor, is randomly selected and each machine is then propagated the change
          3.) the proposal is broadcasted to the machines in the subset
          4.) each machine in the subset is then marked as having received a copy of the data
        else: 
          1.) the proposal is rejected and the node waits for either a new request or a valid data copy
        
    for each redundant copy of the new state change proposal from a unique node after getting the initial copy:
      incremement the number of redundant copies from unique nodes in the network for the data copy:
        if the threshold for minimum number of redundant copies on the quorum of available machines is met:
          1.) stop broadcasting the data copy to new randomly selected nodes
          2.) apply the state change to the local state machine, incrementing the global version to match the version of the proposal
        else:
          1.) select a new random subset of machines to broadcast the proposal to, which is a subset of the propagation factor size from n - the previous subsets of nodes that the copy was forwarded to. Each node that the proposal is forwarded to must be unique as each node can only receive a copy once from each machine.
          2.) broadcast the proposal to the new subset
```

The above can be calculated as follows to approximate the threshold for minimum number of redundant copies an individual node in the cluster must receive to assume that the data copy proposal has been accepted by the minimum number of nodes to achieve quorum.

`Constants`

$$
\begin{align}
  &Q = \text{Quorum Factor: percentage of active nodes required to achieve quorum, set to 0.66}
\end{align}
$$

`Variables`

$$
\begin{align}
  &n = \text{Total Nodes: the number of active nodes in the network} \\ 
  &m = \text{Propagation Factor: the number of nodes to which a node can propagate data} \\
  &k = \text{Network Hops: the current number of network hops, or the level of propagation}
\end{align}
$$

`Result`

$$
\begin{align}
  &T_R = \text{Redundancy Threshold: the threshold hit or exceeded for the number of redundant copies received to determine quorum}
\end{align}
$$

`Assumptions`

$$
\begin{align*}
  &\bullet \text{No initial redundancy} \\
  &\bullet \text{Uniform probability of data distribution between nodes} \\
  &\bullet \text{A node can only send a redundant copy of data to a different node at most once}
\end{align*}
$$

The first step for approximating the threshold number of redundancies required for proposals to be agreed upon by the quorum is to determine the initial spread of the data throughout the network. The initial spread is used to determine the probability of a node in the network receiving a new version proposal within $k$ number of network hops. In a Gossip protocol, data propagation usually grows logarithmically, where the base of the logarithm is the propagation factor, where at each network hop, the total number of nodes reached increases by the propagation factor.

$$
\begin{align}
  &P_{\text{receives new message}} = \text{the binomial probability for calculating} \\
  &X = \text{the random variable representing the number of successes in $k$ network hops} \\
  &i_{\text{spread}} = \text{a range of network hops where $i_{\text{spead}}\leq{k}$} \\
  &p_{\text{receives new message}} = \text{the probability of a success on a single network hop} \\
  &q_{\text{receives new message}} = \text{the probability of a failure on a single network hop}
\end{align} \\
$$

$$
\begin{align*}
  &p_{\text{receives new message}} = \frac{m}{n}
\end{align*}
$$

$$
\begin{align}
  &\text{$q_{\text{receives new message}}$ is the inverse of $p_{\text{receives new message}}$} \\
  &q_{\text{receives new message}} = 1 - p_{\text{receives new message}}
\end{align} \\
$$

A binomial distribution can be used to calculate the probability of a node receiving a new proposal to update the state machine at least once in $k$ rounds, where $i_{\text{spread}}$.

$$
\begin{align}
  &P_{\text{receives new message}}(X = i_{\text{spread}}) = \binom{k}{i_{\text{spread}}} p_{\text{receives new message}}^{i_{\text{spread}}} q_{\text{receives new message}}^{k - i_{\text{spread}}}
\end{align}
$$

However, the the complementary probability $P(X = x)$ of a node not $x$ receiving the message at all in $k$ trials is modeled using the first trial $i=0$:

$$
\begin{align}
  &P_{\text{receives new message}}(X = 0) = \binom{k}{0}p^i_{\text{spread}} q^{k - 0} = q^k
\end{align}
$$

so the probability of receiving the message at least once is the inverse:

$$
\begin{align}
 &P_{\text{receives new message}} = 1 - P_{\text{receives new message}}(X = 0) = 1 - q_{\text{receives new message}}^k
\end{align}
$$

Total redundant messages can then be calculated once the initial spread is known.

$$
\begin{align}
  &i_{\text{recr}} = \text{a range of network hops where $2\leq{i}\leq{k}$} \\
  &R_{i_{\text{recr}}} = \text{the number of unique redundant messages a node has after $i_{\text{recr}}$ rounds} \\
  &p_{\text{recr}, i_{\text{recr}}} = \text{the probability a node receives a new unique redendant message in round network hop $i$} \\
  &q_{\text{recr}, i_{\text{recr}}} = \text{the probability of a failure on a single trial} \\
  &p_{\text{original}, i_{\text{recr}}} = \text{the probability a node has already received the original message}
\end{align}
$$

$$
\begin{align}
  &p_{\text{original}, i_{\text{recr}}} = P_{\text{received new message by $i_{\text{recr}}$}}
\end{align}
$$

$$
\begin{align}
  &\text{$q_{\text{recr}, i_{\text{recr}}}$ is the inverse of $p_{\text{recr}, i_{\text{recr}}}$} \\
  &q_{\text{recr}, i_{\text{recr}}} = 1 - p_{\text{recr}, i_{\text{recr}}}
\end{align}
$$

Since each node can only send a redundant message at most once to each node, initially each node has an $n-1$ chance to receive a message.

A geometric distribution function is created to calculate the probability that a node receives a new unique redundant message in the $i$-th round, with the following signature:

$$
\begin{align}
  &p_{\text{recr}, i_{\text{recr}}} = p_{\text{original}, i_{\text{recr}}}\times{\frac{n - 1 - R_{i_{\text{recr}}} - 1}{n - 1}}
\end{align}
$$

so expected value of $R$ at $k$ network hops, starting $2$ at would be:

$$
\begin{align}
  &R(k) = 1 - \sum{j = 2}^i_{\text{recr}} q_{\text{recr}, \text{j}}
\end{align}
$$

$R(k)$ will approach $1$ as the number of rounds increases.

To get the probability of one node getting exactly $T_R$ redundant messages after $i$ rounds, the following is calculated:

$$
\begin{align}
  &P_{\text{exact}, i_{\text{recr}}}(T_R) = R(i) - R(i - 1)
\end{align}
$$

A binomial distribution is then applied for $Q$.

$$
\begin{align}
  &P(X \ge{Q}) = \sum{x = Q\times{n}}^n \binom{n}{x}[P_{\text{exact}, \text{k}}(T_R)]^x [1 - P_{\text{exact}, \text{k}}(T_R)]^{n - x}
\end{align}
$$


## Global Version Tagging With Verifiable Delay Function

`Verifiable Delay Functions` are a relatively new cryptographic primitive that take a predictable amount of time to compute, even on multi-core systems, but are easily verifiable by viewers of the result. 

The global state machine shared by the network of nodes will utilize a global version tag, inspired by `MVCC` or `multi-version concurrency control` in multi-threaded database systems but adapted for distributed networks of nodes. The global tag will be a verified proof generated by the `vdf` for the `current version of the state + 1` to symbolize a verified version tag for a state change proposal.

Once a version has been accepted by the majority nodes in the cluster using both the `redundancy threshold` and `reputation based conflict resolution`, the volatile state change proposal is persisted to the state machine and the version is incremented. All new and retried requests will then create a new state proposal and then compute the vdf for the current version + 1, appending the verification to the state proposal before propagating it to new nodes in the network. Each node that receives the verification can check it to ensure that it is valid before deciding whether or not to update its own state proposal or to reject the proposal.

For computing the VDF, repeated squaring is chosen. No additional hashing is done on the VDF to apply randomness for security. However, a dynamic difficuly strategy is imployed to create a form of rate limiting. When a single node accumulates successful writes, the difficulty of solving the next vdf increases. This is done to give other nodes in the cluster a chance to submit proposals and introduces a built in rate limiting feature. Once a different node has written to the state machine, the node being rate limited has the total accumulated successful writes set to $0$.

`Variables`

$$
\begin{align}
  &v_{curr} = \text{the current version of the state} \\
  &v_{next} = \text{the proposed version tag for a state machine update} \\
  &s = \text{the total number of successful sequential writes} \\
  &N_{base} = \text{the base number of squaring operations to use} \\
  &N_{s} = \text{N dynamically adujsted after s sequential successful writes}
\end{align}
$$

`Constants`

$$
\begin{align}
  &g = \text{the base element in the group} \\
  &p = \text{a large prime number that defines the group}
\end{align}
$$

$v_{next}$ is computed as follows:

$$
\begin{align}
  &v_{next} = v_{current} + 1
\end{align}
$$

$N_{s}$ is calculated by taking $N_{base}$, where after each sequential successul write to the state machine causes a gradual growth in the difficulty of the vdf.

$$
\begin{align}
  &N_{s} = N_{base}\times{log_{10}{s}}
\end{align}
$$

Where the computation for the `vdf` is as follows. $s$, the total successful sequential writes, is used to increase the value of $N$, where $N$ is multiplied by the square root of the number of successfuly sequential writes, so that each attempt to write multiple proposals to the state machine will take increasingly longer, but not at such a rapid growth that clusters with smaller number of machines will be hurt at each iteration.

$$
\begin{align}
  &VDF(g, p, v_{next}, s) = g^{v_{next}^{2^{N_{s}}}}\mod{p}
\end{align}
$$

The value for $N$ is a major determinant in the total time to solve the VDF, so in this situation a smaller value should be selected where the total iterations results in a delay of between 50-100ms, so that responses to clients on write requests return within a reasonable wait time. 


## Heartbeating Through Propagation

Since there is no leader to enforce a strict heartbeat protocol, nodes in the network chatter with a subset of the network, the propogation factor defined above, to perform a liveness check. The nodes, in the liveness check, return a globally calculated reputation score for the current responder. If the node receives a response over a given number of retries from a node it is checking on, it will set the node to alive and update its local lookup of node ID to reputation score. The node broadcasts a payload containing both the the status and reputation score for each node it checked on to a new random subset of nodes not including any that were checked on. This concept is meant to mimic group connectivity in social situations. If a group of interconnected people have not heard from a member after a period of time, individuals are likely to perform a "liveness check" on the individual by attempting to interact with them.

The heartbeat is only for periods where requests are not being processed, as each proposal can act as a liveness check. Nodes receiving a proposal from a request will see that a node in the cluster is active. On state change, since the original proposer's unique ID is embedded inthe proposal, every node will increment the total count of successful writes the proposer has, where the reputation score will be recalculated for subsequent writes.

The liveness checks will propagate in a similar manner to all other requests using the gossip protocol, but the information is distributed with more randomness.


## Mitigating Gaming The System

One issue in distributed consensus models is gamification of the system, which is when a malicious actor attempts to manipulate the state to gain an advantage over other individuals in the network.

However, there are a couple key design choices in the proposal that help to mitigate these issues.

`Sybil Attacks and Spam Attacks`

Sybil Attacks are when bad actors attempt to create a large number of random node entities and use them to gain reputation in the system by spamming requests.

The use of a verifiable delay function helps to mitigate potential for spam and Sybil attacks on the network, making it computationally difficult to spam the network with malicious activity, as each subsequent request must compute the next verification for each state change proposal, where each must take at least a set amount of time to compute. The use of a `VDF` acts as both a method for maintaining a global ordering of events while also a deterrent for potential bad actors.

`Logarithmic Reputation Growth`

The use of a logarithmic growith rate on reputation helps mitigate reputation inflation, where a single node gains an overweighted reputation due to an enormous amount of successful writes. The logarithmic growth rate also helps give new nodes or nodes with low activity a higher chance of increasing their reputation score, creating a more democratic system.