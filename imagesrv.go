/*---------------------------------------------------------------------+
 | Copyright (c) 2021 by Mark W. Kernodle
 |
 | FILE:     imagesrv.go
 | PURPOSE:  Serve up output of tree and file command for a directory in HTML
 |
 | NOTES:  For NetApp coding challenge 11/28/2021
 +--------------------------------------------------------------------*/

package main

import (
    "log"
    "net/http"
    "os"
    "fmt"
    "os/exec"
    "os/signal"
    "syscall"
    "bytes"
    "strings"
    "github.com/gorilla/mux"
)

const theTitle string = "NetApp tree for candidate files"

func main() {

    var theLines [][]byte
    var theLine string
    var err error
    var fragment string
    var muxIt = mux.NewRouter()
    var scanDir string
    var extens [3]string

// Parse arguments to get the directory to browse.
// Copy image tree info a temp directory, and then construct
// an index.html comprising only the image files. 
// Then serve up this file with the FileServer abstraction.

    scanDir = os.Args[1]

    extens[0] = ".jpg"
    extens[1] = ".png"
    extens[2] = ".gif"

    tmpDir,err := os.MkdirTemp("/var/tmp", "NetApp-test")
    if err != nil {
        log.Fatal(err)
    }

// We must copy the source tree into place due to a limitation of
// path routes from '/' in http.Fileserver. Probably better to have added
// a custome Handler to surmount this restriction.

    srcDir := scanDir
    theCmd := exec.Command("cp", "-R", srcDir, tmpDir)
    _,err = theCmd.Output()   

    if err != nil {
        log.Fatal(err)
    }

    theIndex := tmpDir + "/index.html"
    indexFile,err := os.Create(theIndex)
    if err != nil {
        log.Fatal(err)
    }

    preDir := strings.TrimPrefix(scanDir, "/var/tmp")
    srcDir =  tmpDir + "/" + preDir
    theCmd = exec.Command("tree", "-i", "-T", theTitle, "-h", "--du", "-H", preDir, srcDir)
    treeOut,err := theCmd.Output()
    if err != nil {
        log.Fatal(err)
    }

    theLines = bytes.Split(treeOut, []byte("\n"))
    for i,_ := range theLines {
        theLine = string(theLines[i])
        if  theLine != "" {
            fragment = fmt.Sprintf("%s\n", theLine)
            gotIt := strings.Contains(fragment,extens[0])
            if gotIt == false {
                gotIt = strings.Contains(fragment,extens[1])
            }
            if gotIt == false {
                gotIt = strings.Contains(fragment,extens[2])
            }

            if gotIt {
                hrefMark := strings.Index(fragment,"href=")
                fileN := fragment[hrefMark +  6:]
                fileL := strings.Index(fileN,"\"")
                fName := fileN[:fileL]
                fName = tmpDir + fName
// sanitize blank characters that have been subsumed with %20
                fName = strings.Replace (fName, "%20", " ", -1)
                theCmd:= exec.Command("file", "-b", fName)
                fileCmdOut,err := theCmd.Output()
                if err != nil {
                    log.Fatal(err)
                }
                fileCmdString := string(fileCmdOut)
                fileCmdString = strings.TrimRight(fileCmdString,"\n")
                snippet :=  " - " + "<small>" + fileCmdString + "</small>" + "<br>"
                enhanced := strings.Replace(fragment, "<br>", snippet, 1)
                bytesWritten, err := indexFile.Write([]byte(enhanced))
                if err != nil || bytesWritten <= 0 {
                    log.Fatal(err)
                }
            } else {
                bytesWritten, err := indexFile.Write([]byte(theLine))
                if err != nil || bytesWritten <= 0 {
                    log.Fatal(err)
                }
            }
        }
    }

    indexFile.Close()

    imgDir := http.Dir(tmpDir)
    fileViewer := http.FileServer(imgDir)

// Break out of the Listen/Server loop upon ctrl-c

    ourC := make(chan os.Signal)
    signal.Notify(ourC, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-ourC
        os.RemoveAll(tmpDir)
        os.Exit(0)
    }()

    muxIt.PathPrefix("/").Handler(http.StripPrefix("/", fileViewer))
    log.Fatal(http.ListenAndServe(":9080", muxIt))

}
