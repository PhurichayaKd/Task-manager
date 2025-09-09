package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	ggoogle "golang.org/x/oauth2/google"
)

// ------- Data from Google userinfo -------
type GoogleUser struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// ------- Wrapper -------
type GoogleOAuth struct {
	cfg *oauth2.Config
}

// โหลดค่า CLIENT_ID/SECRET/REDIRECT จาก .env
func NewGoogleOAuthFromEnv() *GoogleOAuth {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	secret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirect := os.Getenv("GOOGLE_REDIRECT_URL")
	if clientID == "" || secret == "" || redirect == "" {
		panic("missing GOOGLE_CLIENT_ID / GOOGLE_CLIENT_SECRET / GOOGLE_REDIRECT_URL")
	}

	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: secret,
		RedirectURL:  redirect,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     ggoogle.Endpoint,
	}
	return &GoogleOAuth{cfg: cfg}
}

// สร้างลิงก์ไปหน้าล็อกอินของ Google
func (g *GoogleOAuth) LoginURL(state string) string {
	return g.cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// แลก code เป็น token แล้วดึงข้อมูลผู้ใช้จาก Google
func (g *GoogleOAuth) FetchUser(ctx context.Context, code string) (*GoogleUser, error) {
	if code == "" {
		return nil, errors.New("empty code")
	}
	
	// สร้าง context ที่มี timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	tok, err := g.cfg.Exchange(ctxWithTimeout, code)
	if err != nil {
		return nil, err
	}
	
	// สร้าง HTTP client ที่มี timeout
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	ctxWithClient := context.WithValue(ctxWithTimeout, oauth2.HTTPClient, httpClient)
	client := g.cfg.Client(ctxWithClient, tok)

	// Create request with proper headers
	req, err := http.NewRequestWithContext(ctxWithTimeout, "GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("google userinfo: bad status")
	}

	var u GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}
