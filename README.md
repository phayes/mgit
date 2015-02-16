# mgit - the git multiplexer

Work on multiple git repositories at the same time.

### example usage

```bash
$ mgit clone git@github.com/phayes/* # check out all phayes repositories
$ mgit checkout testing # Check out testing branch on all repositories
```

### installing

```bash
$ wget "https://phayes.github.io/bin/current/mgit/linux/mgit.gz" #replace `linux` with `mac` for MacOSX version.
$ gunzip mgit.gz
$ sudo cp mgit /usr/bin     
$ sudo chmod a+x /usr/bin
$ echo "export GITHUB_API_TOKEN=mygithubapitoken" >> ~/.bash_profile && source ~/.profile # for github integration
```
