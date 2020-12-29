package main

import (
	"encoding/json"
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

<<<<<<< HEAD
type UploadRequest struct {
	Width  int
	Height int
}

func getPictureHash(parts [3]string) string {
	hashString := strings.Join(parts[:], ":")
	h := md5.Sum([]byte(hashString))
	return hex.EncodeToString(h[:])
}

func doResize(file multipart.File, filename string, width string, height string, imageType string) string {
	var parts [3]string
	parts[0] = filename
	parts[1] = width
	parts[2] = height
	hash := getPictureHash(parts)
	path := fmt.Sprintf("media/%s%s", hash, filepath.Ext(filename))
	if _, err := os.Stat(path); err == nil {
		return path
	}

	fmt.Println("doResize")
	if imageType == "jpg" || imageType == "jpeg" || imageType == "png" {
		img, _, err := image.Decode(file)
		width64, _ := strconv.ParseUint(width, 10, 32)
		height64, _ := strconv.ParseUint(height, 10, 32)
		m := resize.Resize(uint(width64), uint(height64), img, resize.Lanczos3)

		out, err := os.Create(path)
		if err != nil {
			return path
		}
		defer out.Close()
		if imageType == "jpeg" || imageType == "jpg" {
			jpeg.Encode(out, m, nil)
		} else if imageType == "png" {
			png.Encode(out, m)
		}
	} else if imageType == "gif" {
		newGifImg := gif.GIF{}
		width64, _ := strconv.ParseUint(width, 10, 32)
		height64, _ := strconv.ParseUint(height, 10, 32)
		gifImg, err := gif.DecodeAll(file)
		if err != nil {
			log.Fatal(err)
		}

		for _, img := range gifImg.Image {
			resizedGifImg := resize.Resize(uint(width64), uint(height64), img, resize.Lanczos2)
			palettedImg := image.NewPaletted(resizedGifImg.Bounds(), img.Palette)
			draw.FloydSteinberg.Draw(palettedImg, resizedGifImg.Bounds(), resizedGifImg, image.ZP)

			newGifImg.Image = append(newGifImg.Image, palettedImg)
			newGifImg.Delay = append(newGifImg.Delay, 25)
		}
		out, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		gif.EncodeAll(out, &newGifImg)
		if err != nil {
			log.Fatal(err)
		}
	}
	return path
}

=======
>>>>>>> 0c63fd01a7a1f30d67ab5804dd3dfee2b2068c69
func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	r.ParseMultipartForm(0)
	file, handler, err := r.FormFile("image")
	width := r.FormValue("width")
	height := r.FormValue("height")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
<<<<<<< HEAD
	imageExtension := handler.Header.Get("Content-type")

	if imageExtension == "image/jpeg" || imageExtension == "image/jpg" {
		path := doResize(file, handler.Filename, width, height, "jpeg")

		data := UploadResponse{Path: path}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
	} else if imageExtension == "image/png" {
		path := doResize(file, handler.Filename, width, height, "png")

		data := UploadResponse{Path: path}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)

	} else if imageExtension == "image/gif" {
		path := doResize(file, handler.Filename, width, height, "gif")

		data := UploadResponse{Path: path}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
=======

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
>>>>>>> 0c63fd01a7a1f30d67ab5804dd3dfee2b2068c69
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
