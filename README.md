Command line interface
======================

### License
[The MIT License (MIT)](https://zh.wikipedia.org/wiki/MIT許可證)


### Features

1. Based on golang tag. Support three tags: cli,usage,dft
2. Support specify default value: use dft tag
3. Support required declaration: cli tag with prefix `*`
4. Support multi flag name for same one field: like `cli:"h,help"`

### Usage
	
First, you should define a struct, like this:
```go
type Args struct {
	OnlySingle     bool    `cli:"v" usage:"only single char"`
	ManySingle     string  `cli:"X,Y" usage:"many single char"`
	SingleAndMulti int     `cli:"s,single-and-multi" usage:"single and multi"`
	OnlyMulti      uint    `cli:"only-multi" usage:"only multi"`
	Required       int8    `cli:"*required" usage:"required value"`
	Default        uint8   `cli:"id" usage:"default value" dft:"102"`
	Ignored        int16   `cli:"-" usage:"ignored field"`
	UnName         uint16  `usage:"unname field"`
	Int32          int32   `cli:"i32" usage:"type int32"`
	Uint32         uint32  `cli:"u32" usage:"type uint32"`
	Int64          int64   `cli:"i64" usage:"type int64"`
	Uint64         int64   `cli:"u64" usage:"type uint64"`
	Float32        float32 `cli:"f32" usage:"type float32"`
	Float64        float64 `cli:"f364" usage:"type float64"`
}
```

Then, call `cli.Parse` function:
```go
t := new(Args)
flagSet := cli.Parse(os.Args[1:], t)
if flagSet.Error != nil {
	//TODO: handle the error
}
//NOTE: show help
// fmt.Printf("Usage of `%s'`: \n%s", os.Args[0], flagSet.Usage)
```

If you only want show help, you can directly call `cli.Usage` function:
```go
usage := cli.Usage(new(Args))
fmt.Printf("Usage of `%s'`: \n%s", os.Args[0], usage)
```

### Tags

#### cli

*cli* tag support singlechar format and multichar format, e.g.

```go
Help    bool    `cli:"h,help"`
Version string  `cli:"version"`
Port    int     `cli:"p"`
XYZ     bool    `cli:"x,y,z,xyz,XYZ"` 
```

The argument is required if *cli* tag has prefix `*`, e.g.

```go
Required string `cli:"*required"`
```

#### usage

*usage* tag describe the argument. If the argument is required, describe string has `*` prefix while show usage(`*` is red on unix-like os).

#### dft
*dft* tag specify argument default value.
