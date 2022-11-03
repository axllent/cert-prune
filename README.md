# Cert Prune - delete obsolete Let's Encrypt certificates

This is a simple no-frills CLI utility to delete obsolete Let's Encrypt certificate files from your system. Every time certificates are registered or renewed, [certbot](https://certbot.eff.org/) generates new certificates in `/etc/letsencrypt`. It never deletes the old expired ones ([see GitHub issue](https://github.com/certbot/certbot/issues/4635)).

Whilst the physical storage of these certificates is not the issue (they do not take up much space), over time there can be literally tens of thousands of redundant files left within the folder structure.


## Installation

You can download a download a static binary from the [releases](https://github.com/axllent/cert-prune/releases/latest), or install from source using `go install https://github.com/axllent/cert-prune@latest`.



## Options

```shell
$ cert-prune -h
A utility to delete expired Let's Encrypt certficates.

All unused certificates, and (by default) all csrs & keys older than 60 days are deleted.

If no path is provided then /etc/letsencrypt is assumed.

Support:
  https://github.com/axllent/cert-prune

Usage:
  cert-prune [path] [flags]

Flags:
  -n, --nr-days int   Delete generation CSRs and Keys older than X days (default 60)
  -v, --verbose       Verbose logging
```


## Example usage

```shell
$ du -hs /etc/letsencrypt
191M	/etc/letsencrypt

$ cert-prune 
INFO Certs deleted:   27136                       
INFO CSRs  deleted:   8787                        
INFO Keys  deleted:   8787 

$ du -hs /etc/letsencrypt
7.9M	/etc/letsencrypt
```
