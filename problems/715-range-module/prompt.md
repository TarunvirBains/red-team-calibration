# 715 Range Module

Task summary: Implement a structure that tracks covered half-open integer
ranges. It must support adding coverage, querying whether a range is fully
covered, and removing coverage.

Review pressure points:

- half-open interval semantics;
- adjacent range merge behavior;
- split behavior on removal;
- idempotent add/remove operations;
- empty or degenerate ranges if the implementation accepts them;
- invariant preservation across long operation sequences.

