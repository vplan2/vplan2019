+++
title = "Frontend"
weight = 4
sort_by = "weight"
insert_anchor_links = "right"
+++

## System:
- if(possible) [Bootstrap] for simple maintance / change of themes / crossbrowser support
  -> https://medialoot.com/preview/bootstrap-4-dashboard/widgets.html ?
- Mobile First

## Theme:
- ([Material]-) [Photon] design system one of [Awesome Design Systems]
- [Inclusive Web Design Checklist]

## HTML5:
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

## (Web)APP:
- [ServiceWorker][sw] (requires [TLS]) for retrieving and storing data
- [WebAssembly][wasm] written with [Rust] to boost performamce of JavaScript

## Sites / REST:
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

## Calls:
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

## Supported Browsers:
- [Firefox] 66+ with flag: `svg.context-properties.content.enabled` set to `true`