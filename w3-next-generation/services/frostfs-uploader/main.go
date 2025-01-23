package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nspcc-dev/neofs-sdk-go/client"
	"github.com/nspcc-dev/neofs-sdk-go/client/object"
	"github.com/nspcc-dev/neofs-sdk-go/container"
	// ...
)

var (
	frostfsEndpoint  = os.Getenv("FROSTFS_ENDPOINT") // типа "grpcs://localhost:8080"
	frostfsPrivKey   = os.Getenv("FROSTFS_PRIVKEY")
	frostfsContainer = os.Getenv("FROSTFS_CONTAINER_ID")
)

func main() {
	// Инициализация
	r := mux.NewRouter()
	r.HandleFunc("/upload", handleUpload).Methods("POST")

	// HTTP-сервер
	addr := ":8081"
	log.Println("Starting FrostFS uploader on", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	// Пример: берем файл из multipart/form-data
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file error: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Читаем в []byte (для простоты, но лучше стримить)
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "read error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Подключаемся к FrostFS
	cli, err := client.New(client.WithDefaultPrivateKeyStr(frostfsPrivKey))
	if err != nil {
		http.Error(w, "client init error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer cli.Close()

	err = cli.Dial(frostfsEndpoint)
	if err != nil {
		http.Error(w, "dial error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	cntr, err := container.IDFromString(frostfsContainer)
	if err != nil {
		http.Error(w, "invalid container: "+err.Error(), http.StatusInternalServerError)
		return
	}

	oid, err := uploadFileToFrostFS(r.Context(), cli, cntr, data)
	if err != nil {
		http.Error(w, "frostfs upload error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok","object_id":"%s","filename":"%s"}`, oid, header.Filename)
}

func uploadFileToFrostFS(ctx context.Context, cli *client.Client, cntr container.ID, data []byte) (string, error) {
	obj := object.New()
	obj.SetPayload(data)

	writer, err := cli.ObjectPutInit(ctx, cntr, obj)
	if err != nil {
		return "", err
	}
	err = writer.Write(data)
	if err != nil {
		return "", err
	}

	oid, err := writer.Close()
	if err != nil {
		return "", err
	}
	return oid.String(), nil
}
