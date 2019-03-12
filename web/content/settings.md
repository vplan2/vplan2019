+++
title = "Einstellungen"
description = ""
template = "settings.html"
date = 2019-03-11T10:00:00

[extra]
script = "scripts/logins.js"
+++

<div class="mb-3">
	<form class="form-signin" id="settings" method="POST" action="{{ get_url(path="api/authenticate/") | safe }}">
		<label for="class">Klasse</label>
		<input type="text" class="form-control" id="class" name="class" placeholder="ITF17B"/>
		<label for="theme">Theme</label>
		<input type="text" class="form-control" id="theme" name="theme" placeholder="dark|light" list="themes" />
		<datalist id="themes">
			<option value="dark">dark</option>
			<option value="light">light</option>
			<!-- … -->
		</datalist>
		<label for="edited">Editiert</label>
		<input type="text" class="form-control" id="edited" placeholder="datetime as defined in rfc 3339" disabled />
		<label for="submit">Senden</label>
		<input class="btn btn-lg btn-primary btn-block" type="submit" value="Speichern" id="submit"/>
	</form>
</div>
<div class="mb-3">
	<h4 class="d-flex justify-content-between align-items-center mb-3">Anmeldungen</h4>
	<ul class="list-group mb-3" id="logins">
	</ul>
</div>
