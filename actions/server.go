package actions

import (
	"archive/tar"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/codegangsta/cli"
)

func build(w http.ResponseWriter, r *http.Request) {
	tr := tar.NewReader(r.Body)
	defer r.Body.Close()
	tmpDir, err := ioutil.TempDir("", "gci_server_builds")
	if err != nil {
		// FAIL!
	}
	defer os.Remove(tmpDir)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Contents of %s:\n", hdr.Name)
		if _, err := io.Copy(os.Stdout, tr); err != nil {
			log.Fatalln(err)
		}
		fmt.Println()
	}

}

func Server(c *cli.Context) {
	mux := http.NewServeMux()
	mux.HandleFunc("/build", build)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
