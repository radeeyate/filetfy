# filetfy

![hackatime stats](https://hackatime-badge.hackclub.com/U06V886HQLS/filetfy)

easy & simple notifications for when files are created or deleted.

## how to use

1. install via `go install github.com/radeeyate/filetfy@v1.0.1`
1. create a [`.env` file](./.env.example) and modify to your hearts desire.
1. run `filetfy` to start watching for files

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
1. create a new branch
1. make your changes
1. create a pull request
1. thanks!
