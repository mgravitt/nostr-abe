# Quick Start
```
go run main.go
```

# Attribute-based Encryption

Attribute-based Encryption (ABE) is a technique that allows messages to be decrypted based on a series of attributes of a specific user, e.g. if they meet certain criteria.

For example, if a business is making hiring plans, a specific document may be readable by directors in accounting, HR, or anyone with a C-level role. In this case, the policy would look like: 
```
(((level == directory AND (department == accounting) OR (department == HR)) OR level == C-level)
```

In the above example, however, a single company can control the attributes of all of the users.

# Multi-authority ABE 
Multi-authority ABE allows unrelated parties to designate or class users with certain attributes independently.

The repo has a sample application where authorities designate users by:
- Region: americas, europe, asia, australia, icelandic
- Job: plumber, electrician, scientist, programmer, nurse
- Tier: platinum, gold, silver, bronze

The hard-coded policy for encrypting the plain text is:
```
((region:americas AND job:scientist) OR (region:icelandic AND job:plumber)) OR (tier:platinum)
```

It doesn't really make it sense, but it just an example.

The program creates 10 random users and sees which ones are able to decrypt based on the policy.

```
User            :  Rudy
-- Region       :  region:americas
-- Job          :  job:scientist
-- Tier         :  tier:gold
Can decrypt?    :  true

User            :  Cicero
-- Region       :  region:europe
-- Job          :  job:scientist
-- Tier         :  tier:platinum
Can decrypt?    :  true

User            :  Peyton
-- Region       :  region:europe
-- Job          :  job:programmer
-- Tier         :  tier:silver
Can decrypt?    :  false
```

# Nostr Badge Integration

[NIP-58](https://github.com/nostr-protocol/nips/blob/master/58.md) providers support for badges, also described in [this article](https://thebitcoinmanual.com/articles/nostr-badges/).

Each badge is granted or designated by a particular badge owner. 

Multi-party ABE could be used to send messages to a surgical set of badge holders based on flexible policies. For example, a specific message may only be readable by accounts that are designated by badge owners as something like: 

```
((shit_poster AND influencer) OR (founder AND honor_badge) OR mega_donor))
``` 

I don't know if this would be useful, but the code works.