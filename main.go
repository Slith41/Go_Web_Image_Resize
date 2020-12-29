package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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
	"strings"
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

type UploadRequest struct {
	Width  int
	Height int
}

func getPictureHash(parts [3]string) string {
	hashString := strings.Join(parts[:], ":")
	h := md5.Sum([]byte(hashString))
	return hex.EncodeToString(h[:])
}

func doResize(file multipart.File, filename string, width string, height string) string {
	var parts [3]string //здесь мы создаем структуру, которая включает в себя элементы, которые будут входить в параметры кэша.
	parts[0] = filename
	parts[1] = width //width и height представлены в виде строк так как будут вводиться пользователем
	parts[2] = height
	hash := getPictureHash(parts)
	path := fmt.Sprintf("media/%s%s", hash, filepath.Ext(filename))
	if _, err := os.Stat(path); err == nil {
		return path
	}

	fmt.Println("doResize")
	if filepath.Ext(filename) == ".jpeg" || filepath.Ext(filename) == ".jpg" || filepath.Ext(filename) == ".png" {
		img, _, err := image.Decode(file)
		width64, _ := strconv.ParseUint(width, 10, 32) //конвертация введенных данных в uint для передачи в функцию ресайза
		height64, _ := strconv.ParseUint(height, 10, 32)
		m := resize.Resize(uint(width64), uint(height64), img, resize.Lanczos3)

		out, err := os.Create(path)
		if err != nil {
			return path
		}
		defer out.Close()
		if filepath.Ext(filename) == ".jpeg" || filepath.Ext(filename) == ".jpg" { //так как Encode jpeg и png отличвается, мы делаем
			jpeg.Encode(out, m, nil) //проверку на принадлежность к одному из типов и затем используем соответствующие функции
		} else if filepath.Ext(filename) == ".png" {
			png.Encode(out, m)
		}
	} else if filepath.Ext(filename) == ".gif" { //гифка вынесена, так как требует немного другой обработки
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
	if filepath.Ext(handler.Filename) == ".jpeg" || filepath.Ext(handler.Filename) == ".jpg" {
		path := doResize(file, handler.Filename, width, height)

		data := UploadResponse{Path: path}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
	} else if filepath.Ext(handler.Filename) == ".png" {
		path := doResize(file, handler.Filename, width, height)

		data := UploadResponse{Path: path}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)

	} else if filepath.Ext(handler.Filename) == ".gif" {
		path := doResize(file, handler.Filename, width, height)

		data := UploadResponse{Path: path} //Следующий блок - работа с JSON, который в итоге передаст путь к сохраненному файлу в HTML
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
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
