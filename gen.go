package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

//go:generate go run gen.go

const (
	domain    = ".roblox.com"
	domainLen = len(domain)
)

var endpoints = []string{
	"accountsettings.roblox.com",
	"assetdelivery.roblox.com",
	"auth.roblox.com",
	"avatar.roblox.com",
	"badges.roblox.com",
	"billing.roblox.com",
	"catalog.roblox.com",
	"chat.roblox.com",
	"clientsettings.roblox.com",
	"develop.roblox.com",
	"followings.roblox.com",
	"friends.roblox.com",
	"gameinternationalization.roblox.com",
	"gamejoin.roblox.com",
	"games.roblox.com",
	"groups.roblox.com",
	"inventory.roblox.com",
	"locale.roblox.com",
	"notifications.roblox.com",
	"points.roblox.com",
	"presence.roblox.com",
	"publish.roblox.com",
	"thumbnails.roblox.com",
}

func main() {
	for _, endpoint := range endpoints {
		fmt.Print(endpoint)
		path := endpoint[:len(endpoint)-domainLen]

		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.Mkdir(path, os.ModePerm); err != nil {
				panic(err)
			}
		}

		if err := os.Chdir(path); err != nil {
			panic(err)
		}

		fmt.Print(": Getting versions ... ")
		versions, err := getEndpointVersions(endpoint)
		if err != nil {
			panic(err)
		}
		fmt.Println("Done")

		if len(versions) == 0 {
			panic(fmt.Errorf("failed to find any versions for endpoint %s", endpoint))
		}

		for _, version := range versions {
			fmt.Printf("- %s", version)
			if _, err := os.Stat(version); os.IsNotExist(err) {
				if err := os.Mkdir(version, os.ModePerm); err != nil {
					panic(err)
				}
			} else {
				if _, err := os.Stat(version + "/client"); err == nil {
					if _, err := os.Stat(version + "/models"); err == nil {
						// Skip generating for the already existing endpoints
						fmt.Println(": Skip!")
						continue
					}
				}
			}
			fmt.Println()

			if err := os.Chdir(version); err != nil {
				panic(err)
			}

			if _, err := os.Stat("swagger.json"); os.IsNotExist(err) {
				url := "https://" + endpoint + "/docs/json/" + version
				fmt.Printf("\tDownloading %s ... ", url)
				if err := downloadFile(url, "swagger.json"); err != nil {
					panic(err)
				}
				fmt.Println("Done")
			}

			hasPaths, err := checkIfSchemaHasPaths("swagger.json")
			if err != nil {
				panic(err)
			}

			if hasPaths {
				fmt.Print("\tGenerating source code ... ")
				cmd := exec.Command("swagger", "generate", "client", "--skip-validation")
				if err := cmd.Run(); err != nil {
					panic(err)
				}
				fmt.Println("Done")
			} else {
				fmt.Println("\tHas no paths! Skipping")
			}

			if err := os.Chdir("../"); err != nil {
				panic(err)
			}
		}

		if err := os.Chdir("../"); err != nil {
			panic(err)
		}
	}
}

func getEndpointVersions(endpoint string) ([]string, error) {
	resp, err := http.Get("https://" + endpoint + "/docs")
	if err != nil {
		return nil, err
	}

	z := html.NewTokenizer(resp.Body)

	var versions []string

	var isWithinSelect = false
	var isWithinOption = false

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return nil, errors.New("failed to parse HTML")
		case tt == html.StartTagToken:
			t := z.Token()

			switch t.Data {
			case "select":
				for _, v := range t.Attr {
					if v.Key == "id" && v.Val == "version-selector" {
						isWithinSelect = true
						break
					}
				}
			case "option":
				if isWithinSelect {
					isWithinOption = true
					for _, v := range t.Attr {
						if v.Key == "value" {
							versions = append(versions, v.Val)
						}
					}
				}
			}
		case tt == html.EndTagToken:
			if isWithinOption {
				isWithinOption = false
			} else if isWithinSelect {
				return versions, nil
			}
		}
	}
}

func downloadFile(url string, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get url %s got code %s", url, resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	return err
}

func checkIfSchemaHasPaths(path string) (bool, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}

	var value struct {
		Paths map[string]interface{} `json:"paths"`
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return false, err
	}

	return len(value.Paths) > 0, nil
}
