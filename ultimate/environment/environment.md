# Golang Environment Variables and Commands

# Environment Variables
1. Basic Variables: 
    - GOROOT: Usually point to the golang SDK directory
    
    - GOPATH: 
    Places where to look for the Go code
    > **GOPATH** must be set to get, build and install 
    packages outside the standard Go tree.
    
    - GOBIN:
    The command `go install` build the binary executable file and put it to install directory
    . Tne install directory is controlled by **GOPATH** and **GOBIN** environment variables. 
    The rule of this reference as below:
    >> If GOBIN is set, binaries are installed to that directory. If GOPATH is set, binaries are installed to the bin subdirectory of the first directory in the GOPATH list. Otherwise, binaries are installed to the bin subdirectory of the default GOPATH ($HOME/go or %USERPROFILE%\go).
1. Advanced variable    
