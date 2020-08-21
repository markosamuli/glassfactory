# Go library and CLI for Glass Factory

This is a CLI for [Glass Factory] resource management tool.

[Glass Factory]: https://glassfactory.io/
[github.com/markosamuli/glassfactory]: https://github.com/markosamuli/glassfactory

## Usage

### Build

Build the binary:

```bash
make build
```

Install the binary:

```bash
make install
```

### Configuration

Authenticate and create a configuration file:

```bash
glassfactory auth login
```

### Reports

Generate report for the current fiscal year:

```bash
glassfactory report fy
```

Generate monthly reports for the current calendar year:

```bash
glassfactory report monthly
```

## License

[MIT License](LICENSE)

## Author Information

[markosamuli](https://github.com/markosamuli)
