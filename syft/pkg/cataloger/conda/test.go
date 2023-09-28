package main

 import (
         "compress/bzip2"
         "fmt"
         "io"
         "os"
 )

 func main() {

         inputFile, err := os.Open("/Users/jwillhite/Downloads/syft_test/pydantic-1.10.8-py311h6c40b1e_0.tar.bz2")

         if err != nil {
                 fmt.Println(err)
                 os.Exit(1)
         }

         defer inputFile.Close()

         outputFile, err := os.Create("./file.txt")

         if err != nil {
                 fmt.Println(err)
                 os.Exit(1)
         }

         defer outputFile.Close()

         bzip2reader := bzip2.NewReader(inputFile)

         if err != nil {
                 fmt.Println(err)
                 os.Exit(1)
         }

         io.Copy(outputFile, bzip2reader)

 }