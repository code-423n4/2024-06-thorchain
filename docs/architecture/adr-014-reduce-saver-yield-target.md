# ADR 14: Reduce Saver Yield Synth Target to Match POL Target

## Changelog

- 04/04/2024: Created

## Status

Proposed

## Context

The following are the relevant mimir values at the time of writing:

```text
MaxSynthsForSaversYield=6000
POLTargetSynthPerPoolDepth=3000
POLBuffer=2000
```

These values result in savers receiving positive yield while POL is simultaneously depositing into the pool to lower the synth utilization (when it is over 50%). Suggesting setting the following:

```text
MaxSynthsForSaversYield=5000
```

While more complicated discussions around saver ejections and synth backstop are pending, this is a simple change to enact a small reduction in synth utilization over time (even if savers do not respond to the reduced yield). Additionally, it avoids the current state where a small share of the POL in effect subsidizes the saver yield (increasing synth utilization), while it is trying to reduce it.

## Decision

TBD

## Consequences

Savers will not receive yield on pools where synth utilization is over 50% (above which POL currently deposits).
