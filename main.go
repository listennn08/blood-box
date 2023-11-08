package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v56/github"
	"github.com/joho/godotenv"
)

type Blood struct {
	City string `json:"city"`
	A型   string `json:"A型"`
	B型   string `json:"B型"`
	O型   string `json:"O型"`
	AB型  string `json:"AB型"`
}
type BloodResp struct {
	UpdateTime string  `json:"updateTime"`
	Citys      []Blood `json:"citys"`
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		BLOOD_API_URL = os.Getenv("BLOOD_API_URL")
		GIST_ID       = os.Getenv("GIST_ID")
		TOKEN         = os.Getenv("TOKEN")
		DISPLAY_TEXT  = map[string]string{
			"less":   "偏低",
			"normal": "正常",
			"lack":   "急缺",
		}
		DISPLAY_TEMPLATE = []string{
			"｜血型／城市｜",
			"｜Ａ        ｜",
			"｜Ｂ        ｜",
			"｜Ｏ        ｜",
			"｜ＡB       ｜",
		}
	)

	httpClient := &http.Client{}
	client := github.NewClient(nil).WithAuthToken(TOKEN)
	ctx := context.Background()

	req, err := http.NewRequest("GET", BLOOD_API_URL, nil)
	if err != nil {
		panic(err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic(resp.Status)
	}

	var data BloodResp

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		panic(err)
	}

	filename := fmt.Sprintf("🩸 血液庫存 %s", data.UpdateTime)
	for _, v := range data.Citys {
		DISPLAY_TEMPLATE[0] += fmt.Sprintf("%s｜", v.City)
		DISPLAY_TEMPLATE[1] += fmt.Sprintf("%s｜", DISPLAY_TEXT[v.A型])
		DISPLAY_TEMPLATE[2] += fmt.Sprintf("%s｜", DISPLAY_TEXT[v.B型])
		DISPLAY_TEMPLATE[3] += fmt.Sprintf("%s｜", DISPLAY_TEXT[v.O型])
		DISPLAY_TEMPLATE[4] += fmt.Sprintf("%s｜", DISPLAY_TEXT[v.AB型])
	}
	content := strings.Join(DISPLAY_TEMPLATE, "\n")

	original, _, err := client.Gists.Get(ctx, GIST_ID)
	if err != nil {
		panic(err)
	}

	files := original.Files

	for k := range files {
		if k == github.GistFilename(filename) {
			continue
		}
		client.Gists.Edit(ctx, GIST_ID, &github.Gist{
			Files: map[github.GistFilename]github.GistFile{
				github.GistFilename(k): {
					Content: nil,
				},
			},
		})
		fmt.Printf("Delete file: %s\n", k)
	}

	client.Gists.Edit(ctx, GIST_ID, &github.Gist{
		Files: map[github.GistFilename]github.GistFile{
			github.GistFilename(filename): {
				Content: &content,
			},
		},
	})
	fmt.Println("Create file: ", filename)
}
