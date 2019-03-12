// 'use strict';
var months = ["Januar", "Februar", "März", "April", "Mai", "Juni", "Juli", "August", "September", "Oktober", "November", "Dezember"];
var days = ["Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Sonnabend"];

function formatDate(date) {
	var selectedDate = new Date(date);
	var dateFormat = date.split('-');
	return dateFormat[2].split('T')[0] + '. ' + months[selectedDate.getMonth()] + ' <span class="badge badge-secondary badge-pill" id="day1c">' + days[selectedDate.getDay()] + '</span>';
}

function createLoginEntry(id, entry) {
	var list_item = document.createElement("li");
	list_item.setAttribute('class', 'list-group-item d-flex justify-content-between lh-condensed');

	var desc = document.createElement("div");
	var head = document.createElement("h6");
	head.setAttribute('class', 'my-0');
	head.textContent = entry.useragent;
	desc.appendChild(head);
	var tiny = document.createElement("small");
	tiny.textContent = entry.time + ' - ' + entry.ipaddress;
	desc.appendChild(tiny);
	list_item.appendChild(desc);

	var span = document.createElement("span");
	span.setAttribute('class', 'text-muted');
	span.textContent = entry.ident;
	list_item.appendChild(span);

	_(id).appendChild(list_item);
}

function getLoginsData(method, url, args) {
	getJson(method, url, args, function() {
		console.log(this);
		if(this.data != undefined) {
			this.data.forEach( function(entry) { createLoginEntry("logins", entry); })
		} else if(this.error != undefined) {
			console.log(this.error.code)
		} else {
			// TODO
		}
	});
}
