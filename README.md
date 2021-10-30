# DEPRECATED Go library and CLI for Glass Factory

This library is no longer maintained as I don't have access to Glass Factory
anymore.

![Test](https://github.com/markosamuli/glassfactory/workflows/Test/badge.svg)
![Lint](https://github.com/markosamuli/glassfactory/workflows/Lint/badge.svg)

This is a CLI for [Glass Factory][glassfactory] resource management tool that
I've created for learning purposes.

This is not an official tool supported by the [Glass Factory][glassfactory]
team and it might now work with your organisations's configuration.

You might also need certain permissions to use these APIs. I don't provide any
documentation or support for configuring your accounts.

[glassfactory]: https://glassfactory.io/

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
