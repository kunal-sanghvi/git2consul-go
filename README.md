# Git2Consul GO
    - Populate config from your git repo to consul host
    - Only CLI interface supported
    
# Pre-requisites
  - Make sure you have go (v1.11 and above) installed
  - Clone this repo
  - Rename `sampleConfig.json` to `config.json`
  - Fill in the details in your `config.json`

# CLI Parameters

| Flags | Default | Description | 
| --- | --- | --- |
| c, config | ${PWD}/config.json | Path to your config.json |
| p, pemfile | ~/.ssh/id_rsa     | Path to your private ssh key |

# Setup
```bash
$ go build git2consul.go
$ ./git2consul -c <config.json path> -p <private ssh key path>
```

### To-dos

 - Support for ssh keys with passphrase
