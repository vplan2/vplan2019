# Glossar
Generelle Definition und Beschreibung von im Produkt und Projekt verwendeten Begriffen.

## Rechtegruppen

- **user**  
  Allgemeiner Benutzer welcher nur Leseberechtigung hat.
- **author**  
  Benutzer mit Lese- und Schreibberechtigung (erstellen, editieren & löschen)
- **admin**  
  Benutzer welcher die selben Berechtigungen wie der `author` hat.  
  Zusätzlich sind je Benutzer dazu berechtigt, die Rollen von Benutzern zu ändern.

## Datenbak / Datenstruktur

- **entry**  
  Datenstruktur eines Eintrags in dem Vertretungsplan.  
  > PU = Öffentich Sichtbare Membereigenschaften  
  > PR = Öffentlich nicht Sichtbare Membereigenschaften
    - PU: Klasse
    - PU: Stunde
    - PU: Maßnahmen
    - PU: Verantwortlicher
    - *PR: Timestamp*
    - *PR: `author` Ident*

- **ident**  
  Identifikator, welcher ein-eindeutig einem Nutzer zuordbar ist.

---

© 2019 Justin Trommler, Richard Heidenreich, Ringo Hoffmann  
Covered by MIT Licence.