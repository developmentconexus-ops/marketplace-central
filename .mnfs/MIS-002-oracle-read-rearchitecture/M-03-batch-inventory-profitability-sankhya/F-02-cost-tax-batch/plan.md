# F-02 Cost/Tax Batch Plan

1. Add the internal-read batch capability types and a new Oracle adapter using
   `oraclebatch.Chunks`, one shared semaphore acquisition, safe error wrapping,
   cost/tax SQL, and capped sales-history peeking.
2. Extend the profitability internal-read adapter with optional batch methods;
   gather resolved IDs before item mapping and map nil facts explicitly.
3. Enforce the import ceiling before the first port call and map the transport
   error to the fixed 422 `limit_exceeded` shape.
4. Wire one root semaphore and the batch adapter additively, without touching
   `reader.go`, OpenAPI, or SDK runtime files.
5. Add focused tests for chunk counts, dedupe/partial failure, truncation,
   nil-fact incomputability, and the 200/201 limit behavior; run the registered
   targeted and full Go proofs and record evidence.
