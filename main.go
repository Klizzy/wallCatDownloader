package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gen2brain/beeep"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	fileName      string
	wallCatUrl    string
	wallCatImgUrl string
)

// Set your target folder here
const targetPath = ""

func main() {

	wallCatUrl = "https://beta.wall.cat"

	resp, err := httpClient().Get(wallCatUrl)
	checkError(err)
	defer resp.Body.Close()
	findImgUrl(err, resp)

	// Build fileName from fullPath
	buildFileName()

	// Create blank file
	file := createFile()

	// Put content on file
	putFile(file, httpClient())
}

func putFile(file *os.File, client *http.Client) {
	resp, err := client.Get(wallCatImgUrl)

	checkError(err)

	defer resp.Body.Close()
	throwNotFound(resp)

	size, err := io.Copy(file, resp.Body)
	readableFileSize := ByteToHumanReadable(size)

	defer file.Close()

	checkError(err)

	beeep.Notify("Wallcat Downloader", "Image download successfull ("+readableFileSize+")", "assets/iURyERju_400x400.png")

	fmt.Println("Just Downloaded a file %s with size %d", fileName, readableFileSize)
}
func setVariables(imgUrl string, statusCode bool) {
	wallCatImgUrl = imgUrl
}

func findImgUrl(err error, resp *http.Response) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)

	// Find the review items
	doc.Find("article").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		name, test := s.Attr("data-image-bg-url")
		setVariables(name, test)
	})

}

func throwNotFound(resp *http.Response) {
	if resp.StatusCode == 404 {
		panic("Image not found")
	}
}

func buildFileName() {
	fileUrl, err := url.Parse(wallCatImgUrl)
	checkError(err)

	path := fileUrl.Path
	segments := strings.Split(path, "/")

	fileName = segments[len(segments)-1] + ".jpeg"
}

func httpClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	return &client
}

func createFile() *os.File {
	file, err := os.Create(fileName)
	os.Rename(fileName, targetPath+fileName)

	checkError(err)
	return file
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func ByteToHumanReadable(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
