package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/nfnt/resize"
)

type Page struct {
	Title string
	Body  []byte
}

type UploadResponse struct {
	Path string
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	imageType := handler.Header.Get("Content-Type")

	if imageType == "image/jpeg" || imageType == "image/png" || imageType == "image/gif" {
		switch imageType {
		case "image/jpeg":

			img, _, err := image.Decode(file)
			resizedImage := resize.Resize(1000, 1000, img, resize.Lanczos3)
			path := fmt.Sprintf("media/%s", handler.Filename)

			out, err := os.Create(path)
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			// write new image to file
			jpeg.Encode(out, resizedImage, nil)

			if err != nil {
				fmt.Println(err)
			}

			data := UploadResponse{Path: path}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(data)
			break

		case "image/png":

			img, _, err := image.Decode(file)
			resizedImage := resize.Resize(1000, 1000, img, resize.Lanczos3)
			path := fmt.Sprintf("media/%s", handler.Filename)

			out, err := os.Create(path)
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			// write new image to file
			png.Encode(out, resizedImage)
			if err != nil {
				fmt.Println(err)
			}

			data := UploadResponse{Path: path}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(data)
			break

		case "image/gif":
			newGifImg := gif.GIF{}
			gifImg, err := gif.DecodeAll(file)
			if err != nil {
				log.Fatal(err)
			}

			for _, img := range gifImg.Image {
				resizedGifImg := resize.Resize(500, 0, img, resize.Lanczos2)
				palettedImg := image.NewPaletted(resizedGifImg.Bounds(), img.Palette)
				draw.FloydSteinberg.Draw(palettedImg, resizedGifImg.Bounds(), resizedGifImg, image.ZP)

				newGifImg.Image = append(newGifImg.Image, palettedImg)
				newGifImg.Delay = append(newGifImg.Delay, 25)
			}
			path := fmt.Sprintf("media/%s", handler.Filename)
			out, err := os.Create(path)
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()

			gif.EncodeAll(out, &newGifImg)
			if err != nil {
				log.Fatal(err)
			}

			data := UploadResponse{Path: path}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(data)

			break
		}

	} else {
		fmt.Println(errors.New("Eror.A file should be either png, jpeg or gif"))
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	fmt.Println(title)
	t, _ := template.ParseFiles("templates/home.html")
	t.Execute(w, &Page{Title: "Resizer"})
}

func setupRoutes() {
	fs := http.FileServer(http.Dir("./media"))
	http.HandleFunc("/", homePage)
	http.HandleFunc("/upload", uploadFile)
	http.Handle("/media/", http.StripPrefix("/media/", fs))
	http.ListenAndServe(":8080", nil)
}

func main() {
	setupRoutes()
}
