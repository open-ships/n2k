# ADR 0002: Preserve raw payloads until fields change

Status: accepted

## Context

Generated schemas do not always describe every reserved bit, trailing byte, or
proprietary variation. Encoding only decoded fields can therefore change an
otherwise untouched message. Always replaying the original bytes, however,
silently discards caller mutations.

## Decision

After decode, retain an owned original payload and a canonical encoding of the
decoded fields. On encode, compare the current field encoding with that
canonical value. Return the original bytes only when fields are unchanged;
otherwise return the newly encoded fields.

## Consequences

Passive decode/encode forwarding is wire-faithful, including unknown reserved
content. Normal struct mutation behaves as Go callers expect. Code that wants
to alter only unknown bits must work at the raw frame layer.
