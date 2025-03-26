#### Simple image resizer with no dependencies, written in GO 

1. Execute `go run .`.
2. HTTP POST images as multipart/formdata to root path with the desired height  & width in pixels as url query params.
```http request
HTTP POST localhost:8080/?height=300&width=400
Content-Type: multipart/form-data

...images
```
3. Resized images will be returned; in the case of multiple images they will be zipped. 