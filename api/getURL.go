package api

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/labstack/echo/v4"
)

func (a *API) GetURL(e echo.Context, shortUrl string) error {
	url := Url{ShortUrl: shortUrl}

	response, err := a.db.Client.GetItem(e.Request().Context(), &dynamodb.GetItemInput{
		Key: url.GetKey(), TableName: aws.String(a.db.TableName),
	})
	if err != nil {
		return e.JSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("failed: %v", err),
		})
	}

	var longUrl string

	err = attributevalue.Unmarshal(response.Item["LongUrl"], &longUrl)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, Error{
			Message: fmt.Sprintf("failed: %v", err),
		})
	}

	return e.Redirect(http.StatusFound, longUrl)
}
