<!-- markdownlint-disable MD024 -->

# Versioning

THORNode is following semantic version. MAJOR.MINOR.PATCH(0.77.1)

The MAJOR version currently is updated per soft-fork.

Minor version need to update when the network introduce some none backward compatible changes.

Patch version, is backward compatible, usually changes only in bifrost

## Prepare for release

1. Create a milestone using the release version (e.g. [Release-1.116.0](https://gitlab.com/thorchain/thornode/-/milestones/109))
2. Tag issues & PRs using the milestone, so we can identify which PR is on which version
3. PRs need to get approved by #thornode-team and #thorsec. Once approved, merge to `develop` branch
4. Once all PRs for a version have been merged, create a release branch from `develop` such as: `release-1.116.0`.

## Test release candidate locally

1. From your release branch, run `make build-mocknet`.
1. Create a mocknet cluster using `make reset-mocknet-cluster` (follow [README.md](../README.md)).
1. Sanity check the following features work:
   - [ ] Genesis node start up successfully
   - [ ] Bifrost startup correctly, and start to observe all chains
   - [ ] Create pools for BNB/BTC/BCH/LTC/ETH/USDT
   - [ ] Add liquidity to BNB/BTC/BCH/LTC/ETH/USDT pools
   - [ ] Bond new validator
   - [ ] Set version
   - [ ] Set node keys
   - [ ] Set IP Address
   - [ ] Churn successful, cluster grow from 1 genesis node to 4 nodes
   - [ ] Fund migration successfully
   - [ ] Some swaps, RUNE -> BTC, BTC -> BNB etc.
   - [ ] Mocknet grow from four nodes -> five nodes, which include keygen, migration
   - [ ] Node can leave
1. Identify unexpected log / behaviour, and investigate it.

## Release to stagenet

### Build stagenet

1. Merge release branch e.g. `release-1.116.0` branch -> `stagenet` branch. Once the changes are pushed, the stagenet image should be created automatically by pipeline.
1. Make sure `build-thornode` pipeline is successful, you should be able to see the docker image has been built and tagged successfully:

   ```logs
   Successfully built bbf5fe970c75
   stagenet: digest: sha256:8ec7a9c832ad13fc28d0af440b5cddfec8e21b4a311903ad92fe0cab0433faac
   stagenet-1: digest: sha256:8ec7a9c832ad13fc28d0af440b5cddfec8e21b4a311903ad92fe0cab0433faac
   stagenet-1.112: digest: sha256:8ec7a9c832ad13fc28d0af440b5cddfec8e21b4a311903ad92fe0cab0433faac
   stagenet-1.112.0: digest: sha256:8ec7a9c832ad13fc28d0af440b5cddfec8e21b4a311903ad92fe0cab0433faac
   ```

### Stagenet test plan

1. Create a test plan either mentally, on Discord or on Notion (e.g. [Stagenet 1.99 Test Plan](https://lively-router-cb2.notion.site/Stagenet-1-99-Testing-b5f40b8eac684612bd7fdad78d8e4ae9?pvs=4))
1. Consider what changes have shipped:
   - New features may require a dedicated test plan, as above for Savers. Consider the expected vs. actual result for querier endpoints before/after different transactions are made on-chain.
   - New chains will require the following process:
     1. Ensure all bifrost are running the latest version
     1. Ensure `loadchains.go` has successfully connected to the daemon.
     1. Ensure the new daemon is fully sync'd.
     1. Trigger a churn to create the asgard vault.
     1. Create a pool by sending L1 and RUNE to asgard.
     1. Once the pool is seeded with enough LP to pay outbound fees, churn the network again.
     1. Test inbound TX are observed and scheduled outbound TX are sent by doing swaps.
   - Changes to mimir should be set _after_ the new version is adopted. The same should be done for mainnet.
   - If protos are added or changed, it's a good idea to send messages on-chain with the next proto both before and after consensus is reached on the new version.
   - Include a stagenet store migration if one is to take place in mainnet as well. Be sure to check the pools endpoint (or other to-be-changed state) before and after the version increments.
   - Sanity different UIs (only Asgardex and THORSwap support stagenet at this time).
   - If making changes to chain clients or a chain daemon has been updated (`node-launcher/ci/images`): make sure the daemon is up-to-date and fully sync'd, then trigger both an outbound and observe an inbound for the affected chain(s).

These are just a few examples. Each release may contain unique functionality or infrastructure changes that need to be tested.

### Deploy stagenet

The `stagenet` maintainer does not need to keep the upstream `node-launcher` values for stagenet up-to-date. They are there as a reference. State can be kept locally.

The `node-launcher` repo will require the `stagenet` digest hash (e.g. `8ec7a9c832ad13fc28d0af440b5cddfec8e21b4a311903ad92fe0cab0433faac`). Get this from the `build-thornode` CI step above. Be sure you didn't accidentally copy the `mocknet` digest hash. Using the mocknet image in stagenet will cause a consnesus failure.

1. Apply the `stagenet` image(s) one-by-one.
1. Wait for them to fully initialize and rejoin consensus.
1. Run `make set-version`.
1. Repeat until all validators on latest version.
1. Check: `curl thornode:1317/thorchain/version`

### Validate stagenet

1. Conduct your [Stagenet Test Plan](#stagenet-test-plan).
1. Document any findings or issues in `#stagenet` on Discord.
1. Determine if any changes need to be made to the release candidate.

## Release to mainnet

### Build mainnet

1. Merge the release branch (e.g. `release-1.116.0` -> `mainnet`). Once the changes are pushed, the mainnet image should be created automatically by pipeline (e.g.: https://gitlab.com/thorchain/thornode/-/jobs/4682407839)
1. Make sure `build-thornode` pipeline is successful, you should be able to see the docker image has been built and tagged successfully:

   ```logs
   Successfully built d92da6e9c460
   mainnet-1.116.0: digest: sha256:58df167b2c515a0cf4f4093ca27ca49d85cd1201801f9baa3ffcdafaaa138bcb
   mainnet-1.116.0: digest: sha256:58df167b2c515a0cf4f4093ca27ca49d85cd1201801f9baa3ffcdafaaa138bcb
   mainnet-1.116.0: digest: sha256:58df167b2c515a0cf4f4093ca27ca49d85cd1201801f9baa3ffcdafaaa138bcb
   mainnet-1.116.0: digest: sha256:58df167b2c515a0cf4f4093ca27ca49d85cd1201801f9baa3ffcdafaaa138bcb
   ```

### Raise PR in node-launcher

1. Raise PR to release version to `node-launcher/thornode-stack/mainnet.yaml`, bumping the version and tag according to the last step. (e.g. https://gitlab.com/thorchain/devops/node-launcher/-/merge_requests/876/diffs#16eb49b6065b1a08dae8d22c10d771efcce894af_4_2)
2. Post the PR to #devops channel, and tag @thornode-team @thorsec @Nine Realms teams to approve. It will need at least 4 approvals.

### Release to mainnet

Pre-release check

1. Quickly go through all the PRs in the release.
1. Apply the latest changes to a standby node and monitor the following:
   1. THORNode pod didn't get into `CrashloopBackoff`
   2. Version has been set correctly
   3. Bifrost started correctly.

Release

1. Run the PR log script to collect all of the PRs tagged in this milestone (e.g. `scripts/pr-log.py Release-1.116.0`).
1. Create a tag for the release on the `develop` branch (e.g. https://gitlab.com/thorchain/thornode/-/tags/v1.116.0). Copy and paste the output from the script above into the description.
1. After the tag is created, go to the UI, click `Create Release`. Use the PR log for the description again.
1. Post release announcement in #thornode-mainnet. Use previous messages as a template. Be sure to update the version number and tag URL.
1. For mainnet release, post the release announcement in Telegram #THORNode Announcement
