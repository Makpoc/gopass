[![Build Status](https://travis-ci.org/Makpoc/gopass.svg?branch=master)](https://travis-ci.org/Makpoc/gopass) [![code-coverage](http://gocover.io/_badge/github.com/Makpoc/gopass/generator)](http://gocover.io/github.com/Makpoc/gopass/generator) [![go-doc](https://godoc.org/github.com/Makpoc/gopass/generator?status.svg)](https://godoc.org/github.com/Makpoc/gopass/generator)
# gopass

Inspired by [supergenpass-lib](https://github.com/chriszarate/supergenpass-lib) and [vpass](https://github.com/vladstudio/vpass2) this is an attempt to make a password generator based on master password and domain.

# What is it

This tool generates per domain passwords based on a master secret and some custom properties, controlled by command line arguments. This way users can have different passwords for each site without relying on a online service or password manager (which can be out of sync) or any other tool. All they need is to have the tool locally (or be able to download and build it) and remember their master secret.

# Installation

Execute 
```
go get github.com/makpoc/gopass/ui/gopass-cmd
``` 
to get the commandline interface or 
```
go get github.com/makpoc/gopass/ui/gopass-web
```
for the WebUI.

# Usage
```
Usage of gopass-cmd:
  -additional-info="": Free text to add (e.g. index/timestamp/username if the previous password was compromized)
  -domain="": The domain for which this password is intended
  -log-info=false: Whether to log the parameters that were used for password generation to a file. Note that the password itself will NOT be stored!
  -master="": The master phrase to use for password generation. Do NOT forget to escape any special characters contained in the master phrase (e.g. $, space etc).
  -master-file="": The path to a file, containing the master phrase.
  -password-length=12: Define the length of the password.
  -special-characters=true: Whether to add a known set of special characters to the password
```

# Examples

* Generate password with file:
```
$ gopass-cmd -master-file passphrase.file -domain google.com -log-domain
Your password for google.com is: X68qP6hp@%.;
```
and the _domains.log_ file contains
```
Domain: [google.com], Special Characters: [true], AdditionalInfo: []
```

* Generate password with cmd parameter (warning - this will remain in your shell/cmd history!):
```
$ gopass-cmd -master super-mega-secret-master-phrase -domain github.com -password-length 16
Your password for github.com is: 06r68L1RMlyN)(*$
```
* If your password is compromised - change it completely:
```
$ gopass-cmd -master super-mega-secret-master-phrase -domain github.com -password-length 16 -additional-info rev2
Your password for github.com is: 41sLFIlPOFHI[ -=
```

# Configuration files

:information_source: For the command line tool only:

If neither ```master``` nor ```master-file``` is provided on the command line the program will make one last attempt to get a master password from ```$GOPASS_HOME/master```. If you use this file make sure it's secured, because it needs to contain your master password in plaintext. On unix - make sure you set very restrictive permissions to it (e.g. ```0600```). Encrypting your HDD is also something to consider (not just because of this program, but in general).

The default location for the domains log file is also under ```$GOPASS_HOME/``` with name ```domains.log``` and is currently not configurable.

If ```$GOPASS_HOME``` is not set the config folder is ```$HOME/.gopass```.

# Problems

I have not seen any functional problems so far, but please send feedback/issues if you find any :)
However I took a shortcut by defining a predefined set of special character groups to make sure that the password contains such symbols. I intent to improve this at some point but for now I don't have a clear idea how to make it "random enough".


# TODOs

* Configure the server to run over SSL - who would want to send their password over plain HTTP. :)
```
go run $GOROOT/src/crypto/tls/generate_cert.go --host="localhost"
```
* Mobile UI
