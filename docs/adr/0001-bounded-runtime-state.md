# ADR 0001: Bound runtime queues and reassembly state

Status: accepted

## Context

Marine buses can run indefinitely and can contain malformed, duplicated, or
bursty traffic. User consumers and callback providers may also stop making
progress. Unbounded queues turn either condition into process-wide memory
growth; blocking queues allow application behavior to stall bus citizenship.

## Decision

Write admission, live subscriptions, fast-packet sessions, gateway parser
bodies, and ISO transport sessions are bounded. Slow live subscribers fail
with `ErrReceiveOverflow`; saturated writes fail with `ErrWriteQueueFull`.
Fast-packet partial state has both a TTL and a maximum entry count. Protocol
processing never waits for user message consumption.

## Consequences

Applications must choose retry/drop policy for write saturation and restart a
subscription after inspecting overflow. Buffer sizes are configurable. Memory
and shutdown behavior remain predictable under hostile traffic.
