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
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

func doResize(file multipart.File, filename string, width string, height string, TypeOfImage string) (string, error) {
	var parts [4]string
	parts[0] = filename
	parts[1] = width
	parts[2] = height
	parts[3] = TypeOfImage

	path := fmt.Sprintf("media/%s", filepath.Ext(filename))
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	if TypeOfImage == "jpeg" || TypeOfImage == "jpg" || TypeOfImage == "png" {
		fmt.Println("doResize")
		img, _, err := image.Decode(file)
		width64, _ := strconv.ParseUint(width, 10, 32)
		height64, _ := strconv.ParseUint(height, 10, 32)
		m := resize.Resize(uint(width64), uint(height64), img, resize.Lanczos3)

		out, err := os.Create("media/" + filename)
		if err != nil {
			return path, err
		}
		defer out.Close()

		if TypeOfImage == "jpeg" {
			jpeg.Encode(out, m, nil)
		} else if TypeOfImage == "png" {
			png.Encode(out, m)
		}
		return path, err

	} else {
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
		path := fmt.Sprintf("media/%s", filename)
		out, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		gif.EncodeAll(out, &newGifImg)
		if err != nil {
			log.Fatal(err)
		}
		return path, err
	}
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

			doResize(file, handler.Filename, "100", "0", "jpeg")
			path := fmt.Sprintf("media/%s", filepath.Ext(handler.Filename))
			data := UploadResponse{Path: path}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(data)
			break

		case "image/png":
			doResize(file, handler.Filename, "100", "0", "png")
			path := fmt.Sprintf("media/%s", filepath.Ext(handler.Filename))
			data := UploadResponse{Path: path}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(data)
			break

		case "image/gif":
			doResize(file, handler.Filename, "100", "0", "gif")
			path := fmt.Sprintf("media/%s", filepath.Ext(handler.Filename))

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
