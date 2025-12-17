// Deprecated: renamed to admin_smoke.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/Yw332/campus-moments-go/internal/models"
	"github.com/Yw332/campus-moments-go/internal/routes"
	"github.com/Yw332/campus-moments-go/pkg/config"
	"github.com/Yw332/campus-moments-go/pkg/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func must(v interface{}, err error) interface{} {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return v
}

func main() {
	config.Init()
	database.Init()
	db := database.GetDB()
	if db == nil {
		log.Fatalln("数据库未连接")
	}

	// Ensure admin exists
	pwd, _ := bcrypt.GenerateFromPassword([]byte("Admin12345"), bcrypt.DefaultCost)
	admin := models.Admin{Username: "itestadmin", Password: string(pwd), CreatedAt: time.Now(), UpdatedAt: time.Now()}
	db.Where(models.Admin{Username: admin.Username}).FirstOrCreate(&admin)
	log.Printf("Using admin: %s (id=%d)", admin.Username, admin.ID)

	// Start test HTTP server
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	routes.SetupRoutes(r)
	srv := httptest.NewServer(r)
	defer srv.Close()

	hc := &http.Client{Timeout: 5 * time.Second}

	// Helper for POST JSON
	postJSON := func(path string, body interface{}, headers map[string]string) (*http.Response, []byte) {
		b, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", srv.URL+path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, err := hc.Do(req)
		if err != nil {
			log.Fatalf("POST %s failed: %v", path, err)
		}
		defer resp.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return resp, buf.Bytes()
	}

	// Helper for GET
	doGet := func(method, path string, headers map[string]string) (*http.Response, []byte) {
		req, _ := http.NewRequest(method, srv.URL+path, nil)
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, err := hc.Do(req)
		if err != nil {
			log.Fatalf("GET %s failed: %v", path, err)
		}
		defer resp.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return resp, buf.Bytes()
	}

	// 1. Admin login
	resp, body := postJSON("/admin/login", map[string]string{"account": "itestadmin", "password": "Admin12345"}, nil)
	log.Printf("POST /admin/login -> status=%d body=%s", resp.StatusCode, string(body))
	if resp.StatusCode != 200 {
		log.Fatalln("admin login failed")
	}
	var lr map[string]interface{}
	json.Unmarshal(body, &lr)
	token := lr["data"].(map[string]interface{})["token"].(string)
	auth := map[string]string{"Authorization": "Bearer " + token}

	// 2. Admin menu
	resp, body = doGet("GET", "/admin/menu", auth)
	log.Printf("GET /admin/menu -> status=%d body=%s", resp.StatusCode, string(body))
	if resp.StatusCode != 200 {
		log.Fatalln("admin menu failed")
	}

	// 3. Register normal user
	username := fmt.Sprintf("itestuser%v", time.Now().Unix())
	phone := fmt.Sprintf("138%08d", time.Now().Unix()%100000000)
	resp, body = postJSON("/auth/register", map[string]string{"username": username, "phone": phone, "password": "TestPass123"}, nil)
	log.Printf("POST /auth/register -> status=%d body=%s", resp.StatusCode, string(body))
	if resp.StatusCode != 200 {
		log.Fatalln("register failed")
	}
	var rr map[string]interface{}
	json.Unmarshal(body, &rr)
	userId := int64(rr["data"].(map[string]interface{})["userId"].(float64))

	// 4. Non-admin login
	resp, body = postJSON("/auth/login", map[string]string{"account": username, "password": "TestPass123"}, nil)
	log.Printf("POST /auth/login (user) -> status=%d body=%s", resp.StatusCode, string(body))
	if resp.StatusCode != 200 {
		log.Fatalln("user login failed")
	}
	var ur map[string]interface{}
	json.Unmarshal(body, &ur)
	userToken := ur["data"].(map[string]interface{})["token"].(string)
	userAuth := map[string]string{"Authorization": "Bearer " + userToken}

	// 5. Non-admin accessing admin menu should be 403
	resp, _ = doGet("GET", "/admin/menu", userAuth)
	log.Printf("GET /admin/menu as user -> status=%d", resp.StatusCode)
	if resp.StatusCode != 403 {
		log.Fatalln("non-admin access should be forbidden")
	}

	// 6. Admin list users
	resp, body = doGet("GET", "/admin/users?page=1&pageSize=10", auth)
	log.Printf("GET /admin/users -> status=%d body=%s", resp.StatusCode, string(body))
	if resp.StatusCode != 200 {
		log.Fatalln("list users failed")
	}

	// 7. Get user detail
	resp, body = doGet("GET", fmt.Sprintf("/admin/users/%d", userId), auth)
	log.Printf("GET /admin/users/%d -> status=%d body=%s", userId, resp.StatusCode, string(body))
	if resp.StatusCode != 200 {
		log.Fatalln("get user failed")
	}

	// 8. Reset user password
	resp, body = postJSON(fmt.Sprintf("/admin/users/%d/password", userId), map[string]interface{}{"newPassword": "NewPass123", "confirm": true}, auth)
	log.Printf("POST /admin/users/%d/password -> status=%d body=%s", userId, resp.StatusCode, string(body))
	if resp.StatusCode != 200 {
		log.Fatalln("reset password failed")
	}

	// 9. Login with new password
	resp, body = postJSON("/auth/login", map[string]string{"account": username, "password": "NewPass123"}, nil)
	log.Printf("POST /auth/login (user new pass) -> status=%d body=%s", resp.StatusCode, string(body))
	if resp.StatusCode != 200 {
		log.Fatalln("login with new password failed")
	}
	json.Unmarshal(body, &ur)
	userToken = ur["data"].(map[string]interface{})["token"].(string)
	userAuth = map[string]string{"Authorization": "Bearer " + userToken}

	// 10. Create a moment as normal user
	resp, body = postJSON("/api/moments", map[string]interface{}{"content": "测试内容 from admin smoke", "tags": []string{"smoke"}}, userAuth)
	log.Printf("POST /api/moments -> status=%d body=%s", resp.StatusCode, string(body))
	if resp.StatusCode != 200 {
		log.Fatalln("create moment failed")
	}
	var mr map[string]interface{}
	json.Unmarshal(body, &mr)
	post := mr["data"].(map[string]interface{})
	postId := int64(post["id"].(float64))

	// 11. Admin list contents
	resp, body = doGet("GET", "/admin/contents?page=1&pageSize=10", auth)
	log.Printf("GET /admin/contents -> status=%d body=%s", resp.StatusCode, string(body))
	if resp.StatusCode != 200 {
		log.Fatalln("list contents failed")
	}

	// 12. Get content detail
	resp, body = doGet("GET", fmt.Sprintf("/admin/contents/%d", postId), auth)
	log.Printf("GET /admin/contents/%d -> status=%d body=%s", postId, resp.StatusCode, string(body))
	if resp.StatusCode != 200 {
		log.Fatalln("get content detail failed")
	}

	// 13. Delete content
	resp, body = doGet("DELETE", fmt.Sprintf("/admin/contents/%d?confirm=true", postId), auth)
	log.Printf("DELETE /admin/contents/%d -> status=%d body=%s", postId, resp.StatusCode, string(body))
	if resp.StatusCode != 200 {
		log.Fatalln("delete content failed")
	}

	// 14. Public get moment should be not 200
	resp, _ = doGet("GET", fmt.Sprintf("/api/moments/%d", postId), nil)
	log.Printf("GET /api/moments/%d -> status=%d (expect not 200)", postId, resp.StatusCode)
	if resp.StatusCode == 200 {
		log.Fatalln("deleted content still visible")
	}

	// 15. Delete user
	resp, _ = doGet("DELETE", fmt.Sprintf("/admin/users/%d?confirm=true", userId), auth)
	log.Printf("DELETE /admin/users/%d -> status=%d", userId, resp.StatusCode)
	if resp.StatusCode != 200 {
		log.Fatalln("delete user failed")
	}

	// 16. Verify user deleted
	resp, _ = doGet("GET", fmt.Sprintf("/admin/users/%d", userId), auth)
	log.Printf("GET /admin/users/%d after delete -> status=%d", userId, resp.StatusCode)
	if resp.StatusCode == 200 {
		log.Fatalln("deleted user still accessible")
	}

	// Cleanup: remove test admin
	db.Where("username = ?", admin.Username).Delete(&models.Admin{})
	log.Println("Smoke test completed successfully and test admin removed")
}
