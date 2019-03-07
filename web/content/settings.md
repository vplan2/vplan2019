+++
title = "Settings"
description = ""
template = "blog.html"
date = 2017-04-01T10:00:00
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