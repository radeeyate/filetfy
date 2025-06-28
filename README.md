# filetfy

easy & simple notifications for when a file is created, modified (soon) or deleted.

## how to use

1. clone the repo
2. copy `.env.example` to `.env` and fill out the variables
3. compile via `go build .`
4. run `./filetfy` (or equivelant command, such as `filetfy.exe` on windows) to start watching for files

## configuration

| variable | description |
| --- | --- |
| `NTFY_TOPIC` | the ntfy topic to send notifications to. could be `joes-tacos` or whatever you want. |
| `NTFY_SERVER` | the ntfy server to send notifications to. default is `https://ntfy.sh/`. please include protocl & the trailing slash. |
| `CRON_STRING` | how often to check for changes, should work with an cron string. must be set. |
| `DIRECTORIES` | comma separated list of directories to watch. must be set. |
| `CREATE_INDICATOR` | the indicator to use for created files. default is `+`. |
| `DELETE_INDICATOR` | the indicator to use for deleted files. default is `-`. |

## contributing

1. fork the repo
2. create a new branch
3. make your changes
4. create a pull request
5. thanks!
