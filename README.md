# NASYNC 

a customizable async task pool for golang, (event bus, runtime)

## Fetures

* less memory
* more effective
* max gorutines and memory customizable
* more safe

# WHY

golang is something easy but Fallible language, you may do this 

```
    func yourfucntion() {
    
            go dosomething()  // DO NOT DO THIS IN SERVER SIDE APP !
        
    }

```

you may get "too many open files" error, when your application  in High load, so you need this 


## Simple Usage

```
go get github.com/leenanxi/nasync
```

```
import "github.com/leenanxi/nasync"

func yourfucntion() {
        nasync.Do(func,params...)
}

```


## Advanced Usage
```
import "github.com/leenanxi/nasync"

func main() {
        //new a async pool in max 1000 task in max 1000 gorutines
        async := nasync.New(1000,1000)
        defer async.Close()
        async.Do(yourfunc,yourparams)
}

func doSometing(msg string) string{
	return "i am done by " + msg
}

```



## What if something callback ?

```
import "github.com/leenanxi/nasync"

func main() {
        nasync.Do(func() {
        		result := doSometing("msg")
        		fmt.Println("i am call back by ",result)
        })
}

func doSometing(msg string) string{
	return "i am done by " + msg
}

```




