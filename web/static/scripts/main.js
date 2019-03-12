/*
 * Main Script utf8
 * @author mosscap
 * @version 4.0a
 */
'use strict';
function _(id) {
	return document.getElementById(id);
}
function dgebtn(tagName) {
	return document.getElementsByTagName(tagName);
}
function dgebcn(className) {
	return document.getElementsByClassName(className);
}
// remove-method for elements
Element.prototype.remove = function() {
	this.parentElement.removeChild(this);
}
NodeList.prototype.remove = HTMLCollection.prototype.remove = function() {
	for(var i = this.length - 1; i >= 0; i--) {
		if(this[i] && this[i].parentElement) {
			this[i].parentElement.removeChild(this[i]);
		}
	}
}
function ajaxObj(meth, url) {
	if(window.XMLHttpRequest) { // code for IE7+, Firefox, Chrome, Opera, Safari
		x = new XMLHttpRequest();
	} else { // code for IE6, IE5
		x = new ActiveXObject("Microsoft.XMLHTTP");
	}
	x.open(meth, url, true);
	x.setRequestHeader("Content-type", "application/json; charset=UTF-8");
	return x;
}
function ajaxReturn(x) {
	if(x.readyState == 4) { //  && (x.status == 200 || x.status == 401)
		return true;
	}
}
function getJson(type, url, args, callback) {
	if(typeof callback === "function") {
		var ajax = ajaxObj(type, url);
		ajax.onreadystatechange = function() {
			if(ajaxReturn(ajax) == true) {
				console.log(ajax.responseText);
				callback.call(JSON.parse(ajax.responseText));
			} else {
			//	callback.call(JSON.parse('{"error": {"code": 425, "message": "Too Early"}}'));
			}
		}
		ajax.send(args);
	}
}
