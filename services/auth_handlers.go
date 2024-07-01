package services

import (
	"context"
	"database/sql"
	"ioutil"
	"json"
	"log"
	"math/rand"
	"net/http"

	"github.com/karan2704/kube-deploy/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	DB *sql.DB
	googleOauthConfig *oauth2.Config
	oauthStateString string
) 

const (
	letters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)


func SetDB(dbObject *sql.DB){
	DB = dbObject
}

func AuthInit(){
	cred, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read auth credentials: %v", err)
	}
	//fix auth config
	googleOauthConfig, err   = google.ConfigFromJSON(cred, drive.DriveMetadataReadonlyscope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

}

func LoginHandler(w http.ResponseWriter, r *http.Request){
	url := googleOauthConfig.AuthCodeURL(RandomStringGenerator())
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func OauthCallbackHandler(w http.ResponseWriter, r *http.Request){
	ctx := context.Background()

	if r.FormValue("state") != oauthStateString {
		log.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, r.FormValue("state"))
        http.Error(w, "Invalid state parameter", http.StatusBadRequest)
        return
	}

	code := r.FormValue("code")
	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		log.Printf("code exchange failed: %s\n", err.Error())
		http.Error(w, "Failed to fetch access token", http.StatusInternalServerError)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
        log.Printf("failed getting user info: %s\n", err.Error())
        http.Error(w, "Failed to get user info", http.StatusInternalServerError)
        return
    }
	defer response.Body.Close()

	userInfo, err := ioutil.ReadAll(response.Body)
	if err != nil {
        log.Printf("failed reading response body: %s\n", err.Error())
        http.Error(w, "Failed to read profile details", http.StatusInternalServerError)
        return
    }

	var user map[string]interface{}
	if err := json.Unmarshal(userInfo, &user); err != nil {
		log.Printf("failed to unmarshal user info: %s\n", err.Error())
        http.Error(w, "Failed to unmarshal json", http.StatusInternalServerError)
        return
	}

	if UserExists(user["id"].(string)){
		log.Printf("User %s logged in", user["email"])
	} else{
		if err:=SaveUserDetails(user); err!=nil{
			log.Printf("failed to unmarshal user info: %s\n", err.Error())
        	http.Error(w, "Could not save user details", http.StatusInternalServerError)
			return
		}
		log.Printf("User %s signed up\n", user["email"])
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(user)
}

func UserExists(userId string)(bool){
	var exists bool
	err := DB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE id=$1)", userId).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
        log.Printf("Error checking if user exists: %s\n", err)
    }
	return exists
}

func SaveUserDetails(user map[string]interface{})(error){
	_, err := DB.Exec(models.SaveUserInfo, user["id"], user["name"], user["email"])
	return err
}

func RandomStringGenerator()(string){
	 b := ""
    for i:=0; i<=10; i++ {
        b += string(letters[rand.Intn(len(letters))])
    }
	return b
}