# gopass
Inspired by [supergenpass-lib](https://github.com/chriszarate/supergenpass-lib) this is an attempt to make a (not yet complete) password generator based on master password and domain.

Usage
=====
    Usage of gopass:
      -additional-info="": Free text to add (e.g. index/timestamp if the previous password was compromized)
      -domain="": The domain for which this password is intended
      -log-domain=false: Whether to log the domain and the additional info for each generated password. Note that the password itself will NOT be stored!
      -master="": The master phrase to use for password generation. Required unless master-file is provided. Do NOT forget to escape any special characters contained in the master phrase (e.g. $, space etc).
      -master-file="": The path to a file, containing the master phrase. Required unless master is provided.
      -password-length=12: Define the length of the password. Default: 12
      -special-characters=true: Whether to add a known set of special characters to the password

Examples
========
* Generate password with file:


    $ gopass -master-file passphrase.file -domain google.com -log-domain
    Your password for google.com is: X68qP6hp@%.;

and the domains.log file contains

    Domain: [google.com], Special Characters: [true], AdditionalInfo: []

* Generate password with cmd parameter (warning - this will remain in your shell/cmd history!)


    $ gopass -master super-mega-secret-master-phrase -domain github.com -password-length 16 -additional-info rev3
    Your password for github.com is: 06r68L1RMlyN)(*$
