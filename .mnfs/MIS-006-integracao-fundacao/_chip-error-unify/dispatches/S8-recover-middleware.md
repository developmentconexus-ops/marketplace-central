# S8 — Central panic recover (VC-3)

Blocks 1, 2, 3 and 5 are in `_common-blocks.md` in this same directory. **Read that file
first, in full.** It binds you exactly as if it were inline here. Block 4 (per-module
whole-body pin) does not apply to this slice; your tests are specified below instead.

No other slice is running right now, so `go build ./...` failing is YOUR problem, not
someone else's.

## goal

Today a panic in any handler produces whatever Go's default `net/http` server does: the
connection is closed with no body at all. There is no recover middleware anywhere in the
tree — the root chain is exactly one layer. VC-3 requires a panic to become a normal
`500` envelope, with the panic's own message never reaching the client.

## Where it goes

New function in the EXISTING package `apps/server_core/internal/platform/apierror`, in a
new file `recover.go`:

```go
func Recover(next http.Handler) http.Handler
```

It lives in `apierror` rather than in `httpx` for a concrete reason: it must ANSWER with
the envelope, and `httpx` cannot import `apierror` (the direction is `apierror -> httpx`;
reversing it is an import cycle, which is why `httpx/json.go` and
`httpx/route_deadline.go` each had to hardcode a constant literal instead). Putting it in
`apierror` lets it call `apierror.Write` and keeps the 500 body defined in exactly one
place. Do not add a third literal copy of the envelope anywhere.

## Behaviour, precisely

On a recovered panic:
- status `500`, code `internal_error`, `details` `{}`.
- message: a fixed human pt-BR string. **The recovered value must NOT appear in the
  body** — not the panic message, not a stack trace, not a type name. This is the same
  rule as `httpx.WriteJSON`'s encode-failure fallback; read how that one is worded and
  stay consistent with it.
- The panic value AND a stack trace go to the log via `log/slog` at error level. Match the
  logging idiom already in the tree (`catalog/transport/http_handler.go` logs with
  `slog.Error`). The log is where the detail belongs; the body is where it must not be.
- `http.ErrAbortHandler` must be re-panicked, not swallowed. `net/http` uses it as a
  deliberate abort signal and suppressing it turns an intentional connection abort into a
  bogus 500. Check `errors.Is` against it before handling.

**The hard part — a panic AFTER the handler already started writing.** If the handler has
already called `WriteHeader` or written body bytes, you CANNOT answer 500: the status is
on the wire and a second `WriteHeader` produces Go's "superfluous WriteHeader" warning and
a corrupt body glued onto a partial one. In that case: log it (at error level, saying
explicitly that the response was already committed) and return WITHOUT writing. Do not
invent a body, do not pretend it succeeded.

Detecting that requires wrapping the `http.ResponseWriter` to record whether it was
written to. That is a real second consumer of nothing, so keep it unexported and minimal —
a struct with the embedded `http.ResponseWriter` and one bool, overriding `WriteHeader`
and `Write`. `internal/platform/httpx/route_deadline.go:156` already has a wrapper of this
family; read it, cite it in your report, and follow its shape rather than inventing a
different one. Do NOT add `Flush`/`Hijack`/`Push` pass-throughs unless a test in this
repo actually needs them — that is speculative surface.

## Where it attaches

`apps/server_core/internal/composition/root.go`, the final return (currently around
`:862`):

```go
return &RootRuntime{Handler: httpx.CORSMiddleware(mux), ...}
```

becomes `httpx.CORSMiddleware(apierror.Recover(mux))`.

**The order is not arbitrary and must be exactly this.** `CORSMiddleware`
(`httpx/router.go:18`) sets its headers on the way IN, before calling `next`. With CORS
outermost, a panicking request still carries the CORS headers, so a browser can actually
READ the 500 instead of reporting an opaque CORS failure — which is the difference between
the operator seeing an error toast and seeing a blank screen. Recover inside CORS also
means recover cannot swallow the `OPTIONS` preflight short-circuit. Say in your report
that you checked this ordering and why.

## Tests — file `apps/server_core/internal/platform/apierror/recover_test.go`

1. **Panic becomes the envelope.** A handler that panics with a value containing a
   recognisable secret-ish token (e.g. `panic("boom-SENTINEL-9f3a")`), driven through
   `Recover` with `httptest`. Assert: status 500; the COMPLETE body equals
   `{"error":{"code":"internal_error","message":"<your message>","details":{}}}` compared
   with `trimJSON` on BOTH sides (the package's test file already has `trimJSON` — reuse
   it, do not define a second one); and **assert the body does NOT contain
   `SENTINEL`**. Asserting absence is the point: a body that merely looks right could
   still be carrying the panic text appended somewhere.
2. **Negative control — a sane route is untouched.** A handler that returns `200` with a
   normal JSON body, through the SAME `Recover` wrapper. Assert the status and body come
   through byte-identical, and that the response does NOT have the panic shape. Without
   this, a `Recover` that answered 500 to everything would pass test 1.
3. **Already-committed response.** A handler that writes `200` and some body bytes and
   THEN panics. Assert the status stays `200`, that the partial body is unchanged, and
   that the envelope was NOT appended. This is the case that would otherwise corrupt a
   real streaming response.
4. **`http.ErrAbortHandler` is re-panicked.** A handler that panics with
   `http.ErrAbortHandler`; assert the panic propagates out of `Recover` (use
   `defer func(){ recover() }()` in the test to catch it) rather than becoming a 500.
5. **The panic detail reaches the log.** Capture `slog` output with a
   `slog.New(slog.NewTextHandler(&buf, nil))` and assert the sentinel token IS present
   there. Tests 1 and 5 together are the actual claim: the detail moved from the body to
   the log, rather than simply being destroyed.

Prove each test real per Block 5 step 2 — a captured NAMED `--- FAIL:` before it passes.

## Also

`root.go` has no test file of its own for this. Do NOT create an integration test that
boots anything (booting a server is forbidden for you). The wiring change is proven by
`go build ./...` plus the unit tests above; say exactly that in your report rather than
claiming the wiring is verified end-to-end.

## write_set

- `apps/server_core/internal/platform/apierror/recover.go`
- `apps/server_core/internal/platform/apierror/recover_test.go`
- `apps/server_core/internal/composition/root.go` (the single return line only)

## Your packages for the test command

`./internal/platform/apierror/... ./internal/composition/...`
Then also run the full `go build ./...` and `go vet ./...`.

**open_questions**: none. If you find one, stop and report rather than deciding.
