package paytrgo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// Product model
type Product struct {
	Name     string
	Price    string
	Quantity int
}

// PayTRResponse model
type PayTRResponse struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
	Token  string `json:"token"`
}

// GetToken gerekli bilgileri alıp kullanıcı onayından sonra
// kullanım için bir token bilgisi verir.
func GetToken(list map[string]string, productList []Product) PayTRResponse {
	var basket = getBasket(productList)

	var hashStr = getValuesString(list, basket)
	var paytrToken = createPayTrToken(hashStr, list["merchantSalt"], list["merchantKey"])

	var paytrValues = url.Values{
		"merchant_id":       {list["merchantID"]},
		"user_ip":           {list["userIP"]},
		"merchant_oid":      {list["merchantOid"]},
		"email":             {list["email"]},
		"payment_amount":    {list["paymentAmount"]},
		"user_basket":       {basket},
		"debug_on":          {list["debugOn"]},
		"no_installment":    {list["noInstallment"]},
		"max_installment":   {list["maxInstallment"]},
		"paytr_token":       {paytrToken},
		"user_name":         {list["userName"]},
		"user_address":      {list["userAddress"]},
		"user_phone":        {list["userPhone"]},
		"merchant_ok_url":   {list["merchantOkURL"]},
		"merchant_fail_url": {list["merchantFailURL"]},
		"timeout_limit":     {list["timeoutLimit"]},
		"currency":          {list["currency"]},
		"test_mode":         {list["testMode"]},
	}

	return connect(paytrValues)
}

// CheckHash paytr tarafından gönderilen bildirim mesajını doğrulamak
// icin hash doğrulaması yapar.
func CheckHash(paytrData map[string]string, merchantKey string, merchantSalt string) (bool, string) {
	h := hmac.New(sha256.New, []byte(merchantKey))
	h.Write([]byte(paytrData["merchant_oid"]))
	h.Write([]byte(merchantSalt))
	h.Write([]byte(paytrData["status"]))
	h.Write([]byte(paytrData["total_amount"]))
	var hash = base64.StdEncoding.EncodeToString(h.Sum(nil))

	var check = false
	if hash == paytrData["hash"] {
		check = true
	}

	return check, hash
}

// Ürünleri uygun formatta bir dizge haline getirir.
func getBasket(productList []Product) string {
	var urunler string
	for _, p := range productList {
		urunler += getEncodedProduct(p.Name, p.Price, p.Quantity)
		urunler += ","
	}
	urunler = urunler[0 : len(urunler)-1]

	var encodedBasket = "[" + urunler + "]"
	return base64.StdEncoding.EncodeToString([]byte(encodedBasket))
}

// Ürün bilglerindeki Türkçe karakterleri encode eder.
func getEncodedProduct(name string, price string, quantity int) string {
	return "[" +
		strconv.QuoteToASCII(name) + "," +
		strconv.QuoteToASCII(price) + "," +
		strconv.Itoa(quantity) + "]"
}

// Gerekli bilgileri tek bir dizge haline getirir
func getValuesString(list map[string]string, basket string) string {
	return list["merchantID"] +
		list["userIP"] +
		list["merchantOid"] +
		list["email"] +
		list["paymentAmount"] +
		basket +
		list["noInstallment"] +
		list["maxInstallment"] +
		list["currency"] +
		list["testMode"]
}

// Kimlik dogrulama icin gerekli token dizesini oluşturur.
func createPayTrToken(hashStr string, merchantSalt string, merchantKey string) string {
	h := hmac.New(sha256.New, []byte(merchantKey))
	h.Write([]byte(hashStr))
	h.Write([]byte(merchantSalt))
	var encodedHash = base64.StdEncoding.EncodeToString(h.Sum(nil))

	//log.Println("PayTR icin hmac-base64 hash degeri: ", encodedHash)

	return encodedHash
}

// PayTR api uzerinden bir istek gonderir.
func connect(paytrValues url.Values) PayTRResponse {
	res, err := http.PostForm("https://www.paytr.com/odeme/api/get-token", paytrValues)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	var response PayTRResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("PayTR istek durumu: ", res.Status)
		log.Println("PayTR token istek sonucu: ", string(body))
	}

	return response
}
