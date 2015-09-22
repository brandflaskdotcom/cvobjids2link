package main

import (
	"bufio"
	"flag"
	"fmt"
	fb "github.com/huandu/facebook"
	"log"
	"os"
	"time"
)

const (
	FB_APP_ID             = "1380374468852725"
	FB_APP_SECRET         = "9f11a4d965333c7156981415174fd6fb"
	FB_VALID_ACCESS_TOKEN = ""
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func get_source(session *fb.Session, fbobjid string) {
	fbapi := fmt.Sprintf("/%v/?fields=source,from", fbobjid)

	res, err := session.Get(fbapi, nil)
	if err != nil {
		return
	}
	fid := res["from"].(map[string]interface{})["id"]
	bname := res["from"].(map[string]interface{})["name"]
	fmt.Printf("\"%v\",\"%v\",\"%v\",\"%v\"\n", fbobjid, fid, bname, res["source"])

}

func main() {
	access_token := flag.String("t", FB_VALID_ACCESS_TOKEN, "a valid access token")
	from := flag.Int("f", 0, "a valid access token")
	inputFile := flag.String("i", "default", "a path to read file")

	flag.Parse()

	if *inputFile == "default" {
		return
	}
	file, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// create a global App var to hold app id and secret.
	var globalApp = fb.New(FB_APP_ID, FB_APP_SECRET)

	// facebook asks for a valid redirect uri when parsing signed request.
	// it's a new enforced policy starting in late 2013.
	globalApp.RedirectUri = "https://www.facebook.com/connect/login_success.html"

	// if there is another way to get decoded access token,
	// creates a session directly with the token.
	session := globalApp.Session(*access_token)

	// validate access token. err is nil if token is valid.
	err = session.Validate()
	if err != nil {
		return
	}

	fmt.Println("objid,ownerid,name,imgurl")
	linenumber := 0
	for scanner.Scan() {
		fb_obj_id := scanner.Text()
		log.Printf("%v %v\n", linenumber, fb_obj_id)

		if linenumber >= *from {
			get_source(session, fb_obj_id)
			time.Sleep(2000 * time.Millisecond)
		}
		linenumber += 1
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return

}
