package pattern

import (
	"fmt"
	"time"
)

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

/*
Шаблон Стратегия определяет семейство похожих алгоритмов и свой объект для каждого из них.
Это позволяет клиенту выбирать подходящий алгоритм на этапе выполнения кода.
Клиент освобождён от деталей реализации алгоритмов. Это помогает улучшать алгоритмы в объекте Стратегия или добавлять новые, не требуя изменений клиентского кода. С другой стороны, клиент должен знать, чем различаются существующие алгоритмы, чтобы выбрать подходящий вариант.
Паттерн Стратегия используется, когда:
 - нужно применять разные варианты одного и того же алгоритма;
 - нужно выбирать алгоритм во время выполнения программы;
 - нужно скрывать детали реализации алгоритмов.
Примеры использования:
 - сортировка данных(выбор алгоритма в зависимости от имеющихся данных)
 - выбор алгоритма шифрования
 - определение стратегии поведения в зависимости от разных факторов
*/

// Пример с кешированием данных. Есть две стратегии ротации кеша:
// LRU — вычищаем элементы, которые использовались давно.
// FIFO — удаляем элементы, которые были созданы раньше остальных.

type EvictionAlgo interface {
	Evict(c *Cache)
}

type Fifo struct{}

func (l *Fifo) Evict(c *Cache) {
	curMinTime := time.Now()
	curKey := ""
	for k, v := range c.storage {
		if !v.created.After(curMinTime) {
			curMinTime = v.created
			curKey = k
		}
	}
	delete(c.storage, curKey)
	c.capacity--
	fmt.Println("Evicted by fifo strategy: ", curKey)
}

type Lru struct{}

func (l *Lru) Evict(c *Cache) {
	curMinTime := time.Now()
	curKey := ""
	for k, v := range c.storage {
		if !v.lastUsed.After(curMinTime) {
			curMinTime = v.lastUsed
			curKey = k
		}
	}
	delete(c.storage, curKey)
	c.capacity--
	fmt.Println("Evicted by lru strategy: ", curKey)
}

type Element struct {
	data     string
	created  time.Time
	lastUsed time.Time
}

type Cache struct {
	storage      map[string]Element
	evictionAlgo EvictionAlgo
	capacity     int
	maxCapacity  int
}

func InitCache(e EvictionAlgo) *Cache {
	storage := make(map[string]Element)
	return &Cache{
		storage:      storage,
		evictionAlgo: e,
		capacity:     0,
		maxCapacity:  2,
	}
}

// SetEvictionAlgo определяет алгоритм освобождения памяти.
func (c *Cache) SetEvictionAlgo(e EvictionAlgo) {
	c.evictionAlgo = e
}

func (c *Cache) Add(key, value string) {
	if c.capacity == c.maxCapacity {
		c.Evict()
	}
	c.capacity++
	c.storage[key] = Element{
		data:     value,
		created:  time.Now(),
		lastUsed: time.Now(),
	}
}

func (c *Cache) Get(key string) string {
	e, ok := c.storage[key]
	if !ok {
		return ""
	}
	e.lastUsed = time.Now()
	c.storage[key] = e
	return e.data
}

func (c *Cache) Evict() {
	c.evictionAlgo.Evict(c)
}

func main() {
	lru := &Lru{}
	cache := InitCache(lru)
	cache.Add("first", "data1")
	time.Sleep(time.Second)
	cache.Add("second", "data2")
	time.Sleep(time.Second)
	fmt.Println(cache.Get("first"))
	time.Sleep(time.Second)
	cache.Add("third", "data3")
	fifo := &Fifo{}
	cache.SetEvictionAlgo(fifo)
	cache.Add("fourth", "data4")
}
