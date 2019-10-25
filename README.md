# Go library and CLI for Glass Factory

This is a CLI for [Glass Factory] resource management tool.

Uses APIs implemented in [github.com/markosamuli/glassfactory] package.

[Glass Factory]: https://glassfactory.io/
[github.com/markosamuli/glassfactory]: https://github.com/markosamuli/glassfactory

## Usage

### Build

Build the binary:

```bash
go build github.com/markosamuli/glassfactory-cli
```

### Configuration

Create a configuration file:

```bash
glassfactory config
```

### Reports

Generate report for current fiscal year:

```bash
glassfactory report fy
```

Generate monthly reports for current calendar year:

```bash
glassfactory report monthly
```

## License

* [MIT License](LICENSE)

## Author Information

* [@markosamuli](https://github.com/markosamuli)