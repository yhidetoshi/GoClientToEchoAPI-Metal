package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mackerelio/mackerel-client-go"
)

var (
	//url = "http://localhost:1323/metal"
	url    = os.Getenv("APIURL")
	mkrKey = os.Getenv("MKRKEY")
	client = mackerel.NewClient(mkrKey)
	id = os.Getenv("ID")
	pw = os.Getenv("PW")
)

const (
	serviceName = "Metal"
	timezone    = "Asia/Tokyo"
	offset      = 9 * 60 * 60
)

// Metal set Metal value
type Metal struct {
	Date     time.Time `json:"time"`
	GoldInfo GoldInfo  `json:"goldInfo"`
	Platinum Platinum  `json:"platinum"`
}

// GoldInfo set Gold value
type GoldInfo struct {
	RetailTax   int `json:"retailTax"`
	PurchaseTax int `json:"purchaseTax"`
}

// Platinum set Platinum value
type Platinum struct {
	RetailTax   int `json:"retailTax"`
	PurchaseTax int `json:"purchaseTax"`
}

func main() {
	lambda.Start(Handler)
}

// func main() {

// Handler lambda
func Handler() {

	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(id, pw)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	metal := &Metal{}

	err = json.NewDecoder(resp.Body).Decode(metal)
	if err != nil {
		fmt.Println(err)
	}

	jst := time.FixedZone(timezone, offset)
	nowTime := time.Now().In(jst)

	mkrErr := PostValuesToMackerel(metal.GoldInfo.RetailTax, metal.GoldInfo.PurchaseTax, metal.Platinum.RetailTax, metal.Platinum.PurchaseTax, nowTime)
	if mkrErr != nil {
		fmt.Println(mkrErr)
	}
}

// PostValuesToMackerel Post Metrics to Mackerel
func PostValuesToMackerel(goldRetailTax int, goldPurchaseTax int, platinumRetailTax int, platinumPurchaseTax int, nowTime time.Time) error {
	// Post Gold metrics
	errGold := client.PostServiceMetricValues(serviceName, []*mackerel.MetricValue{
		&mackerel.MetricValue{
			Name:  "Gold.retail_tax",
			Time:  nowTime.Unix(),
			Value: goldRetailTax,
		},
		{
			Name:  "Gold.purchase_tax",
			Time:  nowTime.Unix(),
			Value: goldPurchaseTax,
		},
	})
	if errGold != nil {
		fmt.Println(errGold)
	}

	// Post Platinum metrics
	errPlatinum := client.PostServiceMetricValues(serviceName, []*mackerel.MetricValue{
		&mackerel.MetricValue{
			Name:  "Platinum.retail_tax",
			Time:  nowTime.Unix(),
			Value: platinumRetailTax,
		},
		{
			Name:  "Platinum.purchase_tax",
			Time:  nowTime.Unix(),
			Value: platinumPurchaseTax,
		},
	})
	if errPlatinum != nil {
		fmt.Println(errPlatinum)
	}

	return nil
}
