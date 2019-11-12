package main

import (
        "bytes"
        "flag"
        "fmt"
        "log"
        "os/exec"
        "strings"
)

func main() {
        // get my args.
        flag.Parse()
        myargs := flag.Args()

        // make the ps args with "u" option.
        // set the "u" option means show memory usage.
        psargs := "u" + strings.Join(myargs, "")

        // prepare the ps command.
        cmd := exec.Command("ps", psargs)
        var out bytes.Buffer
        cmd.Stdout = &out

        // execute the ps command.
        err := cmd.Run()
        if err != nil {
                log.Fatal(err)
        }

        // temporary: show raw stdout.
        fmt.Println(out.String())

        // finally: show stdout but RSS and VSZ has human-readable formatted byte as like "2.6Gi" 
        // ** TODO **

}
