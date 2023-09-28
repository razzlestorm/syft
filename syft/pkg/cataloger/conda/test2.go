package main

import (
	"archive/tar"
	"compress/bzip2"
	"path/filepath"
	"fmt"
	"io"
	"os"
)

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func Untar(dst string, r io.Reader) error {

	gzr := bzip2.NewReader(r)

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			
			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
	return nil
}



func main() {
	inputFile, err := os.Open("/Users/jwillhite/Downloads/syft_test/pydantic-1.10.8-py311h6c40b1e_0.tar.bz2")
	//inputFileTwo, err := os.Open("/Users/jwillhite/Downloads/syft_test/airflow-1.10.7-py37_0.conda")
	if err != nil {
		fmt.Println("THIS IS AN ERROR")
	}
	// Create and add some files to the archive.
	// var buf bytes.Buffer
	// tw := tar.NewWriter(&buf)
	// var files = []struct {
	// 	Name, Body string
	// }{
	// 	{"readme.txt", "This archive contains some text files."},
	// 	{"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
	// 	{"todo.txt", "Get animal handling license."},
	// }
	// for _, file := range files {
	// 	hdr := &tar.Header{
	// 		Name: file.Name,
	// 		Mode: 0600,
	// 		Size: int64(len(file.Body)),
	// 	}
	// 	if err := tw.WriteHeader(hdr); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	if _, err := tw.Write([]byte(file.Body)); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// if err := tw.Close(); err != nil {
	// 	log.Fatal(err)
	// }

	// Open and iterate through the files in the archive.
	message := Untar("~/Downloads", inputFile)
	if message != nil {
		fmt.Println("ERROR")
		fmt.Println(message)
	}
	if message == nil {
		fmt.Println("SUCCESS")
	}
	/* 
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Contents of %s:\n", hdr.Name)
		if _, err := io.Copy(os.Stdout, tr); err != nil {
			log.Fatal(err)
		}
		fmt.Println()
	}
	*/

}
