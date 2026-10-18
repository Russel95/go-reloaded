package main

import "testing"

func TestCorrigerTexte(t *testing.T) {
	tests := []struct {
		nom    string
		entree string
		voulu  string
	}{
		{
			nom:    "conversions",
			entree: "Simply add 42 (hex) and 10 (bin) and you will see the result is 68.",
			voulu:  "Simply add 66 and 2 and you will see the result is 68.",
		},
		{
			nom:    "majuscules et nombre",
			entree: "This is so exciting (up, 2)",
			voulu:  "This is SO EXCITING",
		},
		{
			nom:    "modifier avec ponctuation separee",
			entree: "one , two (up, 2)",
			voulu:  "ONE, TWO",
		},
		{
			nom:    "ponctuation",
			entree: "Punctuation tests are ... kinda boring ,what do you think ?",
			voulu:  "Punctuation tests are... kinda boring, what do you think?",
		},
		{
			nom:    "apostrophes",
			entree: "As Elton John said: ' I am the most well-known heterosexual in the world '",
			voulu:  "As Elton John said: 'I am the most well-known heterosexual in the world'",
		},
		{
			nom:    "a devient an",
			entree: "There it was. A amazing rock!",
			voulu:  "There it was. An amazing rock!",
		},
		{
			nom:    "modificateurs simples",
			entree: "it (cap) was the worst of times (up)",
			voulu:  "It was the worst of TIMES",
		},
	}

	for _, test := range tests {
		t.Run(test.nom, func(t *testing.T) {
			resultat, err := corrigerTexte(test.entree)
			if err != nil {
				t.Fatalf("corrigerTexte() erreur = %v", err)
			}
			if resultat != test.voulu {
				t.Errorf("corrigerTexte() = %q, voulu %q", resultat, test.voulu)
			}
		})
	}
}

func TestCorrigerTexteNombreInvalide(t *testing.T) {
	_, err := corrigerTexte("XX (hex)")
	if err == nil {
		t.Fatal("corrigerTexte() devait renvoyer une erreur")
	}
}
