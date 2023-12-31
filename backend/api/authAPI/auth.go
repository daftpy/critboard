package authAPI

import (
	"critboard-backend/database/query/queryUsers"
	"critboard-backend/pkg/auth"
	"encoding/json"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"time"
)

type Token struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       time.Time
}

func (a *AuthHandler) TwitchAuthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := a.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOnline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func (a *AuthHandler) TwitchCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("URL:", r.URL.String())         // Log the URL
		log.Println("Query Params:", r.URL.Query()) // Log the query parameters
		token, err := a.oauthConfig.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
		if err != nil {
			log.Println("Error exchanging code for token:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
		if err != nil {
			log.Println("Error creating request object:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
		req.Header.Set("Client-Id", a.oauthConfig.ClientID)
		res, err := client.Do(req)
		if err != nil {
			log.Println("Error fetching user info:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		defer res.Body.Close()

		// ID (twitchID) is returned as string
		var userInfo struct {
			Data []struct {
				ID    string `json:"id"`
				Login string `json:"login"`
				Email string `json:"email"`
			} `json:"data"`
		}

		if err := json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
			log.Println("Error decoding user info:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if len(userInfo.Data) == 0 {
			log.Println("No user data received")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Encrypt the access token
		encryptedAccess, err := auth.Encrypt([]byte(token.AccessToken), a.encryptionKey)
		if err != nil {
			http.Error(w, "Internal server Error", http.StatusInternalServerError)
			return
		}

		// Encrypt the refresh token
		encryptedRefresh, err := auth.Encrypt([]byte(token.RefreshToken), a.encryptionKey)
		if err != nil {
			http.Error(w, "Internal server Error", http.StatusInternalServerError)
			return
		}

		// Store the access and refresh tokens in memcache
		err = auth.StoreOAuthTokens(a.memcacheClient, userInfo.Data[0].ID, encryptedAccess, encryptedRefresh, token.Expiry)
		if err != nil {
			log.Println("Error storing OAuth tokens:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		_, err = queryUsers.CreateUser(a.db, userInfo.Data[0].ID, userInfo.Data[0].Login, userInfo.Data[0].Email)
		if err != nil {
			log.Println("Error creating user:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Create a user session, store it in the database, and send a session ID to the client
		a.sessionManager.Put(r.Context(), "userID", userInfo.Data[0].ID)

		log.Println("Successfully authenticated and session created")

		redirectURL := os.Getenv("SERVER_DOMAIN")
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}
