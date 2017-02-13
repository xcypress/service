package network

import "io"



//protocol
package main

import (
"encoding/binary"
"errors"
"io"
)

//    [x][x][x][x][x][x][x][x]...
//    |  (int32) || (binary)
//    |  4-byte  || N-byte
//    ------------------------...
//        size       data
func ReadMessage(r io.Reader) ([]byte, error) {
    var msgSize int32

    // message size
    err := binary.Read(r, binary.BigEndian, &msgSize)
    if err != nil {
        return nil, err
    }

    if msgSize > 1000 {
        return nil, err
    }
    // message binary data
    buf := make([]byte, msgSize)
    _, err = io.ReadFull(r, buf)
    if err != nil {
        return nil, err
    }

    return buf, nil
}

func WriteMessage(r io.Reader,b []byte) {

}

//    [x][x][x][x][x][x][x][x]...
//    |  (int32) || (binary)
//    |  4-byte  || N-byte
//    ------------------------...
//      msg ID     data
//
func UnpackMessage(msg []byte) (int32, []byte, error) {
    if len(msg) < 4 {
        return -1, nil, errors.New("length of msg is too small")
    }
    return int32(binary.BigEndian.Uint32(msg)), msg[4:], nil
}

func PackMessage(msgId int32, msg []byte) ([]byte, error) {
    data := make([]byte, 8+len(msg))
    binary.BigEndian.PutUint32(data, uint32(4+len(msg)))
    binary.BigEndian.PutUint32(data[4:], uint32(msgId))
    n := copy(data[8:], msg)
    if n != len(msg) {
        return nil, nil
    }
    return data, nil
}

