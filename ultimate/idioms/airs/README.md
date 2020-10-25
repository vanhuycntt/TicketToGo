
## Fundamentals

- In golang interface implementation is implicitly, it's different with other languages that interface implementation is explicitly.

- The bigger the interface, the weaker the abstraction, ”interfaces with only one or two methods are common in Go code”

    > The io.Reader and io.Writer interfaces are the usual examples for powerful interfaces.
    
    ```
    type Reader interface {
        Read(p []byte) (n int, err error)
    }
    ```
## Idioms

### **Consumer**
   > Accept interface and return concrete type

- The consumer should define the interface. If you’re defining an interface and an implementation in the same package, you may be doing it wrong(**Dave Cheney on Dec 18 2017**)
    > An example is the io.Copy function. It accepts both the Writer and Reader interfaces as arguments defined in the same package.
      
    ```
    func Copy(dst Writer, src Reader) (written int64, err error)
    ```

### **Producer**
    
- If a type exists only to implement an interface and will never have exported methods beyond that interface, there is no need to export the type itself

    > An example is the rand.Source interface which is returned by rand.NewSource. The underlying struct rngSource
     within the constructor only exports the methods needed for the Source and Source64 interfaces so the type itself is not exposed.

- The benefit of returning interface over concrete type:
    - Returning an interface allows you to have functions that can return multiple concrete types. For example, 
    the aes.NewCipher constructor returns a cipher.Block interface. If you look within the constructor, 
    you can see that two different structs are returned.
- The mechanics of this pattern is:
    - The interface returned needs to be small so that there can be multiple implementations.
    - Hold off on returning an interface until you have multiple types in your package implementing only the interface. Multiple types with the same behavior signature gives confidence that you’ve the right abstraction.
- Consider creating a separate interfaces-only package for namespacing and standardization
    > An example is the hash.Hash interface that’s implemented by the packages under the subdirectories of hash/ such as hash/crc32 and hash/adler32. 
    The hash package only exposes interfaces.
    
    Three are two points of benefits in this way
    
    - A better namespace for the interfaces. hash.Hash is easier to understand than adler32.Hash.
    - Standardizing how to implement a functionality. A separate package with only interfaces hints that hash functions 
    should have the methods required by the hash.Hash interface.
    
- Finally, private interfaces don’t have to deal with these considerations as they’re not exposed.    
    > We can have larger interfaces such as gobType from the encoding/gob package without worrying about its contents.
     Interfaces can be duplicated across packages, such as the timeout interface that exists in both the os and net packages, without thinking about placing them in a separate location.