package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"strconv"
	"time"
)

type Loan struct {
	//gorm.Model
	ID uint				`gorm:"primary_key" json:"id"`
	Date time.Time		`json:"date"`
	Ktp string			`json:"ktp"`
	BirthDate time.Time `json:"birth_date"`
	Gender string		`json:"gender"`
	Name string			`json:"name"`
	Amount int			`json:"amount"`
	Period int			`json:"period"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type TrackRecord struct {
	Count int		`json:"count"`
	Summary int		`json:"summary"`
	Avg int			`json:"avg"`
}

type Installment struct {
	MonthInstallment int	`json:"month_installment"`
	DueDate time.Time		`json:"due_date"`
	Capital float32			`json:"capital"`
	Interest float32		`json:"interest"`
	Total float32			`json:"total"`
}

var TimeNow = time.Now()
var zone, _  = TimeNow.Zone()

func AllLoan(w http.ResponseWriter, r *http.Request){
	defer RecoverFunc(w)
	db, err := gorm.Open("mysql", "root:password@/tunaiku_test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("fail to connect to database")
	}
	defer db.Close()

	var loans []Loan
	db.Find(&loans)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&loans)
}

func NewLoan(w http.ResponseWriter, r *http.Request){
	defer RecoverFunc(w)
	db, err := gorm.Open("mysql", "root:password@/tunaiku_test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("fail to connect to database")
	}
	defer db.Close()

	t := time.Now()
	zone, _ := t.Zone()

	date := r.FormValue("date")
	date_time, err := time.Parse(time.RFC3339, date+"T00:00:00"+zone+":00")
	ktp := r.FormValue("ktp")
	birth_date := r.FormValue("birth_date")
	birth_date_time, err := time.Parse("2006-01-02", birth_date)
	gender := r.FormValue("gender")
	name := r.FormValue("name")
	amount := r.FormValue("amount")
	amount_int, _ := strconv.Atoi(amount)
	period := r.FormValue("period")
	period_int, _ := strconv.Atoi(period)

	// Validasi Data KTP
	var ktp_dob string
	var dob string
	var ktp_mob string
	var yob string
	var ktp_yob string
	var mob string
	ktp_dob += string(ktp[6])
	ktp_dob += string(ktp[7])
	dob += string(birth_date[8])
	dob += string(birth_date[9])
	ktp_mob += string(ktp[8])
	ktp_mob += string(ktp[9])
	mob += string(birth_date[5])
	mob += string(birth_date[6])
	ktp_yob += string(ktp[10])
	ktp_yob += string(ktp[11])
	yob += string(birth_date[2])
	yob += string(birth_date[3])
	dob_int, _ := strconv.Atoi(dob)
	ktp_dob_int, _ := strconv.Atoi(ktp_dob)
	ktp_mob_int, _ := strconv.Atoi(ktp_mob)

	if ktp_mob_int > 12{
		panic("data ktp tidak valid")
	} else if mob != ktp_mob{
		panic("data ktp tidak valid")
	} else if ktp_yob != yob{
		panic("data ktp tidak valid")
	} else if gender == "Male"{
		if dob_int != ktp_dob_int{
			panic("data ktp tidak valid")
		}
	} else if gender == "Female"{
		if ktp_dob_int != (dob_int+31){
			panic("data ktp tidak valid")
		}
	}

	insert := db.Create(&Loan{Date:date_time, Ktp:ktp, BirthDate:birth_date_time, Gender:gender, Name:name, Amount:amount_int, Period:period_int})
	if insert.Error != nil{
		CustomResponse(400, "fail to insert", w)
	}else {
		CustomResponse(201, "loan created succesfully", w)
	}
}

func ListLoan(w http.ResponseWriter, r *http.Request){
	defer RecoverFunc(w)
	db, err := gorm.Open("mysql", "root:password@/tunaiku_test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("fail to connect to database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	time_request := r.FormValue("date")
	time_parse, _ := time.Parse(time.RFC3339, time_request+"T00:00:00"+zone+":00")

	summary := 0
	counter := 0
	for i:=0;i<7;i++{
		var loan Loan
		db.Where("date = ?", time_parse).Where("name = ?", name).First(&loan)
		time_parse = time_parse.AddDate(0, 0, -1)
		if loan.ID != 0{
			summary+=loan.Amount
			counter++
		}
	}
	avg := summary/7

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(TrackRecord{Count:counter, Summary:summary, Avg:avg})
}

func InstallmentLoan(w http.ResponseWriter, r *http.Request){
	defer RecoverFunc(w)
	date := r.FormValue("date")
	amount := r.FormValue("amount")
	period := r.FormValue("period")
	date_parse, _ := time.Parse("2006-01-02", date)
	amount_int, _ := strconv.Atoi(amount)
	period_int, _ := strconv.Atoi(period)

	month_installment := 1
	if period_int == 12 || period_int == 18{
		interest := (float32(amount_int)/100)*1.68
		capital := amount_int/period_int
		total := interest+float32(capital)
		var installments []Installment
		for i:=0;i<period_int;i++{
			installment := Installment{MonthInstallment:month_installment, DueDate:date_parse, Capital:float32(capital), Interest:interest, Total:total}
			installments = append(installments, installment)
			month_installment++
			date_parse = date_parse.AddDate(0, 1, 0)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(installments)
	}else if period_int == 24 || period_int == 30 || period_int == 36{
		interest := (float32(amount_int)/100)*1.59
		capital := amount_int/period_int
		total := interest+float32(capital)
		var installments []Installment
		for i:=0;i<period_int;i++{
			installment := Installment{MonthInstallment:month_installment, DueDate:date_parse, Capital:float32(capital), Interest:interest, Total:total}
			installments = append(installments, installment)
			month_installment++
			date_parse = date_parse.AddDate(0, 1, 0)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(installments)
	}

	panic("jumlah periode salah, silahkan pilih 12, 18, 24, 50 atau 36")
}
