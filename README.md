abireport
----------

[![Report](https://goreportcard.com/badge/github.com/clearlinux/abireport)](https://goreportcard.com/report/github.com/clearlinux/abireport) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)


A distro-agnostic tool for creating Application Binary Interface (`ABI`) reports from a set of binary packages. This tool is intended to assist packagers & developers of Linux distributions. It should be run at the end of the build, or even during the final stages of the build using the `scan-tree` root to be as unobtrusive as possible.

The tool, when invoked with `scan-packages`, will explode the packages requested, and perform an analysis of the binaries inside. From this, we create two files:

**symbols**

A `$soname`:`$symbol` mapping which describes the ABI of the package as a whole. This is sorted first by soname, second by symbol name:

        libgdk-3.so.0:gdk_x11_window_set_utf8_property
        libgdk-3.so.0:gdk_x11_xatom_to_atom
        libgdk-3.so.0:gdk_x11_xatom_to_atom_for_display
        libgtk-3.so.0:gtk_about_dialog_add_credit_section
        libgtk-3.so.0:gtk_about_dialog_get_artists
        libgtk-3.so.0:gtk_about_dialog_get_authors


**used_libs**

This file contains an alphabetically sorted list of binary dependencies for the package(s) as a whole. This is taken from the `DT_NEEDED` ELF tag. This helps to verify that a change to the package has really taken, such as using new ABI (soname version) or a new library appearing in the list due to enabling.

        libX11.so.6
        libXcomposite.so.1
        libXcursor.so.1
        libXdamage.so.1
        libXext.so.6


These files, when used with diff tools (i.e. `git diff`) make it very trivial to anticipite and deal with ABI breaks.

**Multiple architectures**

In many distributions, multilib or multiarch is employed. `abireport` will assign a unique suffix to each of these architectures to have a view on a per architecture basis. Currently, an `x86_64` file will have no suffix, and `x86` file will have the `32` suffix, etc. If you need a suffix added, please just open an issue.

Integrating
-----------

After your normal build routine, issue a call to `abireport scan-packages`. Note that this will decompress the binary packages in the directory first, so for maximum performance you should integrate into the build system itself. By default `abireport` is highly parallel so almost all of the time spent executing abireport is the decompression of packages.

Currently, `abireport` knows how to handle 3 package types:

 - `*.deb`
 - `*.rpm`
 - `*.eopkg`

More will be accepted by issue or pull request. In the event of a pull request, please ensure you run `make compliant` before sending, to ensure speedy integration of your code.

You may encapsulate `abireport` by using the `scan-tree` command after having decompressed the files yourself. If used against the true package build root, this zero-copy approach is significantly faster.

You should always store the current abireport results in your tree history, i.e. in git. Subsequent rebuilds will then automatically create a diff so that you can immediately see any actions that should be taken.


Implementation details
----------------------

The dependencies are evaluated using the `DT_NEEDED` tag, thus only direct dependencies are considered. Before emitting the report, `abireport` will check in the library names (`ET_DYN` files) to see if the name is provided. If so, it is omitted.

Symbols are only exported if they meet certain export criteria. That is, they must be an `ET_DYN` ELF with a valid `soname`, and living in a valid library directory. That means that `RPATH`-bound libraries are not exported.

This may affect some package which use a private RPATH'd library. From the viewpoint of `abireport`, such private libraries do not constitute a true ABI, given that many distributions are opposed to the use of `RPATH`. In effect, these are actually plugins (unversioned libraries).

License
-------

`Apache-2.0`

Copyright © 2016-2017 Intel Corporation
