package controllers

import (
	"aurora-im/model"
	"aurora-im/config"

	// "fmt"

	"github.com/gin-gonic/gin"
)

func UpdateProfile(c *gin.Context) {
	uid := c.GetString("uid")

	var req struct {
        Username string `json:"username"`
        Avatar string `json:"avatar"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"msg": "Invalid Input"})
        return
    }

	var user model.User

	user.Username = req.Username
	user.Avatar   = req.Avatar

	result := config.DB.
		Model(&model.User{}).
		Where("id = ?", uid).
		Select("username", "avatar").
		Updates(user)

	if result.Error != nil {
		c.JSON(500, gin.H{"msg": "search failed"})
		return
	}    
}

// func UploadAvatar(c *gin.Context) {
// 	// 1. get user id from JWT
// 	uid := c.GetString("uid")

// 	// 2. get file
// 	file, err := c.FormFile("avatar")
// 	if err != nil {
// 		c.JSON(400, gin.H{"error": "avatar file required"})
// 		return
// 	}

// 	// 3. validate size (max 2MB)
// 	if file.Size > 2*1024*1024 {
// 		c.JSON(400, gin.H{"error": "file too large"})
// 		return
// 	}

// 	// 4. validate type
// 	ext := strings.ToLower(filepath.Ext(file.Filename))
// 	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
// 		c.JSON(400, gin.H{"error": "invalid file type"})
// 		return
// 	}

// 	// 5. create upload dir
// 	dir := "uploads/avatar"
// 	os.MkdirAll(dir, os.ModePerm)

// 	// 6. generate filename
// 	filename := fmt.Sprintf("%d_%d%s", userID, time.Now().Unix(), ext)
// 	path := filepath.Join(dir, filename)

// 	// 7. save file
// 	if err := c.SaveUploadedFile(file, path); err != nil {
// 		c.JSON(500, gin.H{"error": "upload failed"})
// 		return
// 	}

// 	// 8. generate access URL
// 	avatarURL := "/static/avatar/" + filename

// 	// 9. update database
// 	db.Model(&User{}).
// 		Where("id = ?", userID).
// 		Update("avatar_url", avatarURL)

// 	c.JSON(200, gin.H{
// 		"avatar_url": avatarURL,
// 	})
// }
