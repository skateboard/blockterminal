# Arguments

This will explain all the command line arguments you can use on the Blockterminal program.

# Command Line Arguments

| Argument | Type | Default | Description | Example |
|----------|------|---------|-------------|---------|
| `--node` | string | _(empty)_ | the node file to load on startup | `--node=path/to/node.json` |
| `--wallet` | string | _(empty)_ | The wallet file to load on startup | `--wallet=path/to/wallet` |
| `--tor` | bool | `false` | Use a tor proxy for all http connections | `--tor` |
| `--torPort` | int | 9050 | The wallet file to load on startup | `--torPort=9050` |

## Notes

- All arguments are optional and will use defaults if not specified
- The `--wallet` argument expects a path to a valid wallet file
- The `--node` argument expects a path to a valid node JSON file