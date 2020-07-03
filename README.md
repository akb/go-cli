This is yet another library for structuring command-line interfaces. It aims to
utilize the Go standard library as much as possible and provide a simpler
API than other options.

Goals:
- As lightweight as possible; should leverage the Go standard library for
    everything it can
- Should enforce consistency. The framework should encourage developers to
    provide a good user experience, and conform to platform expectations.
- Simple and easy to use. Developers should be able to pick it up with little
    or no documentation.

Features:
- Easy to create nested subcommands
- Enforced contextual help for every command, subcommand, and flag
- Enforced use of Go contexts for traceability
- Patterns for environment and flag parsing
- Assertions for writing tests
