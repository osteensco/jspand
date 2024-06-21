package main

import (
    "fmt"
)



func main() {

    fmt.Println("hello jspand")
    


    // jspand takes a small json file and blows it up so that it can be used for testing stuff
    // benefit of this is that it will retain your schema but scale it up


    // TODO
    // parse args
    //      - ex. jspand [filepath] [optional integer]
    //          - integer arg is optional and if not provided defaults to 50000?
    //          - maybe some other features in the future

    // validate
    //      - verify it's a json file
    //      - verify it's under a certain size

    // take keys and duplicate key: value pairs
    //      - ex. key -> key1, anotherKey -> anotherKey1
    //      - repeat until set amount of keys desired (or default of course)

    // output the file

}




