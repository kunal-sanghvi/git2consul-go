# Git2Consul-GO

[Git2Consul](https://github.com/breser/git2consul) is an amazing tool used by many organisations to mirror a git repository content to [Consul](http://www.consul.io/) KVs

## Requirements
* [Golang](https://golang.org/doc/install) 1.11 or above

## Limitations
* The current version (1.0.0) does not support cloning by http when 2 factor auth is enabled
* Config file should be in JSON format only

## Commands
* `--config FILE, -c FILE `  Specify config FILE, defaults to config.json in current directory
* `--pemfile FILE, -p FILE`  Specify pem FILE, defaults to ~/.ssh/id_rsa in home directory

## Sample config file
```javascript
{
  "host": "<YOUR CONSUL URL>",
  "repos": [
    {
      "name": "<Repo name>",
      "url": "git@github.com:kunal-sanghvi/repo1Name.git",
      "branches": [
	"master",
        "feature-branch1",
	"develop"
      ]
    },
    {
      "name": "<Another repo Name>",
      "url": "git@github.com:kunal-sanghvi/repo2Name.git",
      "branches": [
        "master",
	"staging"
      ]
    }
  ]
}
```

> Made in :heart: with OpenSource