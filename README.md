# hr-royalty-hlf

Solution repository for the **Music Royalty Distribution** Hyperledger Fabric
chaincode challenge (NPCI / HackerRank, Hard).

Standard Fabric **test-network** plus a chaincode skeleton at
[`chaincode/royalty.go`](chaincode/royalty.go). Cloned into the candidate's
environment by the HackerRank Setup Script (via [`setup.sh`](setup.sh)).

## Candidate task
1. Implement the functions in `chaincode/royalty.go`, including a weighted payout
   split that loses no paise to integer rounding.
2. Deploy: `cd test-network && ./network.sh deployCC -ccn royaltycc -ccp ../chaincode -ccl go`
3. Register trk1 with 50/30/20 splits, then distribute 100000 and 50000 paise.

---

Authored by **Dayal Mukati** — [dayalmukati.com](https://dayalmukati.com)
