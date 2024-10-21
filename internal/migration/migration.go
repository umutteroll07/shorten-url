package migration

import (
	"github.com/umutteroll07/ShortenURL/app/model"
	"github.com/umutteroll07/ShortenURL/internal/database"
)

func init(){
	database.Connection().AutoMigrate(&model.Url{})
}