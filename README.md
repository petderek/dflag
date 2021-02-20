# dflag
yet another flag library

dflag is a way to bind cli flags to a struct instead of using individual variables.
It emphasizes simplicity and terseness over expressiveness and deep featuresets. Dflag
mostly replicates the flag stdlib but with less boilerplate.

A minimum use case:

```
var flags = struct {
    Name string
    Age int
}{}

dflag.Parse(&flags)

fmt.Printf("%s is age %d", flags.Name, flags.Age)

```

You can also use struct tags to be more explicit.

```
// struct tags correspond to the flag stdlib name, value, and usage
var flags = struct{
    Foo  string `name:"foo" value:"" usage:"the foo value"`
    Bar  string `name:"foo" value:"blah" usage:"the bar value"`
 }{}
 
 func init() {
     dflag.Parse(&flags)
     fmt.Println(flags.Foo)
 }
 ```
 
 Duplicates, incorrect values in the tags, or other structural issues will panic. Invalid
 inputs will call os.Exit(), just like the stdlib.
 
 
