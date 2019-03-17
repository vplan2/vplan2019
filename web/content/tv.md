+++
title = "Vertretungsplan"
description = ""
template = "tv.html"
date = 2019-03-04T12:00:00

[extra]
stylesheet = "style/vplan.css"
script = "scripts/vplan.js"
+++

<div class="order-md-2 row">
	<div class="col">
		<h2 class="d-flex justify-content-between align-items-center mb-3"><span class="text-muted" id="day0">Heute</span></h2>
		<ul class="list-group mb-3" id="day0e">
		</ul>
	</div>
	<div class="col">
		<h2 class="d-flex justify-content-between align-items-center mb-3"><span class="text-muted" id="day1">Morgen</span></h2>
		<ul class="list-group mb-3" id="day1e">
		</ul>
	</div>
	<div class="col">
		<h2 class="d-flex justify-content-between align-items-center mb-3"><span class="text-muted" id="day2">Übermorgen</span></h2>
		<ul class="list-group mb-3" id="day2e">
		</ul>
	</div>
	<div class="mb-3 col" id="news">
		<div class="alert alert-secondary" role="alert">
		</div>
		<div class="alert alert-secondary" role="alert">
		</div>
		<div class="alert alert-secondary" role="alert">
		</div>
	</div>
</div>
