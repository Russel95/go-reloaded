package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: go run . input.txt output.txt")
		os.Exit(1)
	}

	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	result, err := corrigerTexte(string(input))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := os.WriteFile(os.Args[2], []byte(result), 0644); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func corrigerTexte(texte string) (string, error) {
	mots := strings.Fields(texte)
	sortie := make([]string, 0, len(mots))

	for i := 0; i < len(mots); i++ {
		mot := mots[i]

		if i+1 < len(mots) && strings.Contains(mot, ",") {
			if action, nombre, ok := lireAction(mot, mots[i+1]); ok {
				if err := appliquerAction(&sortie, action, nombre); err != nil {
					return "", err
				}
				i++
				continue
			}
		}

		if action, nombre, ok := lireAction(mot, ""); ok {
			if err := appliquerAction(&sortie, action, nombre); err != nil {
				return "", err
			}
			continue
		}

		sortie = append(sortie, mot)
	}

	ajouterArticles(sortie)
	ajouterApostrophes(sortie)
	sortie = corrigerPonctuation(sortie)

	return strings.Join(sortie, " "), nil
}

func lireAction(mot string, prochain string) (string, int, bool) {
	switch mot {
	case "(hex)", "(bin)":
		return mot[1 : len(mot)-1], 1, true
	case "(up)", "(low)", "(cap)":
		return mot[1 : len(mot)-1], 1, true
	case "(up,", "(low,", "(cap,":
		if strings.HasSuffix(prochain, ")") {
			nombre, err := strconv.Atoi(strings.TrimSuffix(prochain, ")"))
			if err == nil && nombre > 0 {
				return strings.TrimSuffix(mot[1:], ","), nombre, true
			}
		}
	}

	for _, action := range []string{"up", "low", "cap"} {
		debut := "(" + action + ","
		if strings.HasPrefix(mot, debut) && strings.HasSuffix(mot, ")") {
			nombre, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(mot, debut), ")"))
			if err == nil && nombre > 0 {
				return action, nombre, true
			}
		}
	}

	return "", 0, false
}

func appliquerAction(sortie *[]string, action string, nombre int) error {
	indices := make([]int, 0, nombre)
	for i := len(*sortie) - 1; i >= 0 && len(indices) < nombre; i-- {
		if estMot((*sortie)[i]) {
			indices = append(indices, i)
		}
	}

	if len(indices) < nombre {
		return fmt.Errorf("pas assez de mots avant (%s)", action)
	}

	for _, i := range indices {
		mot := (*sortie)[i]
		switch action {
		case "hex", "bin":
			converti, err := convertirNombre(mot, action)
			if err != nil {
				return err
			}
			(*sortie)[i] = converti
		case "up":
			(*sortie)[i] = changerMot(mot, strings.ToUpper)
		case "low":
			(*sortie)[i] = changerMot(mot, strings.ToLower)
		case "cap":
			(*sortie)[i] = changerMot(mot, premiereMajuscule)
		}
	}

	return nil
}

func convertirNombre(mot string, base string) (string, error) {
	debut, fin := decouperMot(mot)
	nombre, err := strconv.ParseInt(mot[debut:fin], map[string]int{"hex": 16, "bin": 2}[base], 64)
	if err != nil {
		return "", fmt.Errorf("nombre %q invalide pour (%s)", mot[debut:fin], base)
	}
	return mot[:debut] + strconv.FormatInt(nombre, 10) + mot[fin:], nil
}

func changerMot(mot string, changer func(string) string) string {
	debut, fin := decouperMot(mot)
	return mot[:debut] + changer(mot[debut:fin]) + mot[fin:]
}

func premiereMajuscule(mot string) string {
	runes := []rune(strings.ToLower(mot))
	if len(runes) == 0 {
		return mot
	}
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func decouperMot(mot string) (int, int) {
	debut := 0
	for debut < len(mot) && mot[debut] == '\'' {
		debut++
	}

	fin := len(mot)
	for fin > debut && estPonctuation(rune(mot[fin-1])) {
		fin--
	}

	return debut, fin
}

func ajouterArticles(mots []string) {
	for i := 0; i+1 < len(mots); i++ {
		if mots[i] != "a" && mots[i] != "A" {
			continue
		}

		suivant := strings.TrimLeft(mots[i+1], "'\"([{")
		if suivant == "" {
			continue
		}

		premier := unicode.ToLower([]rune(suivant)[0])
		if strings.ContainsRune("aeiouh", premier) {
			if mots[i] == "A" {
				mots[i] = "An"
			} else {
				mots[i] = "an"
			}
		}
	}
}

func ajouterApostrophes(mots []string) {
	ouvre := false
	for i := 0; i < len(mots); i++ {
		if mots[i] != "'" {
			continue
		}

		if ouvre {
			if i > 0 {
				mots[i-1] += "'"
			}
		} else if i+1 < len(mots) {
			mots[i+1] = "'" + mots[i+1]
		}
		ouvre = !ouvre
		mots[i] = ""
	}
}

func corrigerPonctuation(mots []string) []string {
	sortie := make([]string, 0, len(mots))
	for _, mot := range mots {
		if mot == "" {
			continue
		}

		partie := strings.Builder{}
		ponctuation := strings.Builder{}
		for _, lettre := range mot {
			if estPonctuation(lettre) {
				if partie.Len() > 0 {
					sortie = append(sortie, partie.String())
					partie.Reset()
				}
				ponctuation.WriteRune(lettre)
			} else {
				if ponctuation.Len() > 0 {
					ajouterPonctuation(&sortie, ponctuation.String())
					ponctuation.Reset()
				}
				partie.WriteRune(lettre)
			}
		}

		if partie.Len() > 0 {
			sortie = append(sortie, partie.String())
		}
		if ponctuation.Len() > 0 {
			ajouterPonctuation(&sortie, ponctuation.String())
		}
	}

	return sortie
}

func ajouterPonctuation(mots *[]string, ponctuation string) {
	if len(*mots) == 0 {
		*mots = append(*mots, ponctuation)
		return
	}
	(*mots)[len(*mots)-1] += ponctuation
}

func estPonctuation(lettre rune) bool {
	return strings.ContainsRune(".,!?;:", lettre)
}

func estMot(mot string) bool {
	for _, lettre := range mot {
		if !estPonctuation(lettre) && lettre != '\'' {
			return true
		}
	}
	return false
}
