# ALTF4_TEST
Основной код в main.go  
Использованная документация:  
https://bablofil.ru/binance-webscokets/  
https://github.com/binance-exchange/binance-official-api-docs/blob/mster/rest-api.md  

Endpoint  
  wss://stream.binance.com:9443/ws/blabla@depth  
Ответ  
  {
  "e": "depthUpdate", // Тип события
  "E": 123456789,     // Время события
  "s": "BNBBTC",      // Пара
  "U": 157,           // Первое ID события
  "u": 160,           // Последнее ID события
  "b": [              // Покупки
    [
      "0.0024",       // Цена
      "10",         // Объем
      []              // Не актуально
    ]
  ],
  "a": [              // Продажи
    [
      "0.0026",       // Цена
      "100",          // Объем
      []              // Не актуально
    ]
  ]
}  
TODO   
  Можно вместо json использовать easyjson  
  Вместо gorilla/websocket использовать gobwas/ws  
  Вместо waitgroup использовать context  
  Вместо структур проще кажется map  
  
  
