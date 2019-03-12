+++
title = "Settings"
description = ""
template = "settings.html"
date = 2019-03-11T10:00:00

[extra]
script = "scripts/logins.js"
+++

<div class="mb-3">
	<div class="custom-control custom-radio">
		<input id="dark-theme" name="theme" type="radio" class="custom-control-input" checked="checked">
		<label class="custom-control-label" for="dark-theme">Dark-Theme</label>
	</div>
	<div class="custom-control custom-radio">
		<input id="light-theme" name="theme" type="radio" class="custom-control-input">
		<label class="custom-control-label" for="light-theme">Light-Theme</label>
	</div>
	<label for="theme">Class</label>
	<input type="text" class="form-control" id="theme" placeholder="ITF17B" list="classes" />
</div>
<div class="mb-3">
	<h4 class="d-flex justify-content-between align-items-center mb-3">Anmeldungen</h4>
	<ul class="list-group mb-3" id="logins">
	</ul>
</div>
