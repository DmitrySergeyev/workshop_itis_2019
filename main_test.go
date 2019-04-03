package main

import (
	"math"
	"net"
	"testing"
	"time"
)

var messages = []string{
	"power_off",
	"power_on",
	"service_error",
	"service_ok",
	"security",
}

const UDPAddr = "127.0.0.1:8081"

func TestUDPSocket(t *testing.T) {
	var i uint8

	// Получение количества элементов в массиве
	var len = uint8(len(messages))

	// Создание тика, который сработает через 1 секунду
	var ticker = time.Tick(1 * time.Second)

	for {
		select {
		case <-ticker:

			// Создание подключения и обработка ошибки
			conn, err := net.Dial("udp4", UDPAddr)
			if err != nil {
				t.Error(err)
				return
			}

			// Отправка сообщения и закрытие соединения
			conn.Write([]byte(messages[i%len]))
			conn.Close()

			if i >= math.MaxUint8 {
				t.Error("Counter was too large")
				return
			}

			// Увеличиваем счетчик и создаем новый тик
			i++
			ticker = time.Tick(1 * time.Second)
		}
	}

}
