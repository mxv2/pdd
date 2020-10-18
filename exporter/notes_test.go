package exporter

import "testing"

var (
	htmlNote = Note{
		deckID: 1595179860403,
		front:  "При приближении к остановившемуся транспортному средству с включенной аварийной сигнализацией, которое имеет опознавательные знаки \"Перевозка детей\", водитель должен<div><br></div><div>* Снизить скорость<div>* При необходимости остановиться и пропустить детей</div><div>* Осуществить все перечисленные действия</div></div>",
		back:   "Осуществить все перечисленные действия",
	}
	imageNote = Note{
		deckID: 1595179860403,
		front:  "<img src=\"8oW0AXl.jpg\"><div><br></div><div>Разрешается ли водителю выполнить объезд грузового автомобиля?<br></div><div><br></div><div>* Разрешается</div><div>* Разрешается, если между  шлагбаумом и остановившимся грузовым автомобилем  расстояние более 5 м</div><div>* Запрещается</div>",
		back:   "Запрещается",
	}
)

func TestNote_GUID(t *testing.T) {
	tests := []struct {
		name string
		note Note
	}{
		{
			"simple note",
			htmlNote,
		},
		{
			"html note",
			imageNote,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.note.GUID(); len(got) != 10 {
				t.Errorf("GUID() = %v, len(%v) = %d, want 10 chars", got, got, len(got))
			}
		})
	}
}

func TestNote_Checksum(t *testing.T) {
	tests := []struct {
		name string
		note Note
		want uint32
	}{
		{
			"simple note",
			htmlNote,
			2001113794,
		},
		{
			"html note",
			imageNote,
			3166763790,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.note.Checksum(); got != tt.want {
				t.Errorf("Checksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stripHtmlPreservingImage(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			"note with html tags",
			htmlNote.front,
			"При приближении к остановившемуся транспортному средству с включенной аварийной сигнализацией, которое имеет опознавательные знаки \"Перевозка детей\", водитель должен* Снизить скорость* При необходимости остановиться и пропустить детей* Осуществить все перечисленные действия",
		},
		{
			"note with html tags and image link",
			imageNote.front,
			" 8oW0AXl.jpg Разрешается ли водителю выполнить объезд грузового автомобиля?* Разрешается* Разрешается, если между  шлагбаумом и остановившимся грузовым автомобилем  расстояние более 5 м* Запрещается",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripHtmlPreservingImage(tt.s); got != tt.want {
				t.Errorf("stripHtmlPreservingImage() = %v, want %v", got, tt.want)
			}
		})
	}
}
