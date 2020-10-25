
### Fundamentals
- Golang is statically typed language, so every variables get a static type when declared. In the following code, 
for instance, is an illustration about the interface type for this statement
  ```go
    // r has static type `io.Reader`, and the variable type of r can hold any value whose type 
    // has a `Read` method
    var r io.Reader
    r = os.Stdin
    r = bufio.NewReader(r)
    r = new(bytes.Buffer)
    // and so on
  ```
  
- 3 laws of reflection
  - Reflection goes from interface value to reflection object
    > There are two types we need to know about in package reflect: `Type` and `Value`. Those two types give access to 
    the contents of an interface variable, and two simple functions, called reflect.TypeOf and reflect.ValueOf, 
    retrieve reflect.Type and reflect.Value pieces out of an interface value
    ```go
       var x float64 = 3.4
       fmt.Println("type:", reflect.TypeOf(x))
        // print
        // type: float64
    ```
    ```go
       var x float64 = 3.4
       fmt.Println("value:", reflect.ValueOf(x).String()) 
       // print
       // value: <float64 Value> 
    ```
    >  a `Kind` method that returns a constant indicating what sort of item is stored: Uint, Float64, Slice, and so on
    
    ```go
       var x float64 = 3.4
       v := reflect.ValueOf(x)
       fmt.Println("type:", v.Type())
       fmt.Println("kind is float64:", v.Kind() == reflect.Float64)
       fmt.Println("value:", v.Float())
       // type: float64
       // kind is float64: true
       // value: 3.4 
    ```
  - Reflection goes from reflection object to interface value
    > Given a reflect.Value we can recover an interface value using the Interface method
    ```go
      var f float64 = 3.4  
      v := reflection.ValueOf(f)
      y := v.Interface().(float64) // y will have type float64.
      fmt.Println(y)
      // print
      // 3.4
    ```
    
  - To modify a reflection object, the value must be settable
   > Settability is a bit like addressability, but stricter. It's the property that a reflection object can modify
   the actual storage that was used to create the reflection object. Settability is determined by whether the reflection object holds the original item
   
   ```go
     // x is not addressable, the reason is that x is passed to reflect.ValueOf(...) as a copy of x
     //, not itself
     var x float64 = 3.4
     v := reflect.ValueOf(x)
     v.SetFloat(7.1)
     // print
     // panic: reflect.Value.SetFloat using unaddressable value
   ```
   > To allow to modify reflection value and see the effect from this modification
   ```go
    
   ```
### APIs

### References
- [Laws of reflection](https://blog.golang.org/laws-of-reflection)
- [Reflection in Go](https://www.integralist.co.uk/posts/reflection-in-go/)
- [Reflection In Usage](https://github.com/a8m/reflect-examples)