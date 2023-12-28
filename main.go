package main

import (
	"github.com/bitmagnet-io/bitmagnet/internal/app"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	//fmt.Printf("%#v\n", unidecode.Unidecode("我"))
	//println(unicode.IsLetter('1'))
	app.New().Run()
	//fmt.Printf("%#v\n", tsvector.Tokenize("hello 555 world2010昨 日 お 店に行ってサンドイッチを購入しました hôtel Київ ℂ 안녕하세요, 샌드위치 하나 먹어도 될까요? مرحباً، هل لي أن أتناول شطيرة؟ 你好，我可以吃三明治吗？ สวัสดีครับ ฉันสามารถกินแซนวิชได้ไหมครับ?"))
	//fmt.Printf("%#v\n", tsvector.Tokenize("مرحباً، هل لي أن أتناول شطيرة؟"))
	//fmt.Printf("%#v\n", tsvector.Tokenize("hello world"))
	//println(tsvector.ParseQuery("Hello world foo 你好，我可以吃!!!|(我可以吃三明治吗？店に行ってサンドイッチを購入しました) | (fat !cat) | mat"))
}
