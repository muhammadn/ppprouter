# ppprouter
Simple multi-homed PPP router

### Why?
All the solutions for Linux requires me to enable multipath networking (recompile) in the Linux kernel source.

However, my stock router hardware comes with a customised Linux kernel and the manufacturer's custom Linux
kernel source code is not updated to the latest.

### Note

* I wanted to make it work quickly, so i had to resort to use commands as mentioned below.
* This is a work in progress. It is not complete as it requires `ifmetric` and `curl`
* Will be replacing `ifmetric` and `curl` with native Go libraries (good for systems that does not have/or want to have those commands installed)
