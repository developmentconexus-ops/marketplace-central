# F-01 stock-batch implementation plan

1. Add the shared chunker/semaphore helper and focused unit tests.
2. Add the internal-read `StockFact` and batch port, then implement the new
   Oracle adapter and fake-driver proof for chunking, deduplication, failure
   atomicity, and semaphore concurrency.
3. Convert stock-risk application flow to collect IDs and make one batch call;
   add missing-stock quality output and service tests.
4. Wire one root semaphore and the batch adapter, preserving the existing
   interactive reader path.
5. Run targeted proof, full `go test ./...`, record exact output in
   `validation.md`, inspect the diff, and create one intentional commit.
