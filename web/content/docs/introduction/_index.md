+++
title = "Introduction - Preamble"
weight = 0
sort_by = "weight"
insert_anchor_links = "right"

template = "docs-section.html"
page_template = "docs-page.html"
+++

## Preamble

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in [RFC 2119].


## Requirements

- [Rust] & [Go] environment
- [Git]
- [TLS] (needed for [ServiceWorker][sw] and for [HTTP/2][RFC 7540] because there are **no** implementations without it, see [1][n1] or [2][n2])
- Websockets (needed for push notifications)
- [Firefox] 60+
- knowledge about authentification


## Corporate design
- Logo [, icons]
- Color(s)
- Typographie
- [Images, ...]
- see [identihub] especially [Firefox Quantum](https://demo.identihub.co/firefox#/)


## Backend

Software:
- Development:
  - [Notepad++] / [Geany]
  - [Visual Studio Code] with Go-Package
  - [Insomia] a [REST] client
- Production:
  - Server or Hostsystem (Windows / Linux)

Server / CGI:
- (required) VPLAN2 [REST]-API
  - written in: [Go]
  - (optional) with [go-daemon]
- (optional) Web Push Server
  - written in: [Rust] (cargo tools: fmt, clippy, rls, update?)
  - (optional) with [rust-daemon]

Tasks:
- Authentification
- Autorisation
- [TLS] as spezified in [RFC 5246] or [RFC 8446]
- Application Programming Interface over 
  - [REST]
  - [CLI]? (like mysql or [sysinfo] example)
- HTTP Status Codes
  - [426] Browser support etc.
  - [418] Wrong Usage (specified in [RFC 2324])

Interface / REST:

```bash
$ curl -X GET <host>/api/vplan
<json-output with vplan infos>
$ curl -X GET <host>/rss/releases
<rss-output with release notes>
```

```bash
$ vplan -u <user> -p <pass>
Welcome to Vplan 2.0 CLI
> /help
<help-output>
> /GET /api/vplan
<json-output>
```


## WebFrontend

System:
- if(possible) [Bootstrap] for simple maintance / change of themes / crossbrowser support
- Mobile First

Theme:
- ([Material]-) [Photon] design system one of [Awesome Design Systems]
- [Inclusive Web Design Checklist]

HTML5:
- Formatting:
  - all tags have to be written like xhtml (closing slash)
  - head & body tags are not to be indented
  - script & style tags should be third elements in the xml hierarchy
  - script & style content should be avoided, but is not forbidden
  - script & style content should not be indented
  - for indention tabs should be used
- Mime-Types:
  - script: application/javascript or application/ecmascript
  - xhtml: application/xhtml+xml
  - html: text/html
  - css: text/css
- Tools:
  - Accessibility: [AXE]
  - [Combat Report] and [Open With]
  - [Html Validator]
  - [Lighthouse] (Apache License, Version 2.0)

(Web)APP:
- [ServiceWorker][sw] (requires [TLS]) for retrieving and storing data
- [WebAssembly][wasm] written with [Rust] to boost performamce of JavaScript

Sites / REST:
- index.html (with login logout)?
- vplan.html
- docs[/*].html
- about.html with Contact, License, GDPR(DSGVO)
- robots.txt
- updates.rss?
- sitemap.xml?
- favicon.ico (32\*32, 64\*64px icon, which must be in the root directory)
- media/ (image[\*.svg, \*.png, \*.webp], audio[\*.flac], video[\*.mkv, \*.webm])
- style/ (fonts, css, scss, icons - no images)
- scripts/ (\*.js, \*.es, \*.wasm, etc.)

Calls:
- 1:
  - loading HTML, CSS, JavaScript
  - displaying skeleton screen
  - installing serviceworker
    - loading and caching content
    - displaying content
- 2:
  - loading cache
  - calling for new infos
    - displays new infos or connection error note
    - update cache

Supported Browsers:
- [Firefox] 66+ with flag: `svg.context-properties.content.enabled` set to `true`


## Documentation

Software:
- [mdBook] or [Gutenberg] with our [design]
- [cloc] Count Lines of Code
- [umlet]

Content:
- Why was this decision made ?
  - What are the pros, what are the cons ?


## Presentation

Software:
- [reveal.js] ?! with our [design]


## SCM

[Git]-Repositories:
- [Backend](#Backend) with Docs
- [Frontend](#WebFrontend) with Docs (& [Presentation](#Presentation))
- Main with Docs and automated [Deployment](#Deployment)

[Git]-Hooks:
- verify code formatting
- verify build success (triggers ci ?)
- notify project member via eg.: [git-update-hook-discord]
- ...


## Deployment

[Git]-Hooks / CI:
- if build.success then project.deploy

[DEB]
- (optional) create own *.deb file like [Hello World][hello.deb]
- (optional) create own debian repository


## Bug/Issue Managment


## License

- **[MIT]** / [MPL 2.0] / [WTFPL]


### Licenses of dependencies

| Dependency   | License      |
| :----------- | :----------- |
| *frontend dependencies* |
| [Icons] | Mozilla Public License 2.0 |
| [Bootstrap] | MIT |
| [reveal.js] | MIT |
| *backend (go) dependencies* |
| github.com/ghodss/yaml | MIT |
| github.com/gorilla/mux | BSD 3-Clause "New" or "Revised" License |
| github.com/gorilla/sessions | BSD 3-Clause "New" or "Revised" License |
| github.com/michaeljs1990/sqlitestore | MIT |
| github.com/op/go-logging | BSD 3-Clause "New" or "Revised" License |
| github.com/zekroTJA/timedmap | MIT |
| golang.org/x/time | BSD 3-Clause "New" or "Revised" License |
| gopkg.in/yaml.v2 | Apache License 2.0 |
| github.com/BurntSushi/toml | MIT |


## Extendable/further projects

- App for Android/IOS/Firefox OS/Redox-OS/...
- Browser Addons/Webextensions
- WordPress-Plugin
- Add timetable
- Add classroom (Usage) Overview
- Add Live-Chat (SSE, WebRTC or XMPP with [JSXC])
- Add virtual classrooms (WebGL/WebVR)
- Rewrite of [Front](#WebFrontend) or [Backend](#Backend) (in [Rust])
- Add cli like wttr.in
- Add an issue management system (like Bugzilla, GitLab/GitHub Issues, etc.)
- ...


## References/Resourcen

- [theDesignProject](https://thedesignproject.co/)
- [corporate design](https://en.wikipedia.org/wiki/Corporate_design)
- [webapp.rs](https://github.com/saschagrunert/webapp.rs)
- https://rust.libhunt.com/mdbook-alternatives
- https://rust.libhunt.com/gutenberg-alternatives
- https://www.getzola.org/
- https://en.wikipedia.org/wiki/Free_software_license
- RFC's [RFC 2324], [RFC 2616], [RFC 6585], [RFC 7168], [RFC 7540], [RFC 5246], [RFC 8446]
- Web Push
  - https://serviceworke.rs/web-push.html
  - https://web-push-book.gauntface.com/
  - https://developers.google.com/web/fundamentals/push-notifications/
  - https://developer.mozilla.org/en-US/docs/Web/API/Push_API
  - https://blog.mozilla.org/services/2016/08/23/sending-vapid-identified-webpush-notifications-via-mozillas-push-service/
  - https://github.com/mozilla-services/push-service
  - https://github.com/mozilla-services/autopush-rs
  - https://github.com/mozilla-services/megaphone
  - https://github.com/pimeys/rust-web-push
  - https://github.com/web-push-libs/web-push-php
  - https://github.com/Minishlink/web-push-php-example


[Go]:       https://golang.com/      "a statically typed, compiled programming language designed at Google"
[go-daemon]:https://github.com/sevlyar/go-daemon    "A library for writing system daemons in golang."
[Rust]:     https://rust-lang.org/   "a multi-paradigm systems programming language focused on safety, especially safe concurrency"
[rust-daemon]: https://github.com/jwilm/rust_daemon_template "Template for writing daemons in Rust."

[mdBook]:   https://github.com/rust-lang-nursery/mdBook "a utility to create modern online books from Markdown files"
[Gutenberg]:https://github.com/getzola/zola         "a fast static site generator in a single binary with everything built-in"
[reveal.js]:https://revealjs.com/                   "The HTML Presentation Framework"

[Notepad++]:https://notepad-plus-plus.org/          "a free source code editor"
[Geany]:    https://www.geany.org/                  "a text editor using the GTK+ toolkit"
[Visual Studio Code]: https://code.visualstudio.com/
[Firefox]:  https://www.mozilla.org/en-US/firefox/  "a modern Browser with Gecko as renderengine"
[Insomia]:  https://insomnia.rest/                  "a REST client"
[REST]:     https://en.wikipedia.org/wiki/REST      "Representational state transfer"
[CLI]:      https://en.wikipedia.org/wiki/Command-line_interface "Command-line interface"
[design]:   #WebFrontend                            "an adaption of the Photon Design System"
[identihub]:https://identihub.co/                   "Open source hosting for your brand and visual asset"
[Icons]:    https://design.firefox.com/icons/viewer/
[Bootstrap]:https://getbootstrap.com/               "an open source toolkit for developing with HTML, CSS, and JS"
[Photon]:   https://design.firefox.com/             "Firefox Photon Design"
[Material]: https://material.io/design/             "Google Material Design"
[Awesome Design Systems]: https://github.com/alexpate/awesome-design-systems
[Inclusive Web Design Checklist]: https://github.com/Heydon/inclusive-design-checklist
[sw]:       https://serviceworke.rs/                "ServiceWorker"
[wasm]:     https://webassembly.org/                "WebAssembly"
[sysinfo]:  https://crates.io/crates/sysinfo        "A system handler to interact with processes."

[MIT]:      https://en.wikipedia.org/wiki/MIT_License
[MPL]:      https://en.wikipedia.org/wiki/Mozilla_Public_License "Mozilla Public License"
[MPL 2.0]:  https://www.mozilla.org/en-US/MPL/2.0/ "Mozilla Public License 2.0"
[WTFPL]:    http://www.wtfpl.net/   "Do What the Fuck You Want to Public License"
[Unlicense]:http://unlicense.org/   "Unlicense Yourself: Set Your Code Free"

[Git]:      https://git-scm.com/                    "a free and open source distributed version control system"
[git-update-hook-discord]: https://gist.github.com/zekroTJA/fed889517f02ab32a1a64ff1c8f2e77b
[cloc]:     https://github.com/AlDanial/cloc        "Count Lines of Code"
[umlet]:    http://www.umlet.com/                   "Free UML Tool for Fast UML Diagrams"

[RFC 2119]: https://tools.ietf.org/html/rfc2119     "Key words for use in RFCs to Indicate Requirement Levels"
[RFC 2324]: https://tools.ietf.org/html/rfc2324     "Hyper Text Coffee Pot Control Protocol (HTCPCP/1.0)"
[RFC 2616]: https://tools.ietf.org/html/rfc2616     "Hypertext Transfer Protocol -- HTTP/1.1"
[RFC 6585]: https://tools.ietf.org/html/rfc6585     "Additional HTTP Status Codes"
[RFC 7168]: https://tools.ietf.org/html/rfc7168     "The Hyper Text Coffee Pot Control Protocol for Tea Efflux Appliances (HTCPCP-TEA)"
[RFC 7540]: https://tools.ietf.org/html/rfc7540     "Hypertext Transfer Protocol Version 2 (HTTP/2)"
[n1]:       https://http2.github.io/faq/#does-http2-require-encryption      "Does HTTP/2 require encryption?"
[n2]:       https://wiki.mozilla.org/Networking/http2                       "Networking/http2"
[418]:      https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/418    "I'm a teapot"
[426]:      https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/426    "Upgrade Required"
[RFC 5246]: https://tools.ietf.org/html/rfc5246     "The Transport Layer Security (TLS) Protocol Version 1.2"
[RFC 8446]: https://tools.ietf.org/html/rfc8446     "The Transport Layer Security (TLS) Protocol Version 1.3"
[TLS]:      https://en.wikipedia.org/wiki/Transport_Layer_Security          "Transport Layer Security"

[AXE]:      https://addons.mozilla.org/de/firefox/addon/axe-devtools/
[Combat Report]: https://addons.mozilla.org/de/firefox/addon/compat-report/
[Open With]:https://addons.mozilla.org/firefox/addon/open-with/
[Html Validator]: https://addons.mozilla.org/de/firefox/addon/html-validator/
[Lighthouse]: https://chrome.google.com/webstore/detail/lighthouse/blipmdconlkpinefehnmjammfjpmpbjk

[hello.deb]:https://git.savannah.gnu.org/cgit/hello.git/tree/               "index : hello.git"
[JSXC]:     https://www.jsxc.org/
