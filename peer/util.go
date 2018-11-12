package peer

func GetFromUserMap(mapData map[string]User) []User {
	values := make([]User, len(mapData))
	idx := 0
	for _, value := range mapData {
		values[idx] = value
		idx++
	}
	return values
}
