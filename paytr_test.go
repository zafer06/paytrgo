package paytrgo

import (
	"os"
	"testing"
	"time"
)

func TestCheckHash(t *testing.T) {
	var hashData = map[string]string{
		"hash":         "c37hFnidanige1rWAi+G6C18Eya+jxCQaJ2jQ2AdiDY=",
		"merchant_oid": "20190925T113214",
		"status":       "success",
		"total_amount": "3588",
	}

	ok, hash := CheckHash(hashData, os.Getenv("PAYTR_MERCHANT_KEY"), os.Getenv("PAYTR_MERCHANT_SALT"))

	if !ok {
		t.Errorf("Hash degerleri uyumsuz. Beklenen: %s, Gelen: %s",
			hashData["hash"], hash)
	}
}

func TestPayTRToken(t *testing.T) {
	var valuesList = map[string]string{
		"merchantID":      os.Getenv("PAYTR_MERCHANT_ID"),
		"merchantKey":     os.Getenv("PAYTR_MERCHANT_KEY"),
		"merchantSalt":    os.Getenv("PAYTR_MERCHANT_SALT"),
		"email":           "info@yahoo.com",
		"paymentAmount":   "3588",
		"merchantOid":     time.Now().Format("150405"),
		"userIP":          "192.168.1.56",
		"userName":        "musteriadivesoyadi",
		"userAddress":     "musteriadres",
		"userPhone":       "musteritelefon",
		"noInstallment":   "0",                                         // Taksit yapılmasını istemiyorsanız, sadece tek çekim sunacaksanız 1 yapın
		"maxInstallment":  "0",                                         // Sayfada görüntülenecek taksit adedi sıfır (0) tüm taksitler
		"currency":        "TL",                                        // İşlem yapılacak para birimi, boş bırakılırsa TR kabul edilir.
		"testMode":        "0",                                         // Mağaza canlı modda iken test işlem yapmak için 1 olarak gönderilebilir.
		"debugOn":         "1",                                         // Hata mesajlarının ekrana basılması için entegrasyon ve test sürecinde 1 olarak bırakın. Daha sonra 0 yapabilirsiniz.
		"timeoutLimit":    "30",                                        // İşlem zaman aşımı süresi - dakika cinsinden
		"merchantFailURL": "http://www.siteniz.com/odeme_hata.php",     // Ödeme sürecinde beklenmedik bir hata oluşması durumunda müşterinizin yönlendirileceği sayfa
		"merchantOkURL":   "http://www.siteniz.com/odeme_basarili.php", // Başarılı ödeme sonrası müşterinizin yönlendirileceği sayfa
	}

	var productList = []Product{
		{"Örnek ürün 1", "18.00", 1},
		{"Örnek ürün 2", "33.25", 2},
		{"Örnek ürün 3", "45.42", 1},
	}

	var res = GetToken(valuesList, productList)

	if res.Status != "success" {
		t.Errorf("Hazirlanan token hatali gorunuyor. Hata: %s", res.Reason)
	}
}
