package handlers

import (
    "context"
    "encoding/json"
    "io"
    "net/http"

    "github.com/nspcc-dev/neofs-sdk-go/client"
    "github.com/nspcc-dev/neofs-sdk-go/container"
    "web3-onlyfans/services/frostfs-uploader/internal/frostuploader"
    "web3-onlyfans/services/frostfs-uploader/internal/utils"
)

func UploadHandler(w http.ResponseWriter, r *http.Request, cfg *utils.Config, logger utils.Logger) {
    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "no file found", http.StatusBadRequest)
        return
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        http.Error(w, "read file error", http.StatusInternalServerError)
        return
    }

    // Подключаемся к FrostFS
    cli, err := client.New(client.WithDefaultPrivateKeyStr(cfg.FrostfsPrivKey))
    if err != nil {
        logger.Errorf("client init: %v", err)
        http.Error(w, "FrostFS init error", http.StatusInternalServerError)
        return
    }
    defer cli.Close()

    if err = cli.Dial(cfg.FrostfsEndpoint); err != nil {
        logger.Errorf("dial error: %v", err)
        http.Error(w, "FrostFS dial error", http.StatusInternalServerError)
        return
    }

    cntr, err := container.IDFromString(cfg.FrostfsContainer)
    if err != nil {
        logger.Errorf("invalid container: %v", err)
        http.Error(w, "invalid container", http.StatusInternalServerError)
        return
    }

    oid, err := frostuploader.UploadFile(context.Background(), cli, cntr, data)
    if err != nil {
        logger.Errorf("upload error: %v", err)
        http.Error(w, "upload error", http.StatusInternalServerError)
        return
    }

    resp := map[string]any{
        "status":    "ok",
        "object_id": oid,
        "filename":  header.Filename,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
