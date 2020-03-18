% ABIREPORT(1)

## NAME

`abireport -- Generate ELF ABI reports`

## SYNOPSIS

`abireport [subcommand] <flags>`


## DESCRIPTION

`abireport(1)` is tool to generate an ABI (Application Binary Interface) report
from ELF binaries. In either operational mode, it will scan all dynamic libraries
and executables that it encounteres, and emitting files to describe the ABI.

These file names may be different if you pass the `-p`,`--prefix` option. Also
note that for a non `x86_64` architecture, a unique suffix will be used for the
report file to enable tracking of multilib/multiarch configurations. The default
filename suffix for `x86` is `32`.

 * `symbols`

    A file containing a `$soname`:`$symbol` mapping for the entire scanned set.
    This file is sorted first by `$soname`, i.e. `libz.so.1`, and all symbols
    for that library are listed, sorted by alphabetical order.

    Running the tool will automatically truncate this file if it exists prior
    to creating a new report.

 * `used_libs`

    A file containing an alphabetically sorted list of dependencies required
    for the given data set. Any symbol that one of the files depends on that
    can be satisfied within the given data, set, even if it is not **exported**
    as a symbol, will omit the `$soname` dependency from the list.

    Used libs are determined by the `DT_NEEDED` section of the ELF file.

    Running the tool will automatically truncate this file if it exists prior
    to creating a new report.

## OPTIONS

These options apply to all subcommands within `abireport(1)`.

 * `-p`, `--prefix`

   Set the prefix used in filenames created by `abireport(1)`.  This may be
   used to assist integration by providing more unique filenames.

 * `-D`, `--output-dir`

   Set the output directory for files created by `abireport(1)`. This is
   limited strictly to the report files.

   This option defaults to the current working directory (`.`).

 * `-h`, `--help`

   Help provides an explanation for any command or subcommand. Without any
   specified subcommands it will list the main subcommands for the application.


## SUBCOMMANDS

Subcommands are mutually exclusive, you may only use one at a time. Note that
all subcommands respect the global flags.

### scan-tree [root]

Generate a report from the contents of the indicated root directory.
It is assumed that this is the true root of that filesystem, containing a
legitimate hierarchy, i.e. `/usr/lib`, etc.


### scan-packages [package] [package]

Generate a report from the contents of the given packages. They
will all be extracted into a temporary directory and analysed in
one go. You may pass multiple file names and directories here.

When using directories, `abireport(1)` will not recurse, it
will only look for a glob pattern of **supported** package types:

 * `*.rpm` - requires `rpm2cpio` and `cpio` on the host
 * `*.deb` - requires `dpkg` on the host
 * `*.eokpg` - requires `uneopkg` on the host.


### version

    Print the version and copyright notice of `abireport(1)` and exit.


## EXIT STATUS

On success, 0 is returned. A non-zero return code signals a failure.


## COPYRIGHT

 * Copyright © 2016 Intel Corporation, License: CC-BY-SA-3.0


## NOTES

Creative Commons Attribution-ShareAlike 3.0 Unported

 * http://creativecommons.org/licenses/by-sa/3.0/
