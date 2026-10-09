package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var ponctuation = regexp.MustCompile(`\s+([.,!?;:]+)`)

func majusculeInitiale(mot string) string {
	lettres := []rune(strings.ToLower(mot))
	if len(lettres) > 0 {
		lettres[0] = unicode.ToUpper(lettres[0])
	}
	return string(lettres)
}

func corrigerLigne(ligne string) string {
	mots := strings.Fields(ligne)
	resultat := make([]string, 0, len(mots))

	for _, mot := range mots {
		switch mot {
		case "(up)", "(low)", "(cap)", "(hex)", "(bin)":
			if len(resultat) == 0 {
				continue
			}

			precedent := len(resultat) - 1

			switch mot {
			case "(up)":
				resultat[precedent] = strings.ToUpper(resultat[precedent])

			case "(low)":
				resultat[precedent] = strings.ToLower(resultat[precedent])

			case "(cap)":
				resultat[precedent] = majusculeInitiale(resultat[precedent])

			case "(hex)", "(bin)":
				base := 16
				if mot == "(bin)" {
					base = 2
				}

				nombre, err := strconv.ParseInt(resultat[precedent], base, 64)
				if err == nil {
					resultat[precedent] = strconv.FormatInt(nombre, 10)
				}
			}

		default:
			resultat = append(resultat, mot)
		}
	}

	phrase := strings.Join(resultat, " ")
	return ponctuation.ReplaceAllString(phrase, "$1")
}

func corrigerTexte(texte string) string {
	lignes := strings.Split(texte, "\n")

	for i, ligne := range lignes {
		finWindows := strings.HasSuffix(ligne, "\r")
		ligne = strings.TrimSuffix(ligne, "\r")

		lignes[i] = corrigerLigne(ligne)

		if finWindows {
			lignes[i] += "\r"
		}
	}

	return strings.Join(lignes, "\n")
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage : go run . entree.txt sortie.txt")
		os.Exit(2)
	}

	entree, err1 := filepath.Abs(os.Args[1])
	sortie, err2 := filepath.Abs(os.Args[2])

	if err1 != nil || err2 != nil || strings.EqualFold(entree, sortie) {
		fmt.Fprintln(os.Stderr, "Erreur : entree et sortie doivent etre des fichiers differents et valides")
		os.Exit(1)
	}

	contenu, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Erreur de lecture :", err)
		os.Exit(1)
	}

	texteCorrige := corrigerTexte(string(contenu))

	if err := os.WriteFile(os.Args[2], []byte(texteCorrige), 0644); err != nil {
		fmt.Fprintln(os.Stderr, "Erreur d'ecriture :", err)
		os.Exit(1)
	}

	fmt.Println("Fichier corrige :", os.Args[2])
}
