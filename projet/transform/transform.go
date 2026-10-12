// Package transform applique les règles de correction au texte.
// Toutes ses fonctions sont pures : aucun fichier lu ou écrit, rien d'affiché.
package transform

import (
	"fmt"
	"projet/parser"
	"strconv"
	"strings"
	"unicode"
)

// Texte transforme chaque ligne du texte.
// Les lignes invalides sont recopiées telles quelles et signalées dans avertissements.
func Texte(texte string) (resultat string, avertissements []string) {
	lignes := strings.Split(texte, "\n")
	for i, ligne := range lignes {
		if strings.TrimSpace(ligne) == "" {
			continue // ligne vide : on la garde telle quelle
		}
		transformee, err := Ligne(ligne)
		if err != nil {
			avertissements = append(avertissements, fmt.Sprintf("ligne %d ignorée : %v", i+1, err))
			continue
		}
		lignes[i] = transformee
	}
	return strings.Join(lignes, "\n"), avertissements
}

// Ligne transforme une ligne de log valide.
func Ligne(ligne string) (string, error) {
	e, err := parser.Analyser(ligne)
	if err != nil {
		return "", err
	}
	champs := []string{e.Date, e.Heure, strings.ToUpper(e.Niveau),
		Capitaliser(e.Source), ConvertirCode(e.Code), e.Lieu}
	return strings.Join(champs, " "), nil
}

// Capitaliser met la première lettre en majuscule et le reste en minuscules.
func Capitaliser(mot string) string {
	lettres := []rune(strings.ToLower(mot))
	if len(lettres) == 0 {
		return mot
	}
	lettres[0] = unicode.ToUpper(lettres[0])
	return string(lettres)
}

// ConvertirCode convertit un code hexadécimal (0x1F4) en décimal.
// Un code invalide (0xZZ) est gardé tel quel.
func ConvertirCode(code string) string {
	if !strings.HasPrefix(strings.ToLower(code), "0x") {
		return code
	}
	return versDecimal(code, "0x", 16)
}

// nettoyer transforme "( up , 2 )" en "(up,2)" pour que le marqueur ne soit pas coupé en mots.
func nettoyer(l string) string {
	for i := 0; i < len(l); i++ {
		if l[i] != '(' {
			continue
		}
		fin := strings.Index(l[i:], ")")
		if fin == -1 {
			break
		}
		compact := strings.Join(strings.Fields(l[i:i+fin+1]), "")
		if _, _, ok := lireMarqueur(compact); ok {
			l = l[:i] + compact + l[i+fin+1:]
		}
	}
	return l
}

// decouper sépare la ligne en mots, et détache la ponctuation collée devant un mot.
func decouper(l string) []string {
	var mots []string
	for _, mot := range strings.Fields(l) {
		debut := strings.TrimLeft(mot, ".,!?:;")
		if debut != mot && debut != "" {
			mots = append(mots, mot[:len(mot)-len(debut)], debut)
		} else {
			mots = append(mots, mot)
		}
	}
	return mots
}

// lireMarqueur regarde si le mot est un marqueur comme (up) ou (cap,3).
// Elle renvoie l'action, le nombre de mots concernés, et true si c'est un marqueur.
func lireMarqueur(mot string) (string, int, bool) {
	if len(mot) < 3 || mot[0] != '(' || mot[len(mot)-1] != ')' {
		return "", 0, false
	}
	action, nombre, avecNombre := strings.Cut(mot[1:len(mot)-1], ",")
	if action != "up" && action != "low" && action != "cap" && action != "hex" && action != "bin" {
		return "", 0, false
	}
	n := 1
	if avecNombre {
		var err error
		n, err = strconv.Atoi(nombre)
		if err != nil || nombre[0] < '0' || nombre[0] > '9' {
			return "", 0, false
		}
	}
	return action, n, true
}

// appliquerMarqueurs retire les marqueurs et les applique aux n mots précédents.
func appliquerMarqueurs(mots []string) []string {
	var resultat []string
	for _, mot := range mots {
		action, n, ok := lireMarqueur(mot)
		if !ok {
			resultat = append(resultat, mot)
			continue
		}
		for i := len(resultat) - 1; i >= 0 && i >= len(resultat)-n; i-- {
			resultat[i] = appliquer(action, resultat[i])
		}
	}
	return resultat
}

// appliquer transforme un mot selon l'action du marqueur.
func appliquer(action, mot string) string {
	switch action {
	case "up":
		return strings.ToUpper(mot)
	case "low":
		return strings.ToLower(mot)
	case "cap":
		lettres := []rune(strings.ToLower(mot))
		lettres[0] = unicode.ToUpper(lettres[0])
		return string(lettres)
	case "hex":
		return versDecimal(mot, "0x", 16)
	default: // "bin"
		return versDecimal(mot, "0b", 2)
	}
}

// versDecimal convertit un nombre hexa/binaire en décimal (le mot reste intact si invalide).
func versDecimal(mot, prefixe string, base int) string {
	chiffres := strings.TrimPrefix(strings.ToLower(mot), prefixe)
	nombre, err := strconv.ParseInt(chiffres, base, 64)
	if err != nil {
		return mot
	}
	return strconv.FormatInt(nombre, 10)
}

// assembler recolle les mots avec des espaces, sauf avant la ponctuation
// et à l'intérieur des guillemets simples ' ... '.
func assembler(mots []string) string {
	var b strings.Builder
	guillemetOuvert, collerSuivant := false, false
	for _, mot := range mots {
		colle := collerSuivant
		collerSuivant = false
		if strings.Trim(mot, ".,!?:;") == "" { // ponctuation seule
			colle = true
		} else if mot == "'" {
			colle = guillemetOuvert // fermant : collé au mot d'avant
			collerSuivant = !colle  // ouvrant : le mot suivant sera collé
			guillemetOuvert = !guillemetOuvert
		}
		if b.Len() > 0 && !colle {
			b.WriteByte(' ')
		}
		b.WriteString(mot)
	}
	return b.String()
}

// majuscules met une majuscule à la première lettre de chaque phrase.
func majuscules(l string) string {
	lettres := []rune(l)
	debutPhrase := true
	for i, r := range lettres {
		if unicode.IsLetter(r) {
			if debutPhrase {
				lettres[i] = unicode.ToUpper(r)
			}
			debutPhrase = false
		} else if strings.ContainsRune(".!?", r) {
			debutPhrase = true
		}
	}
	return string(lettres)
}
