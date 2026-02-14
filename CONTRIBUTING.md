# Contributing

We welcome contributions! Whether it's bug fixes, new mappings, documentation improvements, or test cases — all help is appreciated.

## How to contribute

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-change`)
3. Make your changes
4. Run tests (`go test -v ./...`) and vet (`go vet ./...`)
5. Commit with a clear message
6. Open a pull request

### Adding new mappings

To add mappings for a new framework or vendor:

1. Add entries to `processor/genainormprocessor/defaults.go`
2. Update `docs/compatibility-matrix.md` with the new rows
3. Add a vendor mapping pack YAML in `mappings/` if appropriate
4. Add tests covering the new mappings

## Contributor License Agreement (CLA)

**By submitting a pull request, you agree to the following:**

You grant Nostalgic Skin Co. (the project maintainer) a perpetual, worldwide, non-exclusive, royalty-free, irrevocable license to use, reproduce, modify, distribute, sublicense, and relicense your contributions under any license, including proprietary licenses.

You retain copyright ownership of your contributions. This CLA only grants the project maintainer the right to use your contributions in the project under current or future licenses.

**Why?** This project uses dual licensing (AGPL + Commercial). Without a CLA, we cannot offer commercial licenses to enterprise users, which funds continued development of the open-source project.

This is standard practice for dual-licensed open-source projects (used by MongoDB, Elasticsearch, Grafana, and many others).

## Code style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Add tests for new functionality
- Keep PRs focused — one feature or fix per PR

## Reporting issues

Open a GitHub issue with:
- What you expected to happen
- What actually happened
- Steps to reproduce
- Your collector version and config (redact any secrets)

## Questions?

Open a discussion or email jason.j.shotwell@gmail.com.
