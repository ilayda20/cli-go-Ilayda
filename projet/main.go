package main

import (
	"errors"
	"fmt"
	"os"
	"projet/fileio"
	"projet/transform"
)

func main() {
	err := lancer(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "trace :", err)
		os.Exit(1)
	}
}

// lancer vérifie les arguments puis enchaîne lecture, transformation et écriture.
func lancer(args []string) error {
	if len(args) != 2 {
		return errors.New("usage : go run . <entree> <sortie>")
	}
	entree, sortie := args[0], args[1]
	if entree == sortie {
		return errors.New("l'entrée et la sortie doivent être deux fichiers  différents")
	}
	texte, err := fileio.Lire(entree)
	if err != nil {
		return err
	}

	resultat, avertissements := transform.Texte(texte)
	for _, a := range avertissements {
		fmt.Fprintln(os.Stderr, "attention :", a)
	}

	err = fileio.Ecrire(sortie, resultat)
	if err != nil {
		return err
	}

	fmt.Println("Transformation terminée :", sortie)
	return nil
}
