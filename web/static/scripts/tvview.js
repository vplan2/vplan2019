'use strict';

const refreshTime = 10000;

var tvViewItems = [];
var hasSwapped = false;

function checkOverflow(element, buffer) {
    let rect = element.getBoundingClientRect();
    if (!rect)
        return false;
    return (rect.bottom + buffer >= window.innerHeight);
}

function createVplanEntryTVView(id, entry) {
	var list_item = document.createElement("li");
	list_item.setAttribute('class', 'list-group-item d-flex justify-content-between lh-condensed');

	var desc = document.createElement("div");
	var head = document.createElement("h6");
	head.setAttribute('class', 'my-0 font-weight-bold');
	head.textContent = entry.class;
	desc.appendChild(head);
	var tiny = document.createElement("span");
	tiny.textContent = entry.messures;
	desc.appendChild(tiny);
	list_item.appendChild(desc);

	var span = document.createElement("span");
	span.setAttribute('class', 'text-muted');
	span.textContent = entry.time + ' - ' + entry.responsible;
	list_item.appendChild(span);

    _(id).appendChild(list_item);
    
    if (checkOverflow(list_item, 75)) {
        list_item.style.cssText = 'display: none !important';
        tvViewItems.push({ i: list_item, v: false });
    } else {
        tvViewItems.push({ i: list_item, v: true });
    }

	return list_item;
}

function shuffleItems(cb) {
    let nonVisible = [];
    let visible = [];

    tvViewItems.forEach(function(item) {
        if (item.v) visible.push(item);
        else nonVisible.push(item);
    });

    if (nonVisible.length > 0) {
        for (let i = 0; i < visible.length; i++) {
            let item = visible[i];
            item.i.style.cssText = 'display: none !important';
            tvViewItems.splice(tvViewItems.indexOf(item), 1);
        }
        nonVisible.forEach(function(item) {
            item.i.style.cssText = 'display: flex !important';
            if (checkOverflow(item.i, 75)) {
                item.i.style.cssText = 'display: none !important';
                return;
            }
            item.v = true;
        });
        hasSwapped = true;
        return;
    }

    if (hasSwapped && cb) {
        hasSwapped = false;
        cb();
    }
}

function getDataForVplanTVView(method, url, args) {
	getJson(method, url, args, function() {
		// console.log(this);
		if(this.data != undefined) {
			_("day0").innerHTML = formatDate(this.data[0].date_for);

			// COLUMN DAY 1
			_("day0e").innerHTML = '';
			(this.data[0].header != '') ? createVplanEntryHeader('day0e', this.data[0].header) : console.log('header field is empty');
			(this.data[0].entries == null) ? console.log(this.data[0].entries) : this.data[0].entries.forEach( function(entry) { createVplanEntryTVView("day0e", entry); });
			(this.data[0].footer != '') ? createVplanEntryHeader('day0e', this.data[0].footer) : console.log('footer field is empty');
			_("day1").innerHTML = formatDate(this.data[1].date_for);

			// COLUMN DAY 2
			_("day1e").innerHTML = '';
			(this.data[1].header != '') ? createVplanEntryHeader('day1e', this.data[1].header) : console.log('header field is empty');
			(this.data[1].entries == null) ? console.log(this.data[1].entries) : this.data[1].entries.forEach( function(entry) { createVplanEntryTVView("day1e", entry); });
			(this.data[1].footer != '') ? createVplanEntryHeader('day1e', this.data[1].footer) : console.log('footer field is empty');
			_("day2").innerHTML = formatDate(this.data[2].date_for);

			// COLUMN DAY 3
			_("day2e").innerHTML = '';
			(this.data[2].header != '') ? createVplanEntryHeader('day2e', this.data[2].header) : console.log('header field is empty');
			(this.data[2].entries == null) ? console.log(this.data[2].entries) : this.data[2].entries.forEach( function(entry) { createVplanEntryTVView("day2e", entry); });
			(this.data[2].footer != '') ? createVplanEntryHeader('day2e', this.data[2].footer) : console.log('footer field is empty');

            let timer = setInterval(function() {
                shuffleItems(function() {
                    getDataForVplanTVView(method, url, args);
                    clearInterval(timer);
                });
            }, refreshTime);

		} else if(this.error != undefined) {
			console.log(this.error.code)
		} else {
			// TODO
		}
	});
	// setTimeout(function() {getDataForVplan(method, url, args);}, 20000);
}