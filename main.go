package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"os/user"
	"strings"
	"sync"

	"github.com/gedex/go-instagram/instagram"
)

var instagramName = flag.String("n", "", "'Instangram user name such as: 'kingjames'")
var numOfWorkerPointers = flag.String("c", "2", "Number of concurrent rename workers. default = 2")

var mutex sync.Mutex
var FileIndex int = 0
var client *instagram.Client
var ClientID string

func GetFileIndex() (ret int) {
	mutex.Lock()

	ret = FileIndex
	FileIndex = FileIndex + 1

	mutex.Unlock()

	return ret
}

func init() {
	ClientID = os.Getenv("InstagramID")

	if ClientID != "" {
		log.Fatalln("Please set 'export InstagramID=xxxxx' as your environment variables")
	}
}

func DownloadWorker(destDir string, linkChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for target := range linkChan {
		var imageType string

		if strings.Contains(target, ".png") {
			imageType = ".png"
		} else {
			imageType = ".jpg"
		}

		resp, err := http.Get(target)
		
		if err != nil {
			log.Println("Http.Get\nerror: " + err.Error() + "\ntarget: " + target)
			continue
		}
		
		defer resp.Body.Close()

		m, _, err := image.Decode(resp.Body)
		if err != nil {
			log.Println("image.Decode\nerror: " + err.Error() + "\ntarget: " + target)
			continue
		}

		// Ignore small images
		bounds := m.Bounds()
		if bounds.Size().X > 300 && bounds.Size().Y > 300 {
			imgInfo := fmt.Sprintf("pic%04d", GetFileIndex())
			out, err := os.Create(destDir + "/" + imgInfo + imageType)
			
			if err != nil || !os.IsExist(err) {
				log.Printf("os.Create\nerror: %s", err)
				continue
			}
			defer out.Close()

			if imageType == ".png" {
				png.Encode(out, m)
			} else {
				jpeg.Encode(out, m, nil)
			}

			if FileIndex%30 == 0 {
				fmt.Println(FileIndex, " photos downloaded.")
			}
		}
	}
}

func FindPhotos(ownerName, albumName, userId, baseDir string) {
	totalPhotoNumber := 1
	var mediaList []instagram.Media
	var next *instagram.ResponsePagination
	var optionalParameters *instagram.Parameters
	var err error

	dir := fmt.Sprintf("%v/%v", baseDir, ownerName)
	os.MkdirAll(dir, 0755)

	linkChan := make(chan string)

	wg := new(sync.WaitGroup)

	for i := 1; i < 1; i++ {
		wg.Add(1);
		go DownloadWorker(dir, linkChan, wg);
	}

	for {
		maxId := ""

		if next != nil {
			maxId = next.NextMaxID;
		}

		optionalParameters = &instagram.Parameters{Count: 10, MaxID: maxId};
		mediaList, next, err = client.Users.RecentMedia(userId, optionalParameters);

		if err != nil {
			log.Fatal(err);

			break;
		}

		for _, media := range mediaList {
			totalPhotoNumber = totalPhotoNumber + 1;
			linkChan <- media.Images.StandardResolution.URL;
		}

		if len(mediaList) == 0 || next.NextMaxID == "" {
			break;
		} 
	}
}

func main() {
	flag.Parse();
	var inputUser string;

	if *instagramName == "" {
		log.Fatalln("You need input your name -n=name")
	}

	inputUser = *instagramName;

	user, _ := user.Current();
	baseDir := fmt.Sprintf("%v/Pictures/goInstagram", user.HomeDir);

	client = instagram.NewClient(nil);
	client.ClientID = ClientID;

	var userID string;
	searchUser, _, err := client.Users.Search(inputUser, nil);
	
	for _, user := range searchUser {
		if user.Username == inputUser {
			userID = user.ID;
		}
	}

	if userID == "" {
		log.Fatalln("Can't address user name: ", inputUser, err)
	}

	userFolderName := fmt.Sprintf("[%s]%s", userID, inputUser);
	fmt.Println("Starting download [", userID, "]", inputUser);
	FindPhotos(userFolderName, inputUser, userID, baseDir);
}
