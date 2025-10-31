package main

import (
	"encoding/json" //کار با json
	"fmt"           //چاپ متن و خروجی
	"net/http"      //ساخت و مدیریت سرور
	"os"            //کار با سیستم عامل و خواندن و نوشتن در فایل
	"path/filepath" //کار با مسیر فایل ها
	"strings"       //کار با رشته ها
	"sync"          //همگام سازی
	"time"          //کار با زمان
)

// ایجاد یک ساختار برای نگه داری داده ها و همگام سازی
type storing struct {
	Infirmation map[string]json.RawMessage
	Lock        sync.RWMutex
}

// تعریف ساختار پیشفرض برای درخواست ها
var request struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

// تعریف متغیرهای سراسری
var Memory = storing{Infirmation: make(map[string]json.RawMessage)}
var filename = "/app/data/Store.json"

// تابع نوشتن در فایل که داده ها رو به صورت json  میکنه و دسترسی رو نیز محدود میکند و نوشتن رو برای همه امکان پذیر نمیکند
func Writing() {
	bytes, err := json.MarshalIndent(Memory.Infirmation, "", "  ") //داده ها رو به صورت json  در میاره و برسی میکنه که خطایی در تبدیل وجود داره یا نه
	if err != nil {
		fmt.Println("Error Writing:", err)
		return
	}
	//اینجا مطمئن میشیم که دایرکتوری وجود داره یا نه
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		fmt.Println("Error Writing:", err)
		return
	}
	//در اینجا یک فایل موقت برای نوشتن داده هاایجاد کرده و داده ها رو توی اون مینویسیم
	tr := filename + ".tmp"
	if err := os.WriteFile(tr, bytes, 0644); err != nil {
		fmt.Println("Error Writing:", err)
		return
	}
	//نام فایل موقت رو به فایل اصلی تغییر میدهیم
	if err := os.Rename(tr, filename); err != nil {
		fmt.Println("Error Writing:", err)
		return
	}
	//زمان دقیق بارگزاری داده ها
	fmt.Println("[" + time.Now().Format(time.RFC3339) + "] data loaded")
}

// فایل json رو میخواند و در صورت وجود داده ها آن ها رو در حافظه بارگزاری میکند
func Reading() {
	bytes, err := os.ReadFile(filename)
	//اگر فایلی وجود نداشاه باشه یک حافظه ی خالی میسازیم
	if err != nil {
		Memory.Infirmation = make(map[string]json.RawMessage)
		Writing()
		return
	}
	//و اگر فایل وجود داشت داده ها رو در حافظه بارگزاری میکنیم
	if err := json.Unmarshal(bytes, &Memory.Infirmation); err != nil {
		fmt.Println("Error Reading:", err)
		Memory.Infirmation = make(map[string]json.RawMessage)
	}
}

// همه ی داده های ذخیره شده  رو با این تابع به ما نمایش میدهد
func showalldata(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { //نوع درخواست برسی میشه و فقط درخواست get  مجاز است
		http.Error(w, "Not allowed(Errore 405)", http.StatusMethodNotAllowed)
		return
	}
	// حافظه رو قفل میکنیم تا امکان نوشتن همزمان نباشد
	Memory.Lock.RLock()
	defer Memory.Lock.RUnlock()
	// مقدار رو به صورت json  برمیگردونیم
	w.Header().Set("Content-Type", "application/json")
	if len(Memory.Infirmation) == 0 {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{}")) //و یک شی خالی json  برمیگرداند
		return
	}
	json.NewEncoder(w).Encode(Memory.Infirmation)
	if err := json.NewEncoder(w).Encode(Memory.Infirmation); err != nil {
		http.Error(w, "Error http ", http.StatusInternalServerError)
		fmt.Println("Error http", err)
	}
}

// تابع ذخیره ساز که زمان که درخواست putدریافت کمیکند اجرا میشود
func Storage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Not allowed(Errore 405)", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	//در اینجا نوع محتوا را برسی میکنیم که خالیست  یا با application/json  شروع میشود یا نه
	HG := r.Header.Get("content-type")
	if HG == "" || !strings.HasPrefix(HG, "application/json") {
		http.Error(w, "Not application/json(Errore 415)", http.StatusUnsupportedMediaType)
		return
	}
	//در اینجا محتوای درخواست رو به ساختار از پیش تعریف شده تبدیل میکنیم
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Not JSON", http.StatusBadRequest)
		return
	}
	//در این بخش صحت کلید برسی میشه
	if strings.TrimSpace(request.Key) == "" {
		http.Error(w, "Key is required", http.StatusBadRequest)
		return
	}
	//بخش همگام سازی که در آن حافظه رو فقل میکنیم تا امکان نوشتن همزمان رو از بین ببریم
	Memory.Lock.Lock()
	Memory.Infirmation[request.Key] = request.Value
	Writing()
	Memory.Lock.Unlock()
	w.WriteHeader(http.StatusOK)
}

// تابع بازیابی که زمانی  که درخواست get  میاد اجرا میشه
func Restore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Not allowed(Errore 405)", http.StatusMethodNotAllowed)
		return
	}
	// در اینجا کلید رو با رعایت همزمانی استخراج میکنیم
	key := r.URL.Path[len("/objects/"):]
	Memory.Lock.RLock()
	defer Memory.Lock.RUnlock()
	value, ok := Memory.Infirmation[key]
	//در صورت نبود کلید وارد این قسمت میشیم
	if !ok {
		http.Error(w, "ّFalse (Error 404)", http.StatusNotFound)
		return
	}
	//اگر بود هم مقدارش رو به صورت json  برمیگردونیم
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(value)
}

// تعریف تابع اصلی برای ذخیره و بازیابی داده ها
func main() {
	//در صورت وجود داده ها اون ها رو از حافظه میخونیم
	Reading()
	//پورت رو به صورت متغیر محیطی میخواند و در صورت نبود آن از پورت پیشفرض استفاده میکنه
	port := os.Getenv("PORT")
	if port == "" {
		port = "18080"
	}
	http.HandleFunc("/objects", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut: //در خواست put
			Storage(w, r)
		case http.MethodGet: //درخواست get برای نمایش همه داده ها
			showalldata(w, r)
		default:
			http.Error(w, "Not allowed (Error 405)", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/objects/", Restore) //برای بازیابی داده
	fmt.Printf("Server listen to://127.0.0.1:%s\n", port)
	//گوش دادن سرور روی پورت مشخص شده
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("server error:", err)
	}
}
