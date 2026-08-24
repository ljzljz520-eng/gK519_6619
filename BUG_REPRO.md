# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
ok  	bridge-trajectory/analysis	0.006s
?   	bridge-trajectory/cmd/bridgecli	[no test files]
ok  	bridge-trajectory/api	0.010s
ok  	bridge-trajectory/calc	0.001s
ok  	bridge-trajectory/domain	0.001s
ok  	bridge-trajectory/export	0.001s
ok  	bridge-trajectory/query	0.008s
ok  	bridge-trajectory/render	0.001s
--- FAIL: TestBridgeTrajectoriesKeepSnapshots (0.01s)
    trajectory_bug_test.go:40: snapshot point 1 was rewritten
FAIL
FAIL	bridge-trajectory/service	0.019s
ok  	bridge-trajectory/store	0.015s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/bridgecli): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/bridgecli): exit `0`
