package main

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

func main() {
	log.Printf("Authorizing bot...")

	bot, err := tgbotapi.NewBotAPI("{your_token}")
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		var msg tgbotapi.MessageConfig

		if update.Message.Text == "/help" || update.Message.Text == "/start" {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, getHelpText());
			bot.Send(msg);
		} else if (update.Message.Text == "/currency") {
			currenciesData := getCurrenciesCollection()
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, exchangeRatesToString(*currenciesData));
			bot.Send(msg);
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "I didn't get you, so I think you need currencies.")
			bot.Send(msg);
			currenciesData := getCurrenciesCollection()
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, exchangeRatesToString(*currenciesData));
			bot.Send(msg);
		}
	}
}

func getHelpText() string {
	return "/help - Help \n/currency - List of all currencies"
}

func getCurrenciesCollection() *CurrenciesCollection {

	url := "https://api.privatbank.ua/p24api/pubinfo?json&exchange&coursid=5";

	resp, err := http.Get(url)

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	currencyData, err := parseCurrencies([]byte(body))

	return currencyData
}

func parseCurrencies(body []byte) (*CurrenciesCollection, error) {
	var s = new(CurrenciesCollection)
	err := json.Unmarshal(body, &s)
	if(err != nil){
		fmt.Println("whoops:", err)
	}
	return s, err
}

type CurrenciesCollection []Currency


type Currency struct {
	Ccy string `json:"ccy"`
	BaseCcy string `json:"base_ccy"`
	Buy string `json:"buy"`
	Sale string `json:"sale"`
}

func exchangeRatesToString(currenciesCollection CurrenciesCollection) string {
	var message string;

	for i := 0; i< len(currenciesCollection); i++ {
		currency := currenciesCollection[i];
		fmt.Println(currency.Sale);
		message += currency.Ccy + "\n";
		message += "- sale " + currency.Sale + " " + currency.BaseCcy + "\n"
		message += "- buy " + currency.Buy + " " + currency.BaseCcy + "\n"
	}

	return message
}