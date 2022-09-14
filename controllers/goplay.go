package controllers

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dinirestyan/goplay-debugging/models"
	"github.com/dinirestyan/goplay-debugging/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

func NewGoplayController(db *gorm.DB, repo models.GoplayRepo) {
	_ = &GoplayController{
		DB:   db,
		repo: repo,
	}
}

type GoplayController struct {
	DB   *gorm.DB
	repo models.GoplayRepo
}

func generateJwtToken(userID uuid.UUID, email string) (string, error) {
	expiresAt := time.Now().Add(time.Hour * 1000).Unix()
	tk := &models.Token{
		UserID: userID,
		Email:  email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}
	// generate jwt auth
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		return "", errors.New("can't signed token")
	}

	return tokenString, nil
}

func (g *GoplayController) Login(ctx *gin.Context) {

	var (
		req  models.Login
		user models.User
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res := utils.Response("failed to bind json", nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data, err := g.repo.Login(req.Email)
	if err != nil {
		res := utils.Response("failed to get user", nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	if !utils.IsPasswordValid(data.Password, req.Password) {
		res := utils.Response("incorrect password", nil)
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	tokenString, err := generateJwtToken(user.ID, user.Email)
	if err != nil {
		res := utils.Response(err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	user.Token = tokenString
	res := utils.Response("logged in", user)
	ctx.JSON(http.StatusOK, res)
	return
}

func (g *GoplayController) Upload(ctx *gin.Context) {

	var file models.Upload

	if err := ctx.Bind(&file); err != nil {
		res := utils.Response("cannot bind request: "+err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	user_id, errClaim := utils.GetClaims(ctx)
	if errClaim != nil {
		res := utils.Response("You're not allowed", nil)
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	id, err := uuid.FromString(user_id)
	if err != nil {
		res := utils.Response(err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	for _, file := range file.Documents {
		fname, err := utils.UploadFile(ctx, "gp/test", file.Filename, file)
		if err != nil {
			res := utils.Response("cannot upload file: "+err.Error(), nil)
			ctx.JSON(http.StatusInternalServerError, res)
			return
		}

		record := models.Files{
			FileName: fname,
			Uploader: id,
		}

		errUpload := g.DB.Table("files").Create(&record).Error
		if errUpload != nil {
			res := utils.Response("cannot store file: "+err.Error(), nil)
			ctx.JSON(http.StatusInternalServerError, res)
			return
		}
	}

	res := utils.Response("success", nil)
	ctx.JSON(http.StatusOK, res)
	return
}

func (g *GoplayController) GetFileLists(ctx *gin.Context) {

	var (
		data     []models.Files
		totalRow int
	)

	condition := ctx.DefaultQuery("search", "")
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("limit", "10")
	convPage, _ := strconv.Atoi(page)
	convPerPage, _ := strconv.Atoi(perPage)
	offset := (convPage - 1) * convPerPage

	_, errClaim := utils.GetClaims(ctx)
	if errClaim != nil {
		res := utils.Response("You're not allowed", nil)
		ctx.JSON(http.StatusUnauthorized, res)
		return
	}

	query := g.DB
	if condition != "" {
		query = query.Where("file_name ilike '%" + condition + "%'")
	}
	errGet := query.Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&data).
		Error

	if errGet != nil {
		res := utils.Response(errGet.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	errCount := query.
		Model(&data).
		Count(&totalRow).
		Error

	if errCount != nil {
		res := utils.Response(errGet.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	for i := 0; i < len(data); i++ {
		data[i].FileName = os.Getenv("GOOGLE_CLOUD_PATH") +
			os.Getenv("GOOGLE_BUCKET") + "/" + data[i].FileName
	}

	res := utils.PaginationResponse(totalRow, page, perPage, data)
	ctx.JSON(http.StatusOK, res)
	return
}

func Hahaha(test string) string {
	return "oke"
}
