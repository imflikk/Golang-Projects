# Golang-Projects

Various programs written in Golang from the book [Black Hat Go](https://nostarch.com/blackhatgo) and the site [Gophercises](https://gophercises.com/).  Some are the same as the demoes provided and others have been modified for me to be able to practice doing various things in Go.

## Updates made:
- Black Hat Go
  - ch2/scanner
    - Updated scanner to take command-line arguments instead of hard-coding target/ports.
  - ch3/bing-metadata-scraper
    - Updated dork used to retrieve specific file types as the original doesn't work anymore
    - Page structure of results appears to have changed breaking the original method of retrieving file.  Added regex check to locate URL of the actual document in each link followed and make a request to that before processing file metadata.
    - Added functionality to check any additional pages of results, if they exist.  This isn't perfect at the moment, but works well enough for now.

