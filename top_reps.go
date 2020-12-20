package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocarina/gocsv"
)

var (
	tempRepList []Representative
)

func getTopReps(c *gin.Context) {
	reps := []*Representative{}

	in, err := os.Open("./test_data/top_house_reps.csv")
	if err != nil {
		panic(err)
	}
	defer in.Close()
	if err := gocsv.UnmarshalFile(in, &reps); err != nil {
		panic(err)
	}
	for _, rep := range reps {
		tempRepList = append(tempRepList, addRepToMap(rep))
	}

	in, err = os.Open("./test_data/top_senate_reps.csv")
	if err != nil {
		panic(err)
	}
	defer in.Close()
	if err := gocsv.UnmarshalFile(in, &reps); err != nil {
		panic(err)
	}
	for _, rep := range reps {
		tempRepList = append(tempRepList, addRepToMap(rep))
	}

	msg := map[string]interface{}{"Status": "Ok", "users_rep_list": tempRepList}
	c.JSON(http.StatusOK, msg)
}
