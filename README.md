# go-parallel

## DESCRIPTION

Exécute N fois une commande dont les données nécessaires sont stockées dans un fichier externe. Par défaut, la commande est exécutée 4 fois (-nbr) si le nombre de paramètres présents dans le fichier externe est suffisant. Par défaut, le nombre d'exécutions parallèles correspond à la valeur de -nbr. Si -noWait est activé, alors on ne traite plus par paquet de -nbr, mais de façon constante dont la valeur max est -nbr.

L'application va parcourir ligne par ligne le fichier et, en fonction du séparateur, va remplacer les entrées %[N]s par la ou les valeurs.

## EXEMPLE

```
# Fichier 'sample' contenant les valeurs suivantes :
5
15
hello
5

# cmd :
/usr/bin/sleep %[1]s

Il faut mettre le path complet de la commande que l'on utilise.

# La commande complête :
$./parallel.go -cmd="/usr/bin/sleep %[1]s" -file=sample

# JSON retouné par la commande :
{
  "total": 4,
  "skip": 1,
  "succes": 2,
  "fail": 1
}

```

## OPTIONS

| OPTIONS | TYPE | DÉFAUT | OBLIGATOIRE | DESCRIPTION |
|:-------:|:----:|:------:|:-----------:|:-----------:|
| cmd | string | | oui | Commande à traiter (path complet) |
| file | string | | oui | Fichier contenant les données |
| sep | string | , | | Séparateur utilisé dans le fichier |
| nbr | int | 4 | | Nombre de commandes parallélisées |
| timeout | int | 10 | | Timeout d'exécution d'une commande |
| noWait | | | | Désactive le traitement par lot et active le traitement constant. Le traitement est plus rapide |
| debug | | | | Mode debug |

## EXEMPLES
```
# Command basique
$ parallel -cmd="COMMANDE À TRAITER" -file="FICHIER CONTENANT DONNÉES"
```