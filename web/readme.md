### `web`

Dieser Ordner enthält alle Frontend-Source-Files, welche von dem Webserver exposed werden.

Um die statischen Dateien zu generieren wird [Zola] verwendet `zola build`,
dabei landen im Ordner: `public` die generierten Dateien,
für die Entwicklung kann man mit `zola serve --port 8080` auch eine Live-Vorschau erhalten.

[zola]: https://github.com/getzola/zola "A fast static site generator in a single binary with everything built-in."
