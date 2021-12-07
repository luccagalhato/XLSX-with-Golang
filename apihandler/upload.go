package apihandler

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"social-planilha/models"
	"sync"
	"time"
)

var mu sync.Mutex

var files map[string]*models.Excel = map[string]*models.Excel{}

func addExcel(name string, value io.Reader) string {
	key := genKey()
	mu.Lock()
	files[key] = &models.Excel{Name: name, Value: value}
	mu.Unlock()
	fmt.Println("fazendo Download...")
	go func() {
		time.Sleep(time.Second * 1000)
		mu.Lock()
		delete(files, key)
		mu.Unlock()
	}()

	return key
}

func genKey() string {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, rand.Uint64())
	return base64.URLEncoding.EncodeToString(bs)
}

//GetFile ...
func GetFile(id string) *models.Excel {
	fmt.Println("Buscando id...")
	mu.Lock()
	defer mu.Unlock()
	return files[id]
}
