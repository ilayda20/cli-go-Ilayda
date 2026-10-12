// Package parser découpe une ligne de log en champs et vérifie qu'elle est valide.
package parser

import (
	"errors"
	"fmt"
	"strings"
)

// ErrLigneInvalide signale une ligne qui ne respecte pas le format attendu.
var ErrLigneInvalide = errors.New("ligne invalide")

// NbChamps est le nombre de champs d'une ligne de log.
const NbChamps = 6

// Entree représente une ligne de log découpée en champs.
type Entree struct {
	Date, Heure, Niveau, Source, Code, Lieu string
}

// Analyser découpe une ligne en jetons et renvoie une Entree,
// ou une erreur si la ligne ne respecte pas le format.
func Analyser(ligne string) (Entree, error) {
	jetons := strings.Fields(ligne)
	if len(jetons) != NbChamps {
		return Entree{}, fmt.Errorf("%w : %d champs au lieu de %d",
			ErrLigneInvalide, len(jetons), NbChamps)
	}
	niveau := strings.ToLower(jetons[2])
	if niveau != "info" && niveau != "warn" && niveau != "error" {
		return Entree{}, fmt.Errorf("%w : niveau inconnu %q", ErrLigneInvalide,
			jetons[2])
	}
	return Entree{
		Date: jetons[0], Heure: jetons[1], Niveau: niveau,
		Source: jetons[3], Code: jetons[4], Lieu: jetons[5],
	}, nil
}
