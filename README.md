# dflag
yet another flag library

![](https://github.com/petderek/dflag/actions/workflows/go.yml/badge.svg)

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

The equivalent code in the standard lib would look like this:

```
var name,age string

flag.StringVar(&name, "name", "", "")
flag.IntVar(&age, "age", 0, "")
flag.Parse()

fmt.Printf("%s is age %d", name, age)

```

You can also use struct tags to be more explicit.

```
// struct tags correspond to the flag stdlib name, value, and usage
var flags = struct{
    Foo  string `name:"foo" value:"" usage:"the foo value"`
    Bar  string `name:"bar" value:"blah" usage:"the bar value"`
 }{}
 
 func init() {
     dflag.Parse(&flags)
     fmt.Println(flags.Foo)
 }
 ```
 
 Check out the `testdata` directory for more examples.
 
 Duplicates, incorrect values in the tags, or other structural issues will panic. Invalid
 inputs will call os.Exit(), just like the stdlib.
 
 
