# gopass
Inspired by [supergenpass-lib](https://github.com/chriszarate/supergenpass-lib) and [vpass](https://github.com/vladstudio/vpass2) this is an attempt to make a password generator based on master password and domain.

What it is
===========
This tool generates per domain passwords based on a master secret and some custom properties, controlled by command line arguments. This way users can have different passwords for each site without relying on a online service or password manager (which can be out of sync) or any other tool. All they need is to have the tool locally (or be able to download and build it) and remember their master secret.

Usage
=====
    Usage of gopass:
      -additional-info="": Free text to add (e.g. index/timestamp/username if the previous password was compromized)
      -domain="": The domain for which this password is intended
      -log-domain=false: Whether to log the parameters that were used for password generation to a file. Note that the password itself will NOT be stored!
      -master="": The master phrase to use for password generation. Do NOT forget to escape any special characters contained in the master phrase (e.g. $, space etc).
      -master-file="": The path to a file, containing the master phrase.
      -password-length=12: Define the length of the password.
      -special-characters=true: Whether to add a known set of special characters to the password


If neither master nor master-file is provided we try to find _~/.gopass/master_ file and load the master password from there.

Examples
========
* Generate password with file:
```
$ gopass -master-file passphrase.file -domain google.com -log-domain
Your password for google.com is: X68qP6hp@%.;
```
and the _domains.log_ file contains
```
Domain: [google.com], Special Characters: [true], AdditionalInfo: []
```

* Generate password with cmd parameter (warning - this will remain in your shell/cmd history!):
```
$ gopass -master super-mega-secret-master-phrase -domain github.com -password-length 16
Your password for github.com is: 06r68L1RMlyN)(*$
```
* If your password is compromised - change it completely:
```
$ gopass -master super-mega-secret-master-phrase -domain github.com -password-length 16 -additional-info rev2
Your password for github.com is: 41sLFIlPOFHI[ -=
```

Problems
========
I have not seen any functional problems so far, but please send feedback/issues if you find any :)

However I took a shortcut by defining a predefined set of special character groups to make sure that the password contains such symbols. I intent to improve this at some point but for now I don't have a clear idea how to make it "random enough".
