package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"
)

var money = atomic.Int64{} //usd   
var bank = atomic.Int64{}    //usd

func payHandler(w http.ResponseWriter, r *http.Request){
    httpRequestBody, err:= io.ReadAll(r.Body) //берем тело и вычитываем его полностью, блягодаря этой функции
    if err != nil {
        fmt.Println("fail to read http body:", err.Error())
        return
    }

    httpRequestBodyString := string(httpRequestBody)
    paymentAmount, err:=  strconv.Atoi(httpRequestBodyString)
    if err != nil{
        fmt.Println("Error converting to a number", err)
        return 
    }

    if money.Load() - int64(paymentAmount) >= 0{
        money.Add(-int64(paymentAmount))
        fmt.Println("The payment was successful:", money.Load())
    } else {
        fmt.Println("There are not enough funds to make the payment!")
    }

    
}

func saveHandler(w http.ResponseWriter, r *http.Request){
     httpRequestBody, err :=io.ReadAll(r.Body)
     if err != nil{
        fmt.Println("fail to read http request body:", err)
        return 
     }
     httpRequestBodyString:=string(httpRequestBody)
     saveAmount, err:= strconv.Atoi(httpRequestBodyString) // переводим в число
     if err != nil {
        fmt.Println("Error converting to a number", err)
        return 
     }

     if money.Load() >= int64(saveAmount){
        money.Add(-int64(saveAmount))

        bank.Add(int64(saveAmount))
        fmt.Println("New variable value:", bank.Load())
        fmt.Println("New variable value:", money.Load())

     } else {
        fmt.Println("There is not enough money to deposit into the account!")
     }

     
}

func main(){
    fmt.Println("The program is running")

    money.Add(1000)

    http.HandleFunc("/pay", payHandler)
    http.HandleFunc("/save", saveHandler)

    err:=http.ListenAndServe(":9091", nil) 
    if err != nil {
        fmt.Println("Error", err.Error())
    }
}
