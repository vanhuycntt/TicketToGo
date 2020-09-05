# Environment Variables, Commands, Package in Go

# Environment Variables
1. Basic Variables: 
    - GOROOT: The **GOROOT/src** is where the standard go package located, this directory is also the place where 
    the `import package` is looked for firstly
    
    - GOPATH: 
    Places where to look for the Go code
    > **GOPATH** must be set to get, build and install 
    packages outside the standard Go tree.
    
    - GOBIN:
    The command `go install` build the binary executable file and put it to install directory
    . Tne install directory is controlled by **GOPATH** and **GOBIN** environment variables. 
    The rule that install directory belongs to as below:
    > If **GOBIN** is set, binaries are installed to that directory. If **GOPATH** is set, binaries are installed to 
    the bin subdirectory of the first directory in the **GOPATH** list. Otherwise, binaries are installed to the bin subdirectory of the default GOPATH ($HOME/go or %USERPROFILE%\go).

1. Advanced variables

# Commands

1. go install <package-name>
    - Three are two package types: an executable package and a utility package. 
        + An executable package is the main application. 
        + Utility package is not self-executable, instead, it enhances the functionality of an executable package by providing utility functions 
        and other important assets
    - `go install` compile a package and create a binary executable file or package archive file.
    The package archive file is created to avoid the compilation of the package every single time it is imported in the 
    program. 
        + The `go install` command pre-compiles a package and `Go` refers to .a files.
        + Using nested packages path on the command to create archive file for nested package.
    
1. go get <url_path_dependency>
    The **GO11MODULE** environment variable is added in from version **go.1.11**, and is a preventive measure to deprecate 
    **GOPATH** in the near future.
    Depending on how this command is called from the folder where the source code located.
    When GO111MODULE=off, the go get command simply clones the package from the master repository of the package and puts 
    it inside $GOPATH/src. If **GO111MODULE**=on, then **go get** command puts the package code inside a separate
    directory and it supports versioning
    - `go get -u ` : update all the modules with the latest **minor** or **patch** version of the given **major** version
    - `go get -u=patch`: update only patch version even though **minor** version is available.
    - `go get module@version`: update to precise version of a dependency module
       
  

    
