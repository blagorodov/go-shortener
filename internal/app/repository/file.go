package repository

import (
	"strconv"
)

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
			lastUUID = uint(uuid)
			items = append(items, *item)
		} else {
			break
		}
	}

	return items, nil
}

func SaveToFile(filename string, item *ShortenURL) error {
	p, err := newProducer(filename)
	if err != nil {
		return err
	}
	if item.UUID == "" {
		lastUUID++
		item.UUID = strconv.Itoa(int(lastUUID))
	}
	return p.writeItem(item)
}
