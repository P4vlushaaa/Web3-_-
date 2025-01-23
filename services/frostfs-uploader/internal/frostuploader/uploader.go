package frostuploader

import (
    "context"
    "io"

    "github.com/nspcc-dev/neofs-sdk-go/client"
    "github.com/nspcc-dev/neofs-sdk-go/client/object"
    "github.com/nspcc-dev/neofs-sdk-go/container"
)

func UploadFile(ctx context.Context, cli *client.Client, cntr container.ID, data []byte) (string, error) {
    obj := object.New()
    obj.SetPayload(data)

    writer, err := cli.ObjectPutInit(ctx, cntr, obj)
    if err != nil {
        return "", err
    }
    if _, err := writer.Write(data); err != nil {
        return "", err
    }
    oid, err := writer.Close()
    if err != nil {
        return "", err
    }
    return oid.String(), nil
}
