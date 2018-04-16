'use strict';

var Server = function() {

    this.wsServer = 'ws://' + gParam.ws_server; //172.24.222.54:8989';
	this.websocket;
	this.isConnect;
	this.openCallback;
	this.closeCallback;
	this.messageCallback;
	this.errorCallback;
}

Server.prototype = {

	connect:function() {

		console.log("ready to connect ...");
		if(this.websocket && (this.websocket.readyState == 0 || this.websocket.readyState == 1))
		{
			this.websocket.close();
		}

		this.isConnect = false;

		try {
			this.websocket = new WebSocket(this.wsServer);
		}
		catch (ex) {
			console.log(ex, "ERROR");
			return;
		}

		var that = this;
		this.websocket.onopen = function (evt) { that.onOpen(evt); };
		this.websocket.onclose = function (evt) { that.onClose(evt); };
		this.websocket.onmessage = function (evt) { that.onMessage(evt); };
		this.websocket.onerror = function (evt) { that.onError(evt); };
	},

	registerCallback:function(openCallback, closeCallback, messageCallback, errorCallback) {

		this.openCallback = openCallback;
		this.closeCallback = closeCallback;
		this.messageCallback = messageCallback;
		this.errorCallback = errorCallback;
	},

	onOpen:function(evt) {
		console.log("Connected to WebSocket server.");
		this.isConnect = true;
		var data = {}
		if (evt.data)
		{
			try
			{
				data = JSON.parse(evt.data)
			}
			catch(e)
			{
				console.log(e);
			}
		}
		if(this.openCallback)
		{
			this.openCallback(data)
		}
	},
	onClose:function(evt) {
		console.log("Disconnected");
		var data = {}
		if (evt.data)
		{
			try
			{
				data = JSON.parse(evt.data)
			}
			catch(e)
			{
				console.log(e);
			}
		}
		if(this.closeCallback)
		{
			this.closeCallback(data)
		}
	},
	onMessage:function(evt) {
		console.log('Retrieved data from server: ' + evt.data);
		var data = {}
		if (evt.data)
		{
			try
			{
				data = JSON.parse(evt.data)
			}
			catch(e)
			{
				console.log(e);
			}
		}
		if(this.messageCallback)
		{
			this.messageCallback(data)
		}
	},
	onError:function(evt) {
		console.log('Error occured: ' + evt.data);
		var data = {}
		if (evt.data)
		{
			try
			{
				data = JSON.parse(evt.data)
			}
			catch(e)
			{
				console.log(e);
			}
		}
		if(this.errorCallback)
		{
			this.errorCallback(data)
		}
	},

	quit:function() {
		console.log("ready to disconnect");
		this.websocket.close(1000, "");
	},

	sendCommand:function(jsonData, callback) {

		var jsonStr = ""
		var isOK = true;
		if (jsonData) {
			try{
				jsonStr = JSON.stringify(jsonData)
			}catch(e){
				console.log(e);
				isOK = false;
			}
		};
		console.log("send data:", jsonStr);
		try{
			this.websocket.send(jsonStr);
		}catch(e){
			console.log(e);
			isOK = false;
		}

		if(callback)
		{
			callback(isOK);
		}
	}
}
