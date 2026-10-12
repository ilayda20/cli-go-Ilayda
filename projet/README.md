# Correcteur de textes

Petit outil en ligne de commande, écrit en Go, qui corrige les comptes rendus de l'association
« Les Plumes de Serris » : il lit un fichier texte contenant des marqueurs comme `(up)` ou `(hex)`,
les applique, puis écrit le texte corrigé dans un second fichier.

## Installation

Prérequis : Go 1.22 ou plus récent (`go version`) et Git.

```
git clone https://github.com/arthurcommanay/cli-go-Arthur-Commanay.git
cd cli-go-Arthur-Commanay/projet
```

Toutes les commandes ci-dessous se lancent depuis ce dossier `projet/`.

## Utilisation

```
go run . <entree> <sortie>
```

Exemple :

```
go run . tests/T03-entree.txt sortie.txt
```

| Entrée (`tests/T03-entree.txt`)                      | Sortie attendue (`sortie.txt`)                   |
|------------------------------------------------------|--------------------------------------------------|
| `La présidente marie (cap) DUPONT (cap) ouvre la séance.` | `La présidente Marie Dupont ouvre la séance.` |

Erreurs : les messages s'affichent sur la sortie d'erreur et le programme renvoie le code 1
(nombre d'arguments incorrect, fichier d'entrée introuvable, entrée identique à la sortie).

## Règles

| Marqueur / cas      | Effet                                                        | Avant → après                                   |
|---------------------|--------------------------------------------------------------|-------------------------------------------------|
| `(up)`              | met le mot précédent en majuscules                           | `club (up)` → `CLUB`                            |
| `(low)`             | met le mot précédent en minuscules                           | `PAS (low)` → `pas`                             |
| `(cap)`             | met la première lettre en majuscule, le reste en minuscules  | `marie (cap)` → `Marie`                         |
| `(hex)`             | convertit le mot précédent d'hexadécimal en décimal          | `1E (hex)` → `30`                               |
| `(bin)`             | convertit le mot précédent de binaire en décimal             | `1010 (bin)` → `10`                             |
| `(up, 2)` `(cap, 4)`| applique l'action aux n mots précédents                      | `ateliers d'écriture (up, 2)` → `ATELIERS D'ÉCRITURE` |
| Ponctuation         | collée au mot précédent, espace après                        | `Bonjour , chers` → `Bonjour, chers`            |
| Guillemets `' '`    | espaces retirés à l'intérieur                                | `' souvenirs '` → `'souvenirs'`                 |
| Début de phrase     | majuscule après `.`, `!` ou `?`                              | `. le budget` → `. Le budget`                   |

Cas limites : un marqueur inconnu (`(gras)`) est conservé tel quel, un code invalide (`ZZ (hex)`)
reste intact, les espaces multiples sont réduits, les lignes vides sont conservées.

## État d'avancement

Les fonctions de ces règles existent dans `transform/transform.go` (`nettoyer`, `decouper`,
`appliquerMarqueurs`, `assembler`, `majuscules`) mais **ne sont pas encore branchées** sur
`transform.Texte`, qui applique pour l'instant le traitement de logs du fil rouge `trace`
(lignes à 6 champs). Tant que ce branchement n'est pas fait, `go run .` sur un texte libre
affiche un avertissement par ligne et recopie le texte sans le corriger.

## Architecture

```
main.go          vérifie les arguments, enchaîne lecture → transformation → écriture
 ├── fileio      Lire, Ecrire                       (lecture et écriture de fichiers)
 ├── transform   Texte, Ligne, marqueurs, …         (règles, fonctions pures)
 └── parser      Analyser, Entree                   (découpage et validation d'une ligne de log)
```

Pour modifier une règle de correction, c'est dans `transform/transform.go`.

## Tests

Le dossier `tests/` contient des paires `XXX-entree.txt` / `XXX-attendu.txt` :

```
go run . tests/T01-entree.txt sortie.txt
diff sortie.txt tests/T01-attendu.txt        # Windows : fc sortie.txt tests\T01-attendu.txt
```

Aucune différence = test réussi.

| Fichiers | Couvrent                                                                        |
|----------|---------------------------------------------------------------------------------|
| T01–T08  | règles du socle : `(up)`, `(low)`, `(cap)`, `(hex)`, `(bin)`, ponctuation, texte complet |
| L01–L05  | cas limites : marqueur en début de ligne, code invalide, marqueur inconnu, fichier vide, espaces multiples |
| X01–X04  | extensions : marqueur avec nombre, guillemets, majuscules de phrase, combinaisons |

`exemple.txt` est un vrai compte rendu à corriger : `go run . exemple.txt sortie.txt`.

## Auteur

 Ilayda Eruzun - B1 Informatique, Ynov Val d'Europe - 2026
