// Модуль delay предоставляет задержку.
package delay

import (
	"time"
)

// NewDelay возвращает функцию с замыканием , которая увичивает задержку.
func NewDelay() func() time.Duration {
	attempt := 0
	delay := 1 * time.Second
	delayIncrease := 2 * time.Second
	return func() time.Duration {
		attempt++
		if attempt == 1 {
			return delay
		}
		delay += delayIncrease
		return delay
	}
}
