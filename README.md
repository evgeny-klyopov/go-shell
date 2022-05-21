# Go-shell

Use:
```shell
go get "github.com/evgeny-klyopov/go-shell"
```

### Example
example shell
```shell
#!/bin/sh
i=1
while true
do
  if [ "$i" -gt 30 ]; then
    echo "Stop"
    break;
  fi
  echo "Output - i = $i"
  i=$((i+1))
  sleep 1
done
echo "Finish"
```

Execute
```go
package main

import (
	"fmt"
	shell "github.com/evgeny-klyopov/go-shell"
)

func main() {
	s := shell.New("/bin/sh", "example", true)

	s.Run()
	i := 0
	for {
		val, ok := <-s.GetOutput()
		if ok == false {
			break
		} else {
			fmt.Println(val)
			i++
		}
	}
}
```
