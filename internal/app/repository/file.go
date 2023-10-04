package repository

import (
	"strconv"
)

// LoadFromFile Загрузить список ссылок из файла хранилища
func LoadFromFile(filename string) ([]ShortenURL, error) {
	c, err := newConsumer(filename)
	if err != nil {
		return nil, err
	}
	defer c.close()

	var items []ShortenURL
	for {
		if item, _ := c.readItem(); item != nil {
			uuid, err := strconv.Atoi(item.UUID)
			if err != nil {
				uuid = 1
			}
			lastUUID = uuid
			items = append(items, *item)
		} else {
			break
		}
	}

	return items, nil
}

// SaveToFile Добавить одну ссылку в файл хранилище
func SaveToFile(filename string, key, link string) error {
	p, err := newProducer(filename)
	if err != nil {
		return err
	}
	defer p.close()

	lastUUID++
	item := &ShortenURL{
		UUID:        strconv.Itoa(lastUUID),
		ShortURL:    key,
		OriginalURL: link,
	}
	return p.writeItem(item)
}
